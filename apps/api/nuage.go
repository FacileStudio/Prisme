package main

import (
	"context"
	"log/slog"

	"github.com/FacileStudio/Prisme/apps/api/internal/env"
	nuageclient "github.com/FacileStudio/Prisme/apps/api/internal/nuage"
	"github.com/FacileStudio/Prisme/apps/api/modules/auth"
	"github.com/FacileStudio/Prisme/apps/api/modules/nuage"

	"github.com/FacileStudio/porte/exchange"
)

// nuageAudience is Nuage's Registre client id, the audience Registre's seed
// lets Prisme exchange a user's token for.
const nuageAudience = "9U40zs8fnTkpKpR6MMobIOXk0OkvotC7vd1NECGH"

var nuageTarget = exchange.Target{Audience: nuageAudience, Scopes: []string{"nuage.read", "nuage.write"}}

// withNuage adds the Nuage picker to the suite apps, once porte exists to hold
// the users' sealed Registre tokens.
func withNuage(apps suiteApps, cfg env.Config, authService *auth.Service, appLogger *slog.Logger) suiteApps {
	apps.nuage = newNuagePicker(cfg, authService, appLogger)
	return apps
}

// newNuagePicker builds the Nuage picker over a Registre token exchange on
// the identity store that seals the users' Registre tokens. Without OIDC
// there is no user grant to exchange, so the picker answers unavailable.
func newNuagePicker(cfg env.Config, authService *auth.Service, appLogger *slog.Logger) *nuage.Service {
	porteCfg := cfg.Porte
	if porteCfg.Issuer == "" || porteCfg.ClientID == "" || porteCfg.ClientSecret == "" || cfg.NuageURL == "" {
		return nuage.NewService(nil, appLogger)
	}
	exchanger, err := exchange.New(exchange.Config{
		Issuer: porteCfg.Issuer, ClientID: porteCfg.ClientID, ClientSecret: porteCfg.ClientSecret,
		Identities: authService.Identities(), Logger: appLogger,
	})
	if err != nil {
		appLogger.Error("the nuage picker is off: the token exchange is misconfigured", slog.Any("error", err))
		return nuage.NewService(nil, appLogger)
	}
	client := nuageclient.New(cfg.NuageURL, nuageclient.Tokens{
		Token: func(ctx context.Context, userID int64) (string, error) {
			return exchanger.Token(ctx, userID, nuageTarget)
		},
		Forget: func(userID int64) { exchanger.Forget(userID, nuageTarget) },
	})
	return nuage.NewService(client, appLogger)
}
