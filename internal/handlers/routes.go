package handlers

import (
	"justTest/internal/infrastructure/auth"
	"justTest/internal/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(
	router *gin.Engine,
	authClient *auth.AuthClient,
	transactionHandler *TransactionHandler,
	accountHandler *AccountHandler,
	bankAccountHandler *BankAccountHandler,
	categoryHandler *CategoryHandler,
	budgetHandler *BudgetHandler,
	notificationHandler *NotificationHandler,
) {
	router.Use(middleware.CORSMiddleware())
	v1 := router.Group("/api/v1")
	public := v1.Group("")
	{
		public.GET("/health", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"status":  "OK",
				"service": "fin-core",
			})
		})
	}
	protected := v1.Group("")
	protected.Use(middleware.AuthMiddleware(authClient))
	{
		protected.POST("/account", accountHandler.CreateAccount)
		protected.GET("/account", accountHandler.GetAccount)

		bankAccounts := protected.Group("/bankAccounts")
		{
			bankAccounts.GET("", bankAccountHandler.GetBankAccounts)                 // все банк аккаунты дсотаются
			bankAccounts.POST("", bankAccountHandler.CreateBankAccount)              // просто создание
			bankAccounts.GET("/:bank_account_id", bankAccountHandler.GetBankAccount) // достается конкретный по банк аккаунт айди
			bankAccounts.DELETE("/:bank_account_id", bankAccountHandler.DeleteBankAccount)
			bankAccounts.PUT("/:bank_account_id/deactivate", bankAccountHandler.DeactivateBankAccount)
			bankAccounts.PUT("/:bank_account_id/activate", bankAccountHandler.ActivateBankAccount)
		}
		transactions := protected.Group("/transactions")
		{
			transactions.POST("", transactionHandler.CreateTransaction)
			transactions.GET("", transactionHandler.GetAllTransactions)
			transactions.GET("/:id", transactionHandler.GetTransaction)
			transactions.GET("/by-category/:category_id", transactionHandler.GetAllTransactionsByCategoryID)

		}
		protected.POST("/transfer", transactionHandler.TransferBetweenAccounts)

		protected.GET("/account/:account_id/transactions", transactionHandler.GetTransactionHistory) //  по сути удалить надо
		protected.GET("/bank_accounts/:account_id/balance", transactionHandler.GetBankAccountBalance)

		categories := protected.Group("/categories")
		{
			categories.POST("", categoryHandler.CreateCategory)
			categories.GET("", categoryHandler.GetByAccountID)
			categories.GET("/:category_id", categoryHandler.GetCategoryByID)
			categories.DELETE("/:category_id", categoryHandler.DeleteCategoryByID)
		}
		budgets := protected.Group("/budgets")
		{
			budgets.POST("", budgetHandler.CreateBudget)
			budgets.GET("", budgetHandler.GetBudgets)
			budgets.GET("/:category_id/status", budgetHandler.GetBudgetStatus)
			budgets.GET("/summary", budgetHandler.GetBudgetSummary)

		}
		notification := protected.Group("/notification")
		{

			notification.GET("", notificationHandler.GetUserNotifications)   // ?limit=20&offset=0
			notification.PUT("/:id/read", notificationHandler.MarkAsRead)    // PUT /notifications/123/read
			notification.PUT("/read-all", notificationHandler.MarkAllAsRead) // PUT /notifications/read-all
			notification.GET("/settings", notificationHandler.GetSettings)   // GET /notifications/settings
			notification.PUT("/settings", notificationHandler.SaveSettings)  // PUT /notifications/settings

		}
	}
	optional := v1.Group("/public")
	optional.Use(middleware.OptionalAuthMiddleware(authClient))
	{
		optional.GET("data", func(c *gin.Context) {
			userID, exists := c.Get("user_id")
			if !exists {
				c.JSON(200, gin.H{
					"message": "personal data",
					"user_id": userID,
				})
			} else {
				c.JSON(200, gin.H{
					"message": "public data",
				})
			}
		})
	}
}
