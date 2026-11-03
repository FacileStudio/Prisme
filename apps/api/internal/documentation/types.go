// Package docs names the route-registry types each module fills in.
//
// They are aliases for tronc/apiref, so a module's documentation.go stays
// written against a local package while apiref owns the shape and renders it as
// OpenAPI behind the reference page at /docs.
package docs

import "github.com/FacileStudio/tronc/apiref"

// Response is the whole route registry.
type Response = apiref.Registry

// Module is one module's routes.
type Module = apiref.Module

// Route is one documented route.
type Route = apiref.Route

// Field is one documented body field.
type Field = apiref.Field

// Error is one documented error answer.
type Error = apiref.Error
