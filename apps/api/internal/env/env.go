// Package env reads Prisme's configuration from the environment once, at boot.
package env

import (
	"fmt"
	"strings"

	"github.com/FacileStudio/porte"
	troncenv "github.com/FacileStudio/tronc/env"
)

// Matrix is how the API reaches Synapse and what it tells browsers.
type Matrix struct {
	InternalURL string
	PublicURL   string
	ServerName  string
	ASToken     string
}

// Echo is the call integration. Both values empty turns calls off.
type Echo struct {
	URL string
	Key string
}

// Config is the whole process configuration.
type Config struct {
	troncenv.Core
	Porte       porte.Config
	Matrix      Matrix
	Echo        Echo
	SecretKey   string
	ServiceKeys string
}

// Load reads and validates the environment.
func Load() (Config, error) {
	core, err := troncenv.LoadCore()
	if err != nil {
		return Config{}, err
	}
	cfg := Config{
		Core:        core,
		Matrix:      loadMatrix(),
		Echo:        Echo{URL: strings.TrimRight(troncenv.String("ECHO_URL", ""), "/"), Key: troncenv.String("ECHO_INTEGRATION_KEY", "")},
		SecretKey:   troncenv.String("PRISME_SECRET_KEY", ""),
		ServiceKeys: troncenv.String("PRISME_SERVICE_KEY", ""),
	}
	if cfg.Porte, err = loadPorte(); err != nil {
		return Config{}, err
	}
	return cfg, cfg.validate()
}

func loadMatrix() Matrix {
	return Matrix{
		InternalURL: strings.TrimRight(troncenv.String("MATRIX_INTERNAL_URL", "http://127.0.0.1:8008"), "/"),
		PublicURL:   strings.TrimRight(troncenv.String("MATRIX_PUBLIC_URL", "http://localhost:8008"), "/"),
		ServerName:  troncenv.String("MATRIX_SERVER_NAME", "localhost"),
		ASToken:     troncenv.String("PRISME_AS_TOKEN", ""),
	}
}

func loadPorte() (porte.Config, error) {
	ssoOnly, err := troncenv.Bool("SSO_ONLY", true)
	if err != nil {
		return porte.Config{}, err
	}
	return porte.Config{
		Issuer:          troncenv.String("OIDC_ISSUER", ""),
		ClientID:        troncenv.String("OIDC_CLIENT_ID", ""),
		ClientSecret:    troncenv.String("OIDC_CLIENT_SECRET", ""),
		RedirectURL:     troncenv.String("OIDC_REDIRECT_URL", ""),
		SuccessURL:      troncenv.String("OIDC_SUCCESS_URL", ""),
		MachineAudience: troncenv.String("OIDC_MACHINE_AUDIENCE", ""),
		SSOOnly:         ssoOnly,
	}, nil
}

func (c Config) validate() error {
	if c.Matrix.ASToken == "" {
		return fmt.Errorf("PRISME_AS_TOKEN is required")
	}
	if len(c.SecretKey) < 32 {
		return fmt.Errorf("PRISME_SECRET_KEY must be at least 32 characters")
	}
	return nil
}

// CallsEnabled reports whether the Echo integration is configured.
func (c Config) CallsEnabled() bool { return c.Echo.URL != "" && c.Echo.Key != "" }
