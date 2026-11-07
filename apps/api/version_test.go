package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
)

func TestVersionMatchesPackageJSON(t *testing.T) {
	if version == "dev" {
		t.Skip("built without -X version; scripts/check.sh links it from package.json")
	}
	raw, err := os.ReadFile("../client/package.json")
	if err != nil {
		t.Fatal(err)
	}
	var pkg struct {
		Version string `json:"version"`
	}
	if err := json.Unmarshal(raw, &pkg); err != nil {
		t.Fatal(err)
	}
	want := "v" + pkg.Version
	rec := httptest.NewRecorder()
	versionHandler(rec, httptest.NewRequest(http.MethodGet, "/version", nil))
	if !strings.Contains(rec.Body.String(), `"version":"`+want+`"`) {
		t.Fatalf("/version = %s, package.json says %s", rec.Body.String(), want)
	}
}
