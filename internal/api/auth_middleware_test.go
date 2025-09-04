package api

import "testing"

func TestAuthenticate(t *testing.T) {
	a := AuthMiddleware{}
	if !a.Authenticate("token") {
		t.Error("expected valid token to pass")
	}
	if a.Authenticate("") {
		t.Error("expected empty token to fail")
	}
}
