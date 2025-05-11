package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	os "os"

	authHandler "github.com/TravellerGSF/grpc_distr_calc/internal/http/handlers/auth"
	exprHandler "github.com/TravellerGSF/grpc_distr_calc/internal/http/handlers/expression"
	"github.com/TravellerGSF/grpc_distr_calc/internal/http/middlewares"
	"github.com/TravellerGSF/grpc_distr_calc/internal/storage"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	ctx := context.Background()
	db, err := storage.New("./db/storage.db")
	if err != nil {
		panic(err)
	}
	mux := http.NewServeMux()
	mainPageHandler := middlewares.AuthorizeJWTToken(http.FileServer(http.Dir("frontend/main")), db)
	authPageHandler := http.StripPrefix("/auth", http.FileServer(http.Dir("frontend/auth")))
	mux.Handle("/", mainPageHandler)
	mux.Handle("/auth/", authPageHandler)
	mux.HandleFunc("/auth/signup/", authHandler.RegisterUserHandler(ctx, db))
	mux.HandleFunc("/auth/login/", authHandler.LoginUserHandler(ctx, db))

	mux.Handle("/expression/", middlewares.AuthorizeJWTToken(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			exprHandler.CreateExpressionHandler(ctx, db).ServeHTTP(w, r)
		case http.MethodGet:
			exprHandler.GetExpressionsHandler(ctx, db).ServeHTTP(w, r)
		case http.MethodDelete:
			exprHandler.DeleteExpressionHandler(ctx, db).ServeHTTP(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}), db))

	host, ok := os.LookupEnv("ORCHESTRATOR_HOST")
	if !ok {
		log.Print("ORCHESTRATOR_HOST not set, using 0.0.0.0")
		host = "0.0.0.0"
	}
	port, ok := os.LookupEnv("ORCHESTRATOR_PORT")
	if !ok {
		log.Print("ORCHESTRATOR_PORT not set, using 8080")
		port = "8080"
	}
	addr := fmt.Sprintf("%s:%s", host, port)
	server := &http.Server{
		Addr:    addr,
		Handler: mux,
	}
	log.Printf("running Orchestrator server at %s", addr)
	go log.Fatal(server.ListenAndServe())
	log.Print("Something went wrong...")
}
