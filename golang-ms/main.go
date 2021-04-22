package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/iqnev/golang-rest/ms/data"
	"github.com/iqnev/golang-rest/ms/handlers"
	"google.golang.org/grpc"

	"github.com/go-openapi/runtime/middleware"

	gohandlers "github.com/gorilla/handlers"
	protos "github.com/iqnev/golang-rest/currency/protos/currency"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	//	defer profile.Start(profile.MemProfile, profile.ProfilePath(".")).Stop()

	l := hclog.Default()

	err := godotenv.Load(".env")

	if err != nil {
		l.Error("Error loading .env file")
		os.Exit(1)
	}

	appHost := os.Getenv("APP_HOST")
	appPort := os.Getenv("APP_PORT")

	v := data.NewValidation()

	conn, err := grpc.Dial("localhost:8989", grpc.WithInsecure())
	if err != nil {
		panic(err)
	}

	defer conn.Close()

	cc := protos.NewCurrencyClient(conn)

	db := data.NewProductDB(cc, l)

	ph := handlers.NewProducts(l, v, db)

	sm := mux.NewRouter()

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", ph.ListAll)
	getRouter.HandleFunc("/", ph.ListAll).Queries("currency", "{[A-Z]{3}}")
	getRouter.HandleFunc("/products/{id:[0-9]+}", ph.ListSingle)
	getRouter.HandleFunc("/products/{id:[0-9]+}", ph.ListSingle).Queries("currency", "{[A-Z]{3}}")

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
		l.Info("Starting server on port %s", appPort)
		err := serv.ListenAndServe()

		if err != nil {
			l.Error("Error starting server", "error", err)
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
