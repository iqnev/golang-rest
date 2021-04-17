package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/iqnev/golang-rest/data"
	"github.com/iqnev/golang-rest/handlers"

	"github.com/go-openapi/runtime/middleware"

	gohandlers "github.com/gorilla/handlers"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	//	defer profile.Start(profile.MemProfile, profile.ProfilePath(".")).Stop()

	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	err := godotenv.Load(".env")

	if err != nil {
		l.Fatalf("Error loading .env file")
		os.Exit(1)
	}

	appHost := os.Getenv("APP_HOST")
	appPort := os.Getenv("APP_PORT")

	v := data.NewValidation()

	ph := handlers.NewProducts(l, v)

	sm := mux.NewRouter()

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", ph.ListAll)
	getRouter.HandleFunc("/products/{id:[0-9]+}", ph.ListSingle)

	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.Update)
	putRouter.Use(ph.MiddlewareValidateProduct)

	postRouter := sm.Methods(http.MethodPost).Subrouter()
	postRouter.HandleFunc("/", ph.Create)
	postRouter.Use(ph.MiddlewareValidateProduct)

	deleteRouter := sm.Methods(http.MethodDelete).Subrouter()
	deleteRouter.HandleFunc("/products/{id:[0-9]+}", ph.Delete)

	opts := middleware.RedocOpts{SpecURL: "/swagger.yaml"}
	sh := middleware.Redoc(opts, nil)

	getRouter.Handle("/docs", sh)
	getRouter.Handle("/swagger.yaml", http.FileServer(http.Dir("./")))

	cor := gohandlers.CORS(gohandlers.AllowedOrigins([]string{"*"}))

	serv := &http.Server{
		Addr:         appHost + ":" + appPort,
		Handler:      cor(sm),
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		l.Printf("Starting server on port %s", appPort)
		err := serv.ListenAndServe()

		if err != nil {
			l.Printf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	sigChan := make(chan os.Signal)

	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan

	log.Println("Got signal:", sig)

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)

	serv.Shutdown(ctx)

}
