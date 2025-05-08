package snowflake

import (
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"fmt"
	"os"

	sf "github.com/snowflakedb/gosnowflake"

	"github.com/patterninc/heimdall/internal/pkg/context"
	"github.com/patterninc/heimdall/internal/pkg/object/cluster"
	"github.com/patterninc/heimdall/internal/pkg/object/job"
	"github.com/patterninc/heimdall/internal/pkg/result"
	"github.com/patterninc/heimdall/pkg/plugin"
)

const (
	privateKeyType      = `PRIVATE KEY`
	snowflakeDriverName = `snowflake`
)

var (
	ErrFailedToDecodePEMBlock = fmt.Errorf(`failed to decode PEM block`)
	ErrInvalidKeyType         = fmt.Errorf(`invalida key type`)
)

type snowflakeCommandContext struct {
	Account    string `yaml:"account,omitempty" json:"account,omitempty"`
	User       string `yaml:"user,omitempty" json:"user,omitempty"`
	Database   string `yaml:"database,omitempty" json:"database,omitempty"`
	Warehouse  string `yaml:"warehouse,omitempty" json:"warehouse,omitempty"`
	PrivateKey string `yaml:"private_key,omitempty" json:"private_key,omitempty"`
	dsn        string
}

type snowflakeJobContext struct {
	Query string `yaml:"query,omitempty" json:"query,omitempty"`
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

func New(_ *context.Context) (plugin.Handler, error) {
	s := &snowflakeCommandContext{}
	return s.handler, nil
}

func (s *snowflakeCommandContext) handler(r *plugin.Runtime, j *job.Job, c *cluster.Cluster) error {

	// do we have dsn defined?
	if s.dsn == `` {
		if c.Context != nil {
			if err := c.Context.Unmarshal(s); err != nil {
				return err
			}
		}

		// prepare snowflake config
		privateKeyBytes, err := os.ReadFile(s.PrivateKey)
		if err != nil {
			return err
		}

		// Parse the private key
		privateKey, err := parsePrivateKey(privateKeyBytes)
		if err != nil {
			return err
		}

		if s.dsn, err = sf.DSN(&sf.Config{
			Account:       s.Account,
			User:          s.User,
			Database:      s.Database,
			Warehouse:     s.Warehouse,
			Authenticator: sf.AuthTypeJwt,
			PrivateKey:    privateKey,
			Application:   r.UserAgent,
		}); err != nil {
			return err
		}
	}

	// let's unmarshal job context
	sc := &snowflakeJobContext{}
	if err := j.Context.Unmarshal(sc); err != nil {
		return err
	}

	// open connection
	db, err := sql.Open(snowflakeDriverName, s.dsn)
	if err != nil {
		return err
	}
	defer db.Close()

	rows, err := db.Query(sc.Query)
	if err != nil {
		return err
	}

	if j.Result, err = result.FromRows(rows); err != nil {
		return err
	}

	return nil

}
