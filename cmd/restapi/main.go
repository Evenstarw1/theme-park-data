package handler

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/evenstarw1/theme-park-data/pkg/themepark"
)

// Constantes para las variables de entorno
const (
	envDbConnectionString = "https://amsmpjefzegchbqafouw.supabase.co" 
	envListeningPort      = "THEMEPARK_LISTENING_PORT"       // Puerto en el que escuchará la app
)

// Config estructura para almacenar las configuraciones de entorno
type Config struct {
	databaseConnectionString string
	listeningPort            int
}

// Handler es la función exportada requerida por Vercel para manejar todas las solicitudes
func Handler(w http.ResponseWriter, r *http.Request) {
	config := getEnv()

	// Inicializa el router
	router := mux.NewRouter()

	// Conecta con la base de datos
	databaseStore, err := themepark.NewDatabaseStore(config.databaseConnectionString)
	if err != nil {
		log.Fatalf("Error al conectar con la base de datos: %v", err)
	}
	defer databaseStore.Close()

	// Crea los manejadores de rutas
	handlers := themepark.NewHandlers(databaseStore)

	// Rutas públicas
	pub := router.PathPrefix("/pub").Subrouter()
	pub.HandleFunc("/login", handlers.SignIn).Methods(http.MethodPost)
	pub.HandleFunc("/getCategories", handlers.GetCategories).Methods(http.MethodGet)
	pub.HandleFunc("/register", handlers.AddUser).Methods(http.MethodPost)

	// Rutas privadas
	priv := router.PathPrefix("/priv").Subrouter()
	priv.Use(handlers.AuthMiddleware)
	priv.HandleFunc("/", handlers.ServeHTTP).Methods(http.MethodGet)
	priv.HandleFunc("/users/{id:[0-9]+}", handlers.GetUser).Methods(http.MethodGet)
	priv.HandleFunc("/users/{id:[0-9]+}", handlers.UpdateUser).Methods(http.MethodPatch)
	priv.HandleFunc("/parks", handlers.GetParks).Methods(http.MethodGet)
	priv.HandleFunc("/parks", handlers.InsertThemePark).Methods(http.MethodPost)
	priv.HandleFunc("/park/{id:[0-9]+}", handlers.GetParkDetails).Methods(http.MethodGet)
	priv.HandleFunc("/park/{id:[0-9]+}", handlers.DeleteThemePark).Methods(http.MethodDelete)
	priv.HandleFunc("/park/{id:[0-9]+}", handlers.UpdateThemePark).Methods(http.MethodPatch)
	priv.HandleFunc("/park/comments", handlers.InsertParkComment).Methods(http.MethodPost)
	priv.HandleFunc("/categories", handlers.AddCategory).Methods(http.MethodPost)
	priv.HandleFunc("/users", handlers.GetAllUsers).Methods(http.MethodGet)

	// Manejador para rutas no encontradas
	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Ruta no encontrada: %s", r.URL.Path)
		http.Error(w, "404 - Ruta no encontrada", http.StatusNotFound)
	})

	// Redirige la solicitud al router
	router.ServeHTTP(w, r)
}

// getEnv obtiene las variables de entorno y las configura con valores por defecto si no existen
func getEnv() Config {
	var config Config

	// Cadena de conexión a la base de datos
	config.databaseConnectionString = os.Getenv(envDbConnectionString)
	if config.databaseConnectionString == "" {
		// Valor por defecto
		config.databaseConnectionString = "postgres://localhost:5432/themepark"
	}

	// Puerto de escucha
	portStr := os.Getenv(envListeningPort)
	if portStr == "" {
		// Valor por defecto
		config.listeningPort = 8080
	} else {
		// Convierte la cadena a entero
		config.listeningPort, _ = strconv.Atoi(portStr)
	}

	return config
}
