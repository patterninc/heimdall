package snowflake

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	sf "github.com/snowflakedb/gosnowflake"

	heimdallContext "github.com/patterninc/heimdall/pkg/context"
	"github.com/patterninc/heimdall/pkg/object/cluster"
	"github.com/patterninc/heimdall/pkg/object/job"
	"github.com/patterninc/heimdall/pkg/plugin"
	"github.com/patterninc/heimdall/pkg/result"
)

const (
	privateKeyType      = `PRIVATE KEY`
	snowflakeDriverName = `snowflake`
	oauthTokenTTL       = time.Hour
	oauthScope          = `session:role-any`
)

var (
	ErrFailedToDecodePEMBlock = fmt.Errorf(`failed to decode PEM block`)
	ErrInvalidKeyType         = fmt.Errorf(`invalid key type`)
	ErrMissingOAuthConfig     = fmt.Errorf(`impersonate_user requires oauth_private_key, oauth_issuer, and oauth_audience`)
	ErrMissingJobUser         = fmt.Errorf(`impersonate_user is enabled but job has no user identity`)
)

type commandContext struct {
	DefaultRole string `yaml:"default_role,omitempty" json:"default_role,omitempty"`
}

type jobContext struct {
	Query string `yaml:"query,omitempty" json:"query,omitempty"`
	Role  string `yaml:"role,omitempty" json:"role,omitempty"`
}

type clusterContext struct {
	Account    string `yaml:"account,omitempty" json:"account,omitempty"`
	User       string `yaml:"user,omitempty" json:"user,omitempty"`
	Database   string `yaml:"database,omitempty" json:"database,omitempty"`
	Warehouse  string `yaml:"warehouse,omitempty" json:"warehouse,omitempty"`
	PrivateKey      string `yaml:"private_key,omitempty" json:"private_key,omitempty"`
	ImpersonateUser bool   `yaml:"impersonate_user,omitempty" json:"impersonate_user,omitempty"`
	// OAuthPrivateKey is a path to the PEM-encoded RSA private key used to mint OAuth tokens.
	// Only required when impersonate_user is true.
	OAuthPrivateKey string `yaml:"oauth_private_key,omitempty" json:"oauth_private_key,omitempty"`
	// OAuthIssuer must match the EXTERNAL_OAUTH_ISSUER configured in Snowflake.
	OAuthIssuer string `yaml:"oauth_issuer,omitempty" json:"oauth_issuer,omitempty"`
	// OAuthAudience is the token audience, typically the Snowflake account URL.
	OAuthAudience string `yaml:"oauth_audience,omitempty" json:"oauth_audience,omitempty"`
}

func parsePrivateKey(privateKeyBytes []byte) (*rsa.PrivateKey, error) {

	// Decode the key into a "block"
	privateBlock, _ := pem.Decode(privateKeyBytes)
	if privateBlock == nil || privateBlock.Type != privateKeyType {
		return nil, ErrFailedToDecodePEMBlock
	}

	// Parse the private key from the block
	privateKey, err := x509.ParsePKCS8PrivateKey(privateBlock.Bytes)
	if err != nil {
		return nil, err
	}

	// Check the type of the key
	rsaPrivateKey, ok := privateKey.(*rsa.PrivateKey)
	if !ok {
		return nil, ErrInvalidKeyType
	}

	return rsaPrivateKey, nil

}

func (c *clusterContext) mintOAuthToken(username string) (string, error) {
	privateKeyBytes, err := os.ReadFile(c.OAuthPrivateKey)
	if err != nil {
		return "", err
	}

	privateKey, err := parsePrivateKey(privateKeyBytes)
	if err != nil {
		return "", err
	}

	now := time.Now()
	claims := jwt.MapClaims{
		"sub": username,
		"iss": c.OAuthIssuer,
		"aud": jwt.ClaimStrings{c.OAuthAudience},
		"iat": now.Unix(),
		"exp": now.Add(oauthTokenTTL).Unix(),
		"scp": oauthScope,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	return token.SignedString(privateKey)
}

func New(cmdCtx *heimdallContext.Context) (plugin.Handler, error) {
	s := &commandContext{}

	// Parse command context from YAML config (contains role configuration)
	if cmdCtx != nil {
		if err := cmdCtx.Unmarshal(s); err != nil {
			return nil, err
		}
	}

	return s, nil
}

// Execute implements the plugin.Handler interface
func (s *commandContext) Execute(ctx context.Context, r *plugin.Runtime, j *job.Job, c *cluster.Cluster) error {

	clusterContext := &clusterContext{}
	if c.Context != nil {
		if err := c.Context.Unmarshal(clusterContext); err != nil {
			return err
		}
	}

	// let's unmarshal job context
	jobContext := &jobContext{}
	if err := j.Context.Unmarshal(jobContext); err != nil {
		return err
	}

	// build snowflake DSN — either OAuth impersonation or key-pair JWT
	var dsn string
	if clusterContext.ImpersonateUser {
		if clusterContext.OAuthPrivateKey == `` || clusterContext.OAuthIssuer == `` || clusterContext.OAuthAudience == `` {
			return ErrMissingOAuthConfig
		}
		if j.User == `` {
			return ErrMissingJobUser
		}

		oauthToken, err := clusterContext.mintOAuthToken(j.User)
		if err != nil {
			return err
		}

		role := s.DefaultRole
		if jobContext.Role != `` {
			role = jobContext.Role
		}

		dsn, err = sf.DSN(&sf.Config{
			Account:       clusterContext.Account,
			User:          j.User,
			Authenticator: sf.AuthTypeOAuth,
			Token:         oauthToken,
			Application:   r.UserAgent,
			Database:      clusterContext.Database,
			Warehouse:     clusterContext.Warehouse,
			Role:          role,
		})
		if err != nil {
			return err
		}
	} else {
		// standard key-pair JWT auth as Heimdall service account
		privateKeyBytes, err := os.ReadFile(clusterContext.PrivateKey)
		if err != nil {
			return err
		}

		privateKey, err := parsePrivateKey(privateKeyBytes)
		if err != nil {
			return err
		}

		// s.DefaultRole from command context; empty string = Snowflake uses user's default role
		dsn, err = sf.DSN(&sf.Config{
			Account:       clusterContext.Account,
			User:          clusterContext.User,
			Database:      clusterContext.Database,
			Warehouse:     clusterContext.Warehouse,
			Role:          s.DefaultRole,
			Authenticator: sf.AuthTypeJwt,
			PrivateKey:    privateKey,
			Application:   r.UserAgent,
		})
		if err != nil {
			return err
		}
	}

	// open connection
	db, err := sql.Open(snowflakeDriverName, dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.QueryContext(ctx, jobContext.Query)
	if err != nil {
		return err
	}

	if j.Result, err = result.FromRows(rows); err != nil {
		return err
	}

	return nil

}

// HealthCheck implements the plugin.HealthChecker interface
func (s *commandContext) HealthCheck(ctx context.Context, c *cluster.Cluster) error {
	clusterCtx := &clusterContext{}
	if c.Context != nil {
		if err := c.Context.Unmarshal(clusterCtx); err != nil {
			return err
		}
	}

	privateKeyBytes, err := os.ReadFile(clusterCtx.PrivateKey)
	if err != nil {
		return err
	}

	privateKey, err := parsePrivateKey(privateKeyBytes)
	if err != nil {
		return err
	}

	dsn, err := sf.DSN(&sf.Config{
		Account:       clusterCtx.Account,
		User:          clusterCtx.User,
		Database:      clusterCtx.Database,
		Warehouse:     clusterCtx.Warehouse,
		Role:          s.DefaultRole,
		Authenticator: sf.AuthTypeJwt,
		PrivateKey:    privateKey,
	})
	if err != nil {
		return err
	}

	db, err := sql.Open(snowflakeDriverName, dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.QueryContext(ctx, fmt.Sprintf("SHOW WAREHOUSES LIKE '%s'", clusterCtx.Warehouse))
	if err != nil {
		return fmt.Errorf("snowflake SHOW WAREHOUSES failed: %w", err)
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return err
	}

	stateIdx, autoResumeIdx := -1, -1
	for i, col := range cols {
		switch col {
		case "state":
			stateIdx = i
		case "auto_resume":
			autoResumeIdx = i
		}
	}
	if stateIdx < 0 {
		return fmt.Errorf("snowflake SHOW WAREHOUSES: state column not found")
	}

	if !rows.Next() {
		return fmt.Errorf("snowflake warehouse %q not found", clusterCtx.Warehouse)
	}

	vals := make([]any, len(cols))
	ptrs := make([]any, len(cols))
	for i := range vals {
		ptrs[i] = &vals[i]
	}
	if err := rows.Scan(ptrs...); err != nil {
		return err
	}

	state := scanToString(vals[stateIdx])
	switch state {
	case "STARTED", "RESUMING":
		return nil
	case "SUSPENDED":
		if autoResumeIdx >= 0 && scanToString(vals[autoResumeIdx]) == "true" {
			return nil
		}
		return fmt.Errorf("snowflake warehouse %q is suspended and auto_resume is not enabled", clusterCtx.Warehouse)
	default:
		return fmt.Errorf("snowflake warehouse %q is not ready: state=%s", clusterCtx.Warehouse, state)
	}
}

func scanToString(v any) string {
	switch val := v.(type) {
	case string:
		return val
	case []byte:
		return string(val)
	case bool:
		if val {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", v)
	}
}

func (s *commandContext) Cleanup(ctx context.Context, jobID string, c *cluster.Cluster) error {
	// TODO: Implement cleanup if needed
	return nil
}
