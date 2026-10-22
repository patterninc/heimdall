// Package eksauth builds AWS and Kubernetes client configuration for EKS clusters.
//
// It assumes an optional IAM role, looks the cluster up with eks:DescribeCluster and
// authenticates to the Kubernetes API with the same presigned STS GetCallerIdentity
// bearer token that aws-iam-authenticator and `aws eks get-token` produce, so no
// kubeconfig file or aws CLI is needed.
package eksauth

import (
	"context"
	"encoding/base64"
	"fmt"
	"net/http"
	"regexp"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/eks"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
	smithyhttp "github.com/aws/smithy-go/transport/http"
	"k8s.io/client-go/rest"
)

const (
	// tokenPrefix is the prefix EKS expects on presigned STS bearer tokens.
	tokenPrefix = `k8s-aws-v1.`
	// clusterIDHeader binds the presigned request to one cluster.
	clusterIDHeader = `x-k8s-aws-id`
	// EKS accepts a token for 15 minutes after signing; refresh well before that.
	tokenRefreshAfter = 5 * time.Minute
	// assumed role session length; refreshed by the credentials cache before expiry.
	assumeRoleDuration = time.Hour
	// refresh assumed credentials this long before they expire, so a cached token
	// is never signed with credentials that expire before the token does.
	credentialsExpiryWindow = 10 * time.Minute
	clusterStatusActive     = `ACTIVE`
	defaultSessionName      = `heimdall`
	maxSessionNameLength    = 64
)

var rxSessionNameInvalid = regexp.MustCompile(`[^\w+=,.@-]`)

// Config identifies an EKS cluster and how to reach it.
type Config struct {
	// RoleARN is assumed before calling AWS when set; otherwise the default credentials are used.
	RoleARN string
	// Region overrides the region from the default AWS configuration when set.
	Region string
	// ClusterName is the EKS cluster name.
	ClusterName string
	// SessionName is the STS role session name (sanitised to the allowed character set).
	SessionName string
}

// AWSConfig loads the default AWS configuration, applies Region and, when RoleARN is set,
// replaces the credentials with an auto-refreshing assume-role provider.
func AWSConfig(ctx context.Context, c Config) (aws.Config, error) {

	opts := []func(*awsconfig.LoadOptions) error{}
	if c.Region != `` {
		opts = append(opts, awsconfig.WithRegion(c.Region))
	}

	cfg, err := awsconfig.LoadDefaultConfig(ctx, opts...)
	if err != nil {
		return aws.Config{}, fmt.Errorf("failed to load AWS config: %w", err)
	}

	if c.RoleARN != `` {
		provider := stscreds.NewAssumeRoleProvider(sts.NewFromConfig(cfg), c.RoleARN, func(o *stscreds.AssumeRoleOptions) {
			o.RoleSessionName = SessionName(c.SessionName)
			o.Duration = assumeRoleDuration
		})
		cfg.Credentials = aws.NewCredentialsCache(provider, func(o *aws.CredentialsCacheOptions) {
			o.ExpiryWindow = credentialsExpiryWindow
		})
	}

	return cfg, nil

}

// DescribeCluster calls eks:DescribeCluster. It does not check the cluster status: the API
// server stays usable while a cluster is UPDATING (for example during a version upgrade).
func DescribeCluster(ctx context.Context, cfg aws.Config, clusterName string) (*ekstypes.Cluster, error) {

	if clusterName == `` {
		return nil, fmt.Errorf("EKS cluster name is required")
	}

	out, err := eks.NewFromConfig(cfg).DescribeCluster(ctx, &eks.DescribeClusterInput{
		Name: aws.String(clusterName),
	})
	if err != nil {
		return nil, err
	}
	if out.Cluster == nil {
		return nil, fmt.Errorf("EKS cluster %q not found", clusterName)
	}

	return out.Cluster, nil

}

// DescribeActiveCluster calls eks:DescribeCluster and returns an error unless the cluster is ACTIVE.
func DescribeActiveCluster(ctx context.Context, cfg aws.Config, clusterName string) (*ekstypes.Cluster, error) {

	cluster, err := DescribeCluster(ctx, cfg, clusterName)
	if err != nil {
		return nil, err
	}

	if string(cluster.Status) != clusterStatusActive {
		return nil, fmt.Errorf("EKS cluster %q is not active: %s", clusterName, cluster.Status)
	}

	return cluster, nil

}

// RESTConfig returns a Kubernetes client configuration for an EKS cluster. Requests carry
// a presigned STS bearer token generated from cfg's credentials and refreshed before it expires.
func RESTConfig(cfg aws.Config, cluster *ekstypes.Cluster) (*rest.Config, error) {

	if cluster == nil || cluster.Name == nil {
		return nil, fmt.Errorf("EKS cluster description is empty")
	}
	if cluster.Endpoint == nil || *cluster.Endpoint == `` {
		return nil, fmt.Errorf("EKS cluster %q has no endpoint", *cluster.Name)
	}
	if cluster.CertificateAuthority == nil || cluster.CertificateAuthority.Data == nil {
		return nil, fmt.Errorf("EKS cluster %q has no certificate authority data", *cluster.Name)
	}

	caData, err := base64.StdEncoding.DecodeString(*cluster.CertificateAuthority.Data)
	if err != nil {
		return nil, fmt.Errorf("failed to decode certificate authority of EKS cluster %q: %w", *cluster.Name, err)
	}

	tokens := newTokenSource(&stsPresigner{client: sts.NewPresignClient(sts.NewFromConfig(cfg))}, *cluster.Name)

	return &rest.Config{
		Host:            *cluster.Endpoint,
		TLSClientConfig: rest.TLSClientConfig{CAData: caData},
		WrapTransport: func(rt http.RoundTripper) http.RoundTripper {
			return &bearerRoundTripper{next: rt, tokens: tokens}
		},
	}, nil

}

// NewRESTConfig assumes the role, describes the cluster and returns its Kubernetes client
// configuration. With requireActive it fails unless the cluster status is ACTIVE.
func NewRESTConfig(ctx context.Context, c Config, requireActive bool) (*rest.Config, error) {

	cfg, err := AWSConfig(ctx, c)
	if err != nil {
		return nil, err
	}

	describe := DescribeCluster
	if requireActive {
		describe = DescribeActiveCluster
	}

	cluster, err := describe(ctx, cfg, c.ClusterName)
	if err != nil {
		return nil, err
	}

	return RESTConfig(cfg, cluster)

}

// SessionName sanitises s into a valid STS role session name.
func SessionName(s string) string {
	s = rxSessionNameInvalid.ReplaceAllString(s, `-`)
	if len(s) > maxSessionNameLength {
		s = s[:maxSessionNameLength]
	}
	if len(s) < 2 {
		return defaultSessionName
	}
	return s
}

// presigner produces a presigned STS GetCallerIdentity URL bound to a cluster.
type presigner interface {
	presign(ctx context.Context, clusterName string) (string, error)
}

type stsPresigner struct {
	client *sts.PresignClient
}

func (p *stsPresigner) presign(ctx context.Context, clusterName string) (string, error) {
	req, err := p.client.PresignGetCallerIdentity(ctx, &sts.GetCallerIdentityInput{}, func(o *sts.PresignOptions) {
		o.ClientOptions = append(o.ClientOptions, func(so *sts.Options) {
			so.APIOptions = append(so.APIOptions,
				smithyhttp.AddHeaderValue(clusterIDHeader, clusterName),
				// hoisted into the query string by the presigner
				smithyhttp.AddHeaderValue(`X-Amz-Expires`, `60`),
			)
		})
	})
	if err != nil {
		return ``, err
	}
	return req.URL, nil
}

// tokenSource caches an EKS bearer token and regenerates it before it expires.
type tokenSource struct {
	presigner   presigner
	clusterName string
	now         func() time.Time

	mu      sync.Mutex
	token   string
	created time.Time
}

func newTokenSource(p presigner, clusterName string) *tokenSource {
	return &tokenSource{presigner: p, clusterName: clusterName, now: time.Now}
}

// Token returns a cached token, or a fresh one when the cached token is too old.
func (t *tokenSource) Token(ctx context.Context) (string, error) {

	t.mu.Lock()
	defer t.mu.Unlock()

	if t.token != `` && t.now().Sub(t.created) < tokenRefreshAfter {
		return t.token, nil
	}

	url, err := t.presigner.presign(ctx, t.clusterName)
	if err != nil {
		return ``, fmt.Errorf("failed to generate EKS token for cluster %q: %w", t.clusterName, err)
	}

	t.token = tokenPrefix + base64.RawURLEncoding.EncodeToString([]byte(url))
	t.created = t.now()

	return t.token, nil

}

// Invalidate drops the cached token so the next request signs a new one.
func (t *tokenSource) Invalidate() {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.token = ``
}

// bearerRoundTripper adds the EKS bearer token to every request.
type bearerRoundTripper struct {
	next   http.RoundTripper
	tokens *tokenSource
}

func (b *bearerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {

	token, err := b.tokens.Token(req.Context())
	if err != nil {
		return nil, err
	}

	authed := req.Clone(req.Context())
	authed.Header.Set(`Authorization`, `Bearer `+token)

	resp, err := b.next.RoundTrip(authed)
	if err == nil && resp.StatusCode == http.StatusUnauthorized {
		// a rejected token is not reused; the client's next request gets a new one
		b.tokens.Invalidate()
	}

	return resp, err

}
