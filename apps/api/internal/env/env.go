// Package env reads Prisme's configuration from the environment once, at boot.
package env

import (
	"cmp"
	"fmt"
	"slices"
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
	HSToken     string
}

// Echo is the call integration. Both values empty turns calls off.
type Echo struct {
	URL string
	Key string
}

// DefaultMobileRedirect is the deep link the Android and desktop apps register
// for porte's ?flow=mobile sign-in. It applies only while OIDC is on: porte
// refuses a mobile redirect without an issuer.
const DefaultMobileRedirect = "studio.facile.prisme://auth"

// appOrigins are the origins the installed apps serve their pages from: the
// Capacitor Android app, then the Tauri desktop app on macOS and Linux, then on
// Windows.
var appOrigins = []string{"https://localhost", "tauri://localhost", "http://tauri.localhost"}

// Config is the whole process configuration. PublicURL is where browsers and
// Echo reach Prisme; it defaults to the Matrix public URL, which shares the
// host.
type Config struct {
	troncenv.Core
	Porte       porte.Config
	Matrix      Matrix
	Echo        Echo
	PublicURL   string
	LiveKitURL  string
	SecretKey   string
	ServiceKeys string
	FCMAccount  string
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
	cfg.LiveKitURL = strings.TrimRight(troncenv.String("LIVEKIT_URL", ""), "/")
	cfg.PublicURL = strings.TrimRight(troncenv.String("PRISME_PUBLIC_URL", cfg.Matrix.PublicURL), "/")
	cfg.Porte = loadPorte()
	cfg.FCMAccount = troncenv.String("FCM_SERVICE_ACCOUNT", "")
	if cfg.Porte.Issuer != "" {
		cfg.Porte.MobileRedirect = cmp.Or(cfg.Porte.MobileRedirect, DefaultMobileRedirect)
		cfg.CORSAllowedOrigins = withAppOrigins(cfg.CORSAllowedOrigins)
	}
	return cfg, cfg.validate()
}

func loadMatrix() Matrix {
	return Matrix{
		InternalURL: strings.TrimRight(troncenv.String("MATRIX_INTERNAL_URL", "http://127.0.0.1:8008"), "/"),
		PublicURL:   strings.TrimRight(troncenv.String("MATRIX_PUBLIC_URL", "http://localhost:8008"), "/"),
		ServerName:  troncenv.String("MATRIX_SERVER_NAME", "localhost"),
		ASToken:     troncenv.String("PRISME_AS_TOKEN", ""),
		HSToken:     troncenv.String("PRISME_HS_TOKEN", ""),
	}
}

func loadPorte() porte.Config {
	return porte.Config{
		Issuer:          troncenv.String("OIDC_ISSUER", ""),
		ClientID:        troncenv.String("OIDC_CLIENT_ID", ""),
		ClientSecret:    troncenv.String("OIDC_CLIENT_SECRET", ""),
		RedirectURL:     troncenv.String("OIDC_REDIRECT_URL", ""),
		SuccessURL:      troncenv.String("OIDC_SUCCESS_URL", ""),
		MachineAudience: troncenv.String("OIDC_MACHINE_AUDIENCE", ""),
		MobileRedirect:  troncenv.String("OIDC_MOBILE_REDIRECT", ""),
	}
}

// withAppOrigins adds the installed apps' origins to the CORS list, so a deploy
// that overrides CORS_ALLOWED_ORIGINS cannot lock the apps out.
func withAppOrigins(origins []string) []string {
	for _, origin := range appOrigins {
		if !slices.Contains(origins, origin) {
			origins = append(origins, origin)
		}
	}
	return origins
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
