package pkg

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/coreos/go-oidc"
)

var (
	auth0Domain   = os.Getenv("AUTH0_DOMAIN")
	auth0Audience = os.Getenv("AUTH0_AUDIENCE")
)

type contextKey string

func AuthMiddleware(next http.Handler) http.Handler {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, "https://"+auth0Domain+"/")
	if err != nil {
		panic(err)
	}

	verifier := provider.Verifier(&oidc.Config{ClientID: auth0Audience})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Authorization header missing", http.StatusUnauthorized)
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Verify the token
		idToken, err := verifier.Verify(r.Context(), tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Extract custom claims or pass token down
		var claims struct {
			Email string `json:"email"`
			Sub   string `json:"sub"`
		}
		if err := idToken.Claims(&claims); err != nil {
			http.Error(w, "Failed to parse claims", http.StatusUnauthorized)
			return
		}
		// Add claims to context
		ctx := context.WithValue(r.Context(), contextKey("user"), claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
