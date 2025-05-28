package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/chizheg/forum/proto"
	"google.golang.org/grpc"
)

type AuthMiddleware struct {
	authClient proto.AuthServiceClient
}

func NewAuthMiddleware(authConn *grpc.ClientConn) *AuthMiddleware {
	return &AuthMiddleware{
		authClient: proto.NewAuthServiceClient(authConn),
	}
}

func (m *AuthMiddleware) Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "missing authorization header", http.StatusUnauthorized)
			return
		}

		// Extract token from Bearer header
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, "invalid authorization header", http.StatusUnauthorized)
			return
		}

		// Validate token with auth service
		resp, err := m.authClient.ValidateToken(r.Context(), &proto.ValidateTokenRequest{
			Token: parts[1],
		})

		if err != nil || !resp.Valid {
			http.Error(w, "invalid token", http.StatusUnauthorized)
			return
		}

		// Add user ID to context
		ctx := context.WithValue(r.Context(), "userID", int(resp.UserId))
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func (m *AuthMiddleware) OptionalAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			next.ServeHTTP(w, r)
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			next.ServeHTTP(w, r)
			return
		}

		resp, err := m.authClient.ValidateToken(r.Context(), &proto.ValidateTokenRequest{
			Token: parts[1],
		})

		if err != nil || !resp.Valid {
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), "userID", int(resp.UserId))
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}
