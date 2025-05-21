package middlewares

import (
	"context"
	"log"
	"time"

	"github.com/gofiber/fiber/v2"
	"mobilka/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DBHealthCheckMiddleware creates middleware that ensures database connection is healthy
func DBHealthCheckMiddleware(dbManager *utils.DBManager) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Check if DB connection is healthy
		reconnected, err := dbManager.EnsureConnected(ctx)
		if err != nil {
			log.Printf("Database health check failed: %v", err)
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status":  "error",
				"message": "Database connection unavailable",
			})
		}

		if reconnected {
			log.Println("Database connection reestablished during health check")
		}

		return c.Next()
	}
}

// WithDBRetry wraps a handler with database retry logic
func WithDBRetry(dbManager *utils.DBManager, handler func(c *fiber.Ctx) error) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Set a reasonable timeout for the entire operation
		ctx, cancel := context.WithTimeout(c.Context(), 30*time.Second)
		defer cancel()

		var handlerErr error
		operation := string(c.Request().URI().Path())

		// Adjust the function signature to match ExecuteWithRetry's expected type
		err := dbManager.ExecuteWithRetry(ctx, operation, func(ctx context.Context, pool *pgxpool.Pool) error {
			// Call the original handler
			handlerErr = handler(c)
			return handlerErr
		})

		// Check if response has been sent using Fiber's context state
		if err != nil && c.Response().StatusCode() == 0 {
			return c.Status(fiber.StatusServiceUnavailable).JSON(fiber.Map{
				"status":  "error",
				"message": "Database operation failed after retries",
			})
		}

		return handlerErr
	}
}
