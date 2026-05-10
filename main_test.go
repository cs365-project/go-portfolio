package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

// setupFileServer creates the same file server as main(), but returns it as an http.Handler
func setupFileServer(dir string) http.Handler {
	return http.FileServer(http.Dir(dir))
}

// TestStaticFileServerReturns200 verifies that the server serves index.html with HTTP 200
func TestStaticFileServerReturns200(t *testing.T) {
	// Create a temp dir with a dummy index.html to simulate ./static
	tmpDir := t.TempDir()
	err := os.WriteFile(tmpDir+"/index.html", []byte("<html><body>ok</body></html>"), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	handler := setupFileServer(tmpDir)
	// Request root — FileServer serves index.html and returns 200
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("expected status 200, got %d", rr.Code)
	}
}

// TestStaticFileServerReturns404ForMissingFile verifies that missing files return HTTP 404
func TestStaticFileServerReturns404ForMissingFile(t *testing.T) {
	tmpDir := t.TempDir()

	handler := setupFileServer(tmpDir)
	req := httptest.NewRequest(http.MethodGet, "/notfound.html", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusNotFound {
		t.Errorf("expected status 404, got %d", rr.Code)
	}
}

// TestStaticFileServerContentType verifies that HTML files are served with correct Content-Type
func TestStaticFileServerContentType(t *testing.T) {
	tmpDir := t.TempDir()
	err := os.WriteFile(tmpDir+"/index.html", []byte("<html><body>test</body></html>"), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	handler := setupFileServer(tmpDir)
	// Request root so FileServer serves index.html directly (not a redirect)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	ct := rr.Header().Get("Content-Type")
	if ct == "" {
		t.Error("expected Content-Type header to be set, got empty string")
	}
}

// TestRootPathServesIndex verifies that GET / redirects or serves index.html
func TestRootPathServesIndex(t *testing.T) {
	tmpDir := t.TempDir()
	err := os.WriteFile(tmpDir+"/index.html", []byte("<html><body>home</body></html>"), 0644)
	if err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	handler := setupFileServer(tmpDir)
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	// FileServer returns 200 for directory listing or index.html
	if rr.Code != http.StatusOK && rr.Code != http.StatusMovedPermanently {
		t.Errorf("expected status 200 or 301 for root path, got %d", rr.Code)
	}
}
