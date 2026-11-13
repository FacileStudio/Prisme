package main

import (
	"context"
	"net"
	"net/http"
	"strconv"
	"strings"

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
// signed-in user, or on the address without a session. signIn guards porte's
// routes per address. service is the separate bucket of a service-key caller
// (Bulle), keyed on which key it presents, so a busy bot and the people behind
// one address never drain each other.
type budgets struct {
	visitor budget
	write   budget
	signIn  budget
	service budget
}

// productionBudgets sizes the buckets. Every /api route but porte's config
// reaches Synapse or the database, so visitor and write shape load per person.
// signIn is a strict bound of 30 per minute (15 + 15/60*60). service covers
// Bulle's provisioning pass, about three calls per person every ten minutes.
var productionBudgets = budgets{
	visitor: budget{limit: 240.0 / 60, burst: 120},
	write:   budget{limit: 30.0 / 60, burst: 30},
	signIn:  budget{limit: 15.0 / 60, burst: 15},
	service: budget{limit: 1200.0 / 60, burst: 600},
}

// callers tells the limiter who sends a request: a service key, or the user
// whose session it carries. auth.Service is one.
type callers interface {
	ServiceCaller(r *http.Request) int
	SessionUser(r *http.Request) (int64, bool)
}

// visitorKeyContext carries the visitor's bucket key from the one session
// lookup per request to both limiters that spend on it.
type visitorKeyContext struct{}

// signInPath is where porte's routes live. They stay counted per address.
const signInPath = "/api/auth/"

// rateLimits are the two middlewares mountAPI applies: requests on the whole
// /api group, signIn on porte's routes as well.
type rateLimits struct {
	requests func(http.Handler) http.Handler
	signIn   func(http.Handler) http.Handler
}

// newRateLimits builds the limiters. A signed-in user's requests spend from
// the user's own buckets, wherever they come from, so people behind one office
// address never drain each other. Everything else spends from its address's.
func newRateLimits(b budgets, who callers) rateLimits {
	visitor := limiter(b.visitor, visitorKey)
	write := limiter(b.write, visitorKey)
	service := limiter(b.service, func(r *http.Request) string {
		return "service-key:" + strconv.Itoa(who.ServiceCaller(r))
	})
	requests := func(next http.Handler) http.Handler {
		byService, byReader, byWriter := service(next), visitor(next), visitor(write(next))
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if who.ServiceCaller(r) != 0 {
				byService.ServeHTTP(w, r)
				return
			}
			r = r.WithContext(context.WithValue(r.Context(), visitorKeyContext{}, visitorOf(who, r)))
			if isRead(r.Method) {
				byReader.ServeHTTP(w, r)
				return
			}
			byWriter.ServeHTTP(w, r)
		})
	}
	return rateLimits{requests: requests, signIn: limiter(b.signIn, addressKey)}
}

// visitorOf is the bucket key of a request that presents no service key: the
// user its session resolves to, or its address. porte's routes always take
// the address, since they are where a session starts or ends, and they skip
// the session lookup with it.
func visitorOf(who callers, r *http.Request) string {
	if !strings.HasPrefix(r.URL.Path, signInPath) {
		if userID, ok := who.SessionUser(r); ok {
			return "user:" + strconv.FormatInt(userID, 10)
		}
	}
	return addressKey(r)
}

// visitorKey reads the key visitorOf put on the request.
func visitorKey(r *http.Request) string {
	if key, ok := r.Context().Value(visitorKeyContext{}).(string); ok {
		return key
	}
	return addressKey(r)
}

// addressKey is the host of RemoteAddr as httpx.NewRouter's RealIP rewrote it
// from TRUSTED_PROXIES, so behind Cloudflare it is Cf-Connecting-Ip read from
// a trusted chain, never a header the caller wrote.
func addressKey(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		host = r.RemoteAddr
	}
	return "ip:" + host
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
