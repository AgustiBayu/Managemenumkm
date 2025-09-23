package middleware

import (
	"Managemenumkm/domain"
	"Managemenumkm/helper"
	"context"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// Define a key for the context value
type ContextKey string

var ClaimsKey = ContextKey("claims")

// Authorize is a middleware that checks for a valid JWT and authorizes based on roles.
func Authorize(h httprouter.Handle, requiredRoles ...domain.Role) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		claims, err := helper.ValidateJWT(cookie.Value)
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		// Check if the user's role is one of the required roles
		isAuthorized := false
		for _, requiredRole := range requiredRoles {
			if string(requiredRole) == claims.Role {
				isAuthorized = true
				break
			}
		}

		if !isAuthorized {
			http.Error(w, "Forbidden: You don't have access to this page.", http.StatusForbidden)
			return
		}

		// Add claims to the request context to be used by handlers
		ctx := context.WithValue(r.Context(), ClaimsKey, claims)
		h(w, r.WithContext(ctx), ps)
	}
}
