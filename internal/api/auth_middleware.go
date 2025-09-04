package api

// AuthMiddleware performs simple token checks.
type AuthMiddleware struct{}

// Authenticate validates a token.
func (a *AuthMiddleware) Authenticate(token string) bool {
	return token != ""
}
