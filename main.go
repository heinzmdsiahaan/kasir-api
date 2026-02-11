package main

import (
	"encoding/json"
	"fmt"
	"kasir-api/database"
	"kasir-api/handlers"
	"kasir-api/middlewares"
	"kasir-api/repositories"
	"kasir-api/services"
	"net/http"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type Config struct {
	Port    string `mapstructure:"PORT"`
	DB_CONN string `mapstructure:"DB_CONN"`
	APIKey  string `mapstructure:"API_KEY"`
}

func main() {
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if _, err := os.Stat(".env"); err == nil {
		viper.SetConfigFile(".env")
		_ = viper.ReadInConfig()
	}

	config := Config{
		Port:    viper.GetString("PORT"),
		DB_CONN: viper.GetString("DB_CONN"),
		APIKey:  viper.GetString("API_KEY"),
	}

	// Setup Database
	db, err := database.InitDB(config.DB_CONN)
	if err != nil {
		fmt.Println("Failed to connect to database:", err)
		return
	}
	defer db.Close()

	apiKeyMiddleware := middlewares.APIKey(config.APIKey)

	// Initialize Layers
	productRepo := repositories.NewProductRepository(db)
	productService := services.NewProductService(productRepo)
	productHandler := handlers.NewProductHandler(productService)

	categoryRepo := repositories.NewCategoryRepository(db)
	categoryService := services.NewCategoryService(categoryRepo)
	categoryHandler := handlers.NewCategoryHandler(categoryService)

	// Setup routes
	http.HandleFunc("/api/produk", middlewares.CORS(middlewares.Logger(productHandler.HandleProducts)))
	http.HandleFunc("/api/produk/", middlewares.CORS(middlewares.Logger(apiKeyMiddleware(productHandler.HandleProductByID))))

	http.HandleFunc("/api/product", middlewares.CORS(middlewares.Logger(productHandler.HandleProducts)))
	http.HandleFunc("/api/product/", middlewares.CORS(middlewares.Logger(apiKeyMiddleware(productHandler.HandleProductByID))))

	http.HandleFunc("/api/category", middlewares.CORS(middlewares.Logger(categoryHandler.HandleCategories)))
	http.HandleFunc("/api/category/", middlewares.CORS(middlewares.Logger(apiKeyMiddleware(categoryHandler.HandleCategoryByID))))

	transactionRepo := repositories.NewTransactionRepository(db)
	transactionService := services.NewTransactionService(transactionRepo)
	transactionHandler := handlers.NewTransactionHandler(transactionService)

	http.HandleFunc("/api/checkout", middlewares.CORS(middlewares.Logger(apiKeyMiddleware(transactionHandler.Checkout))))

	reportRepo := repositories.NewReportRepository(db)
	reportService := services.NewReportService(reportRepo)
	reportHandler := handlers.NewReportHandler(reportService)

	http.HandleFunc("/api/report", reportHandler.HandleReport)

	// GET localhost:8080/api/health
	http.HandleFunc("/api/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"status":  "OK",
			"message": "API Running",
		})
	})

	addr := "0.0.0.0:" + config.Port
	fmt.Println("Starting server on " + addr)

	err = http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Println("Error starting server:", err)
	}
}
