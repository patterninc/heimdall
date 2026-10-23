package eksauth

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	ekstypes "github.com/aws/aws-sdk-go-v2/service/eks/types"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type countingPresigner struct {
	calls int
}

func (p *countingPresigner) presign(_ context.Context, clusterName string) (string, error) {
	p.calls++
	return "https://sts.example.com/?cluster=" + clusterName + "&n=" + string(rune('0'+p.calls)), nil
}

func TestTokenSourceCachesAndRefreshes(t *testing.T) {
	p := &countingPresigner{}
	ts := newTokenSource(p, "my-cluster")
	now := time.Unix(1000, 0)
	ts.now = func() time.Time { return now }

	first, err := ts.Token(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(first, tokenPrefix) {
		t.Fatalf("token %q missing prefix %q", first, tokenPrefix)
	}
	decoded, err := base64.RawURLEncoding.DecodeString(strings.TrimPrefix(first, tokenPrefix))
	if err != nil || !strings.Contains(string(decoded), "cluster=my-cluster") {
		t.Fatalf("token does not encode presigned URL: %q (%v)", decoded, err)
	}

	now = now.Add(tokenRefreshAfter - time.Second)
	if again, _ := ts.Token(context.Background()); again != first || p.calls != 1 {
		t.Fatalf("expected cached token, presign calls=%d", p.calls)
	}

	now = now.Add(2 * time.Second)
	if refreshed, _ := ts.Token(context.Background()); refreshed == first || p.calls != 2 {
		t.Fatalf("expected refreshed token, presign calls=%d", p.calls)
	}

	ts.Invalidate()
	if _, _ = ts.Token(context.Background()); p.calls != 3 {
		t.Fatalf("expected new token after invalidate, presign calls=%d", p.calls)
	}
}

type recordingTransport struct {
	auth   []string
	status int
}

func (r *recordingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	r.auth = append(r.auth, req.Header.Get("Authorization"))
	return &http.Response{StatusCode: r.status, Body: http.NoBody, Request: req}, nil
}

func TestBearerRoundTripperSetsTokenAndInvalidatesOn401(t *testing.T) {
	p := &countingPresigner{}
	next := &recordingTransport{status: http.StatusUnauthorized}
	rt := &bearerRoundTripper{next: next, tokens: newTokenSource(p, "c")}

	req, _ := http.NewRequest(http.MethodGet, "https://k8s.example.com/api", nil)
	if _, err := rt.RoundTrip(req); err != nil {
		t.Fatal(err)
	}
	if req.Header.Get("Authorization") != "" {
		t.Fatal("original request must not be mutated")
	}
	if !strings.HasPrefix(next.auth[0], "Bearer "+tokenPrefix) {
		t.Fatalf("unexpected Authorization header %q", next.auth[0])
	}

	next.status = http.StatusOK
	if _, err := rt.RoundTrip(req); err != nil {
		t.Fatal(err)
	}
	if p.calls != 2 {
		t.Fatalf("expected token to be regenerated after a 401, presign calls=%d", p.calls)
	}
}

func TestSTSPresignerBindsClusterAndExpiry(t *testing.T) {
	cfg := aws.Config{
		Region:      "us-west-2",
		Credentials: credentials.NewStaticCredentialsProvider("AKIDEXAMPLE", "secret", "session"),
	}
	p := &stsPresigner{client: sts.NewPresignClient(sts.NewFromConfig(cfg))}

	raw, err := p.presign(context.Background(), "my-cluster")
	if err != nil {
		t.Fatal(err)
	}
	u, err := url.Parse(raw)
	if err != nil {
		t.Fatal(err)
	}
	q := u.Query()
	if q.Get("Action") != "GetCallerIdentity" {
		t.Errorf("Action = %q", q.Get("Action"))
	}
	if q.Get("X-Amz-Expires") != "60" {
		t.Errorf("X-Amz-Expires = %q", q.Get("X-Amz-Expires"))
	}
	if !strings.Contains(q.Get("X-Amz-SignedHeaders"), clusterIDHeader) {
		t.Errorf("signed headers %q do not include %s", q.Get("X-Amz-SignedHeaders"), clusterIDHeader)
	}
	if !strings.Contains(u.Host, "us-west-2") {
		t.Errorf("expected regional STS endpoint, got %s", u.Host)
	}
}

func TestRESTConfig(t *testing.T) {
	ca := base64.StdEncoding.EncodeToString([]byte("ca-bytes"))
	cluster := &ekstypes.Cluster{
		Name:                 aws.String("c"),
		Endpoint:             aws.String("https://ABC.gr7.us-west-2.eks.amazonaws.com"),
		CertificateAuthority: &ekstypes.Certificate{Data: aws.String(ca)},
	}
	rc, err := RESTConfig(aws.Config{Region: "us-west-2"}, cluster)
	if err != nil {
		t.Fatal(err)
	}
	if rc.Host != *cluster.Endpoint || string(rc.TLSClientConfig.CAData) != "ca-bytes" || rc.WrapTransport == nil {
		t.Fatalf("unexpected rest config: %+v", rc)
	}

	cluster.CertificateAuthority = nil
	if _, err := RESTConfig(aws.Config{}, cluster); err == nil {
		t.Fatal("expected error for missing certificate authority")
	}
	if _, err := RESTConfig(aws.Config{}, nil); err == nil {
		t.Fatal("expected error for nil cluster")
	}
}

func TestSessionName(t *testing.T) {
	cases := map[string]string{
		"":                                     defaultSessionName,
		"x":                                    defaultSessionName,
		"heimdall-health":                      "heimdall-health",
		"2b1c0c5e-8f0a-4b5e-9d1e-0e0c9b9a1f00": "2b1c0c5e-8f0a-4b5e-9d1e-0e0c9b9a1f00",
		"a b/c":                                "a-b-c",
		strings.Repeat("a", 80):                strings.Repeat("a", maxSessionNameLength),
	}
	for in, want := range cases {
		if got := SessionName(in); got != want {
			t.Errorf("SessionName(%q) = %q, want %q", in, got, want)
		}
	}
}
