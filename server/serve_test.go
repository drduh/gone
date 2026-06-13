package server

import "testing"

// TestNewServer tests HTTP server setup.
func TestNewServer(t *testing.T) {
	app := newTestApp(nil)
	srv := newServer(app)

	if srv == nil {
		t.Fatal("new server is nil")
	}
	if srv.Handler == nil {
		t.Fatal("server handler is nil")
	}
	if srv.Addr != app.GetAddr() {
		t.Fatalf("addr = %q, want %q",
			srv.Addr, app.GetAddr())
	}
}
