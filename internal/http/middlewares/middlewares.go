package middlewares

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/TravellerGSF/grpc_distr_calc/internal/storage"
	"github.com/TravellerGSF/grpc_distr_calc/internal/utils/orchestrator/jwts"
)

func AuthorizeJWTToken(next http.Handler, db storage.UserInteractor) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("auth_token")
		if err != nil {
			http.Redirect(w, r, "/auth", http.StatusSeeOther)
			log.Printf("no cookie found")
			return
		}
		tokenString := cookie.Value
		tokenValue, err := jwts.VerifyJWTToken(tokenString)
		if err != nil {
			http.Redirect(w, r, "/auth", http.StatusSeeOther)
			log.Printf("error: %v", err)
			return
		}
		userID, err := strconv.ParseInt(tokenValue, 10, 64)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Printf("error: %v", err)
			return
		}

		ctx := context.WithValue(r.Context(), "userid", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUsernameHandler(db storage.UserInteractor) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		tokenString := authHeader[len("Bearer "):]
		tokenValue, err := jwts.VerifyJWTToken(tokenString)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		userID, err := strconv.ParseInt(tokenValue, 10, 64)
		if err != nil {
			http.Error(w, "Invalid user ID in token", http.StatusUnauthorized)
			return
		}

		user, err := db.GetUserByID(userID)
		if err != nil {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"username": user.Name})
	}
}
