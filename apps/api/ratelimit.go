package main

import (
	"net/http"
	"strconv"

	"github.com/FacileStudio/tronc/errors"
	"github.com/FacileStudio/tronc/httpjson"
	troncmiddleware "github.com/FacileStudio/tronc/middleware"
)

// budget is one token bucket: limit tokens per second refill a bucket of
// burst. Over any span of W seconds it admits at most burst + limit*W.
type budget struct {
	limit float64
	burst int
}

// budgets are the buckets in front of /api. visitor spends on every browser
// request and write on top of it for anything but a read, both keyed on the
// visitor's address. signIn guards porte's routes. service is the separate
// bucket of a service-key caller (Bulle), keyed on which key it presents, so a
// busy bot and the people behind one address never drain each other.
type budgets struct {
	visitor budget
	write   budget
	signIn  budget
	service budget
}

// productionBudgets sizes the buckets. Every /api route but porte's config
// reaches Synapse or the database, so visitor and write shape load: an office
// behind one address still fits. signIn is a strict bound of 30 per minute
// (15 + 15/60*60). service covers Bulle's provisioning pass, about three calls
// per person every ten minutes.
var productionBudgets = budgets{
	visitor: budget{limit: 240.0 / 60, burst: 120},
	write:   budget{limit: 30.0 / 60, burst: 30},
	signIn:  budget{limit: 15.0 / 60, burst: 15},
	service: budget{limit: 1200.0 / 60, burst: 600},
}

// serviceCallers tells a service-key request apart. auth.Service is one.
type serviceCallers interface {
	ServiceCaller(r *http.Request) int
}

// rateLimits are the two middlewares mountAPI applies: requests on the whole
// /api group, signIn on porte's routes as well.
type rateLimits struct {
	requests func(http.Handler) http.Handler
	signIn   func(http.Handler) http.Handler
}

// newRateLimits builds the limiters. A visitor's key is the host of RemoteAddr
// as httpx.NewRouter's RealIP rewrote it from TRUSTED_PROXIES, so behind
// Cloudflare it is Cf-Connecting-Ip read from a trusted chain, never a header
// the caller wrote.
func newRateLimits(b budgets, callers serviceCallers) rateLimits {
	visitor := limiter(b.visitor, nil)
	write := limiter(b.write, nil)
	service := limiter(b.service, func(r *http.Request) string {
		return "service-key:" + strconv.Itoa(callers.ServiceCaller(r))
	})
	requests := func(next http.Handler) http.Handler {
		byService, byReader, byWriter := service(next), visitor(next), visitor(write(next))
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch {
			case callers.ServiceCaller(r) != 0:
				byService.ServeHTTP(w, r)
			case isRead(r.Method):
				byReader.ServeHTTP(w, r)
			default:
				byWriter.ServeHTTP(w, r)
			}
		})
	}
	return rateLimits{requests: requests, signIn: limiter(b.signIn, nil)}
}

func limiter(b budget, key func(*http.Request) string) func(http.Handler) http.Handler {
	return troncmiddleware.Ratelimit(troncmiddleware.RatelimitConfig{
		Limit: b.limit, Burst: b.burst, KeyFunc: key, MaxKeys: 50000, OnRefuse: refuseRate,
	})
}

func refuseRate(w http.ResponseWriter, _ *http.Request) {
	httpjson.WriteError(w, errors.RateLimited("too many requests, try again later"))
}

func isRead(method string) bool {
	return method == http.MethodGet || method == http.MethodHead || method == http.MethodOptions
}
