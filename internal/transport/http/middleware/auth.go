package middleware

import (
	"context"
	"crud/internal/services/user"
	httpapi "crud/internal/transport/http/helpers"
	"errors"
	"net/http"
)

type contextKey string

const userIDKey contextKey = "userID"

type AuthMiddleware struct {
	sessionStore user.SessionStore
}

func NewAuthMiddleware(sessionStore user.SessionStore) *AuthMiddleware {
	return &AuthMiddleware{
		sessionStore: sessionStore,
	}
}

func (s *AuthMiddleware) RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			if errors.Is(err, http.ErrNoCookie) {
				httpapi.WriteError(w, http.StatusUnauthorized, "missing session")
				return
			}
			httpapi.WriteError(w, http.StatusBadRequest, "invalid cookie")
			return
		}

		ctx := r.Context()
		sessionID := cookie.Value
		session, err := s.sessionStore.Get(ctx, sessionID)
		if err != nil {
			httpapi.WriteError(w, 401, "Not authorized")
			return
		}
		userID := session.UserID
		ctx = context.WithValue(ctx, userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func UserIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(userIDKey).(string)
	return id, ok
}
