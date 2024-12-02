package main

import (
	"fmt"
	"github.com/gorilla/mux"
	"github.com/evenstarw1/theme-park-data/pkg/themepark"
	"net/http"
	"os"
	"strconv"
)

const (
	envDbConnectionString = "THEMEPARK_DB_CONNECTION_STRING"
	envListeningPort      = "THEMEPARK_LISTENING_PORT"
)

type Config struct {
	databaseConnectionString string
	listeningPort            int
}

func main() {
	config := getEnv()

	router := mux.NewRouter()
	databaseStore, err := themepark.NewDatabaseStore(config.databaseConnectionString)
	if err != nil {
		os.Exit(1)
	}

	defer databaseStore.Close()

	handlers := themepark.NewHandlers(databaseStore)

	pub := router.PathPrefix("/pub").Subrouter()
	pub.HandleFunc("/login", handlers.SignIn).
		Methods(http.MethodPost)
	pub.HandleFunc("/getCategories", handlers.GetCategories).
		Methods(http.MethodGet)
	pub.HandleFunc("/register", handlers.AddUser).
		Methods(http.MethodPost)

	priv := router.PathPrefix("/priv").Subrouter()
	priv.Use(handlers.AuthMiddleware)
	priv.HandleFunc("/", handlers.ServeHTTP).
		Methods(http.MethodGet)
	priv.HandleFunc("/users/{id:[0-9]+}", handlers.GetUser).
		Methods(http.MethodGet)
	priv.HandleFunc("/users/{id:[0-9]+}", handlers.UpdateUser).
		Methods(http.MethodPatch)

	http.ListenAndServe(fmt.Sprintf(":%d", config.listeningPort), router)
}

func getEnv() Config {
	var config Config
	config.databaseConnectionString = os.Getenv(envDbConnectionString)
	if config.databaseConnectionString == "" {
		config.databaseConnectionString = "postgres://localhost:5432/themepark"
	}

	portStr := os.Getenv(envListeningPort)
	if portStr == "" {
		config.listeningPort = 8080
	} else {
		config.listeningPort, _ = strconv.Atoi(portStr) // Not so production ready lol
	}

	return config
}
