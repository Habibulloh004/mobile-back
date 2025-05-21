package utils

import (
	"context"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

// DBManager handles database connection management with reconnection logic
type DBManager struct {
	pool         *pgxpool.Pool
	connString   string
	maxRetries   int
	retryDelay   time.Duration
	lastConnTime time.Time
	mutex        sync.RWMutex
}

// NewDBManager creates a new database connection manager
func NewDBManager(connString string) *DBManager {
	return &DBManager{
		connString: connString,
		maxRetries: 5,
		retryDelay: 3 * time.Second,
	}
}

// Connect establishes a connection to the database with proper pool configuration
func (m *DBManager) Connect(ctx context.Context) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// Configure connection pool
	poolConfig, err := pgxpool.ParseConfig(m.connString)
	if err != nil {
		return fmt.Errorf("error parsing connection string: %w", err)
	}

	// Set optimal pool settings
	poolConfig.MaxConns = 10                                // Maximum number of connections
	poolConfig.MinConns = 2                                 // Minimum number of idle connections
	poolConfig.MaxConnLifetime = 1 * time.Hour              // Maximum connection lifetime
	poolConfig.MaxConnIdleTime = 30 * time.Minute           // Maximum idle time
	poolConfig.HealthCheckPeriod = 1 * time.Minute          // How often to check connection health
	poolConfig.ConnConfig.ConnectTimeout = 10 * time.Second // Connection timeout

	// Connect to database with retry logic
	var pool *pgxpool.Pool
	var lastErr error

	for i := 0; i < m.maxRetries; i++ {
		log.Printf("Attempting to connect to database (attempt %d/%d)...", i+1, m.maxRetries)

		pool, err = pgxpool.NewWithConfig(ctx, poolConfig)
		if err == nil {
			// Test connection
			if err = pool.Ping(ctx); err == nil {
				log.Println("Connected to database successfully")
				m.pool = pool
				m.lastConnTime = time.Now()
				return nil
			}
		}

		lastErr = err
		log.Printf("Failed to connect to database: %v. Retrying in %v...", err, m.retryDelay)
		time.Sleep(m.retryDelay)
	}

	return fmt.Errorf("failed to connect to database after %d attempts: %w", m.maxRetries, lastErr)
}

// GetPool returns the connection pool
func (m *DBManager) GetPool() *pgxpool.Pool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.pool
}

// EnsureConnected makes sure the database connection is alive
// Returns true if the connection was reestablished, false if it was already good
func (m *DBManager) EnsureConnected(ctx context.Context) (bool, error) {
	if m.pool == nil {
		if err := m.Connect(ctx); err != nil {
			return false, err
		}
		return true, nil
	}

	// Check if connection is still good
	if err := m.pool.Ping(ctx); err != nil {
		log.Println("Database connection lost, attempting to reconnect...")

		// Close the current pool
		if m.pool != nil {
			m.pool.Close()
			m.pool = nil
		}

		// Reconnect
		if err := m.Connect(ctx); err != nil {
			return false, err
		}
		return true, nil
	}

	return false, nil
}

// Close closes the connection pool
func (m *DBManager) Close() {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	if m.pool != nil {
		m.pool.Close()
		m.pool = nil
	}
}

// ExecuteWithRetry executes a database function with retry logic
func (m *DBManager) ExecuteWithRetry(ctx context.Context, operation string, fn func(ctx context.Context, pool *pgxpool.Pool) error) error {
	for i := 0; i < m.maxRetries; i++ {
		// Ensure we have a connection
		if _, err := m.EnsureConnected(ctx); err != nil {
			log.Printf("Failed to ensure database connection for %s: %v", operation, err)
			return err
		}

		// Get the pool
		pool := m.GetPool()
		if pool == nil {
			log.Printf("Pool is nil for %s operation", operation)
			time.Sleep(m.retryDelay)
			continue
		}

		// Execute the function
		err := fn(ctx, pool)
		if err == nil {
			return nil
		}

		log.Printf("Database operation %s failed: %v. Retry %d/%d",
			operation, err, i+1, m.maxRetries)

		// If we get a connection error, force reconnection on next attempt
		m.mutex.Lock()
		m.pool = nil
		m.mutex.Unlock()

		if i < m.maxRetries-1 {
			time.Sleep(m.retryDelay)
		}
	}

	return fmt.Errorf("operation %s failed after %d attempts", operation, m.maxRetries)
}
