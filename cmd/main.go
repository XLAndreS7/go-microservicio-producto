package main

import (
	"context"
	"fmt"
	"go-microservicio-producto/controllers"
	"go-microservicio-producto/repository"
	"go-microservicio-producto/services"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	// Cargar archivo .env
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  Advertencia: No se pudo cargar el archivo .env, se usará configuración por defecto")
	}

	// Leer variables de entorno
	port := os.Getenv("PORT")
	if port == "" {
		port = "8084" // Valor por defecto
	}

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		log.Fatal("❌ Error: La variable de entorno MONGO_URI no está definida")
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		log.Fatal("❌ Error: La variable de entorno DB_NAME no está definida")
	}

	// Conexión a MongoDB
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		log.Fatalf("❌ Error creando cliente de MongoDB: %v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		log.Fatalf("❌ Error conectando a MongoDB: %v", err)
	}

	// Verificación de conexión
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("❌ No se pudo conectar con MongoDB: %v", err)
	}
	fmt.Println("✅ Conectado a MongoDB correctamente")

	db := client.Database(dbName)

	// Inyección de dependencias
	productoRepo := repository.NewProductoRepository(db)
	productoService := services.NewProductoService(productoRepo)
	productoController := controllers.NewProductoController(productoService)

	// Usar gorilla/mux para rutas dinámicas
	router := mux.NewRouter()

	// Rutas RESTful
// Rutas RESTful
router.HandleFunc("/productos", productoController.CrearProducto).Methods("POST")
router.HandleFunc("/productos", productoController.ObtenerProductos).Methods("GET")
router.HandleFunc("/productos/{id}", productoController.ObtenerProductoPorID).Methods("GET") // <- esta es la nueva
router.HandleFunc("/productos/{id}", productoController.ActualizarProducto).Methods("PUT")
         router.HandleFunc("/productos/{id}", productoController.EliminarProducto).Methods("DELETE")

	// Iniciar servidor
	fmt.Printf("🚀 Servidor corriendo en http://localhost:%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, router))
}

