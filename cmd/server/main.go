// @title Fin-Core API
// @version 1.0
// @description Personal Finance Management API
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.url http://www.swagger.io/support
// @contact.email support@swagger.io

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8080
// @BasePath /api/v1
// @schemes http

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

package main

import (
	"database/sql"
	"fmt"
	"justTest/internal/events"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "justTest/docs"
	"justTest/internal/handlers"
	"justTest/internal/infrastructure/auth"
	"justTest/internal/repo"
	"justTest/internal/services"
	"log"
	"os"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}
	db, err := initDB()
	if err != nil {
		log.Fatal("Error connecting to database")
	}
	defer db.Close()
	authClient := auth.NewAuthClient(os.Getenv("AUTH_SERVICE_URL"))
	accountRepo := repo.NewAccountRepository(db)
	bankAccountRepo := repo.NewBankAccountRepository(db)
	transactionRepo := repo.NewTransactionRepository(db)
	categoryRepo := repo.NewCategoryRepository(db)
	budgetRepo := repo.NewBudgetRepository(db)
	notificationRepo := repo.NewNotificationRepository(db)
	settingRepo := repo.NewUserNotificationSettingsRepository(db)

	publisher, err := events.NewPublisher("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatal(err)
	}
	defer publisher.Close()
	accountService := services.NewAccountService(accountRepo, bankAccountRepo, transactionRepo, authClient)
	bankAccService := services.NewBankAccService(bankAccountRepo, accountRepo)
	transactionService := services.NewTransactionService(transactionRepo, bankAccountRepo, categoryRepo, accountRepo)
	categoryService := services.NewCategoryService(accountRepo, categoryRepo, authClient)
	budgetService := services.NewBudgetService(budgetRepo, transactionRepo, accountRepo, categoryRepo, publisher)
	notificationService := services.NewNotificationService(notificationRepo, settingRepo, publisher)

	consumer, err := events.NewConsumer("amqp://guest:guest@localhost:5672/", notificationService, budgetService)
	if err != nil {
		log.Fatal("Failed to create consumer:", err)
	}
	defer consumer.Close()
	accountHandler := handlers.NewAccountHandler(accountService)
	bankAccountHandler := handlers.NewBankAccountHandler(bankAccService)
	transactionHandler := handlers.NewTransactionHandler(transactionService, publisher)
	categoryHandler := handlers.NewCategoryHandler(categoryService)
	budgetHandler := handlers.NewBudgetHandler(budgetService)
	notificationHandler := handlers.NewNotificationHandler(notificationService)

	router := gin.Default()

	handlers.RegisterRoutes(
		router,
		authClient,
		transactionHandler,
		accountHandler,
		bankAccountHandler,
		categoryHandler,
		budgetHandler,
		notificationHandler,
	)

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	go func() {
		log.Println("Starting consumer")
		consumer.ConsumeAll()
	}()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
func initDB() (*sql.DB, error) {
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")
	sslMode := os.Getenv("DB_SSL_MODE")
	if dbHost == "" {
		dbHost = "localhost"
	}
	if dbPort == "" {
		dbPort = "5432"
	}
	if dbUser == "" {
		dbUser = "postgres"
	}
	if dbPassword == "" {
		dbPassword = "postgres"
	}
	if dbName == "" {
		dbName = "fin_db"
	}
	if sslMode == "" {
		sslMode = "disable"
	}
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbPassword, dbName, sslMode)

	log.Printf("Connecting to database: host=%s port=%s user=%s dbname=%s sslmode=%s",
		dbHost, dbPort, dbUser, dbName, sslMode)

	db, err := sql.Open("postgres", connStr)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetConnMaxIdleTime(1 * time.Minute)

	if err != nil {
		log.Printf("Failed to open database connection: %v", err)
		return nil, fmt.Errorf("failed to open database: %v", err)
	}

	log.Println("Database connection opened, testing ping...")
	err = db.Ping()
	if err != nil {
		log.Printf("Failed to ping database: %v", err)
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}
	log.Println("Database connection established successfully")

	return db, nil
}
