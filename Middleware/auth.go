package Middleware

import (
	"context"
	"database/sql"
	"github.com/golang-jwt/jwt/v5"
	"github.com/pkg/errors"
	"net/http"
	"os"
	"strings"
	"todo/Database/dbHelper"
	"todo/Models"
	"todo/Utils"
)

type ContextKey string

const (
	userContext ContextKey = "userContext"
)

func Authentication(n http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			Utils.RespondError(w, http.StatusUnauthorized, nil, "authorization header missing")
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader {
			Utils.RespondError(w, http.StatusUnauthorized, nil, "bearer token missing")
			return
		}

		token, parseErr := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("invalid signing method") // Invalid signing method error
			}
			return []byte(os.Getenv("JWT_SECRET_KEY")), nil
		})

		if parseErr != nil || !token.Valid {
			Utils.RespondError(w, http.StatusUnauthorized, parseErr, "invalid token")
			return
		}

		claimValues, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			Utils.RespondError(w, http.StatusUnauthorized, nil, "invalid token claims")
			return
		}

		sessionID := claimValues["sessionID"].(string)

		archivedAt, err := dbHelper.GetArchivedAt(sessionID)
		if err != nil && !errors.Is(err, sql.ErrNoRows) {
			Utils.RespondError(w, http.StatusInternalServerError, err, "internal server error")
			return
		}

		if archivedAt != nil {
			Utils.RespondError(w, http.StatusUnauthorized, nil, "invalid token")
			return
		}

		user := &Models.UserCtx{
			UserID:    claimValues["userID"].(string),
			Name:      claimValues["name"].(string),
			Email:     claimValues["email"].(string),
			SessionID: sessionID,
		}

		ctx := context.WithValue(r.Context(), userContext, user)
		r = r.WithContext(ctx)

		n.ServeHTTP(w, r)
	})
}

func UserContext(r *http.Request) *Models.UserCtx {
	if user, ok := r.Context().Value(userContext).(*Models.UserCtx); ok {
		return user
	}
	return nil
}
