package tasks

import (
	"context"
	"log"
	"time"

	"mobilka/internal/service"
)

// SubscriptionChecker periodically checks for expired subscriptions
type SubscriptionChecker struct {
	paymentService *service.PaymentService
	interval       time.Duration
	stopChan       chan struct{}
}

// NewSubscriptionChecker creates a new subscription checker
func NewSubscriptionChecker(paymentService *service.PaymentService, interval time.Duration) *SubscriptionChecker {
	return &SubscriptionChecker{
		paymentService: paymentService,
		interval:       interval,
		stopChan:       make(chan struct{}),
	}
}

// Start starts the subscription checker
func (sc *SubscriptionChecker) Start() {
	go func() {
		ticker := time.NewTicker(sc.interval)
		defer ticker.Stop()

		// Run immediately on start
		sc.checkExpiredSubscriptions()

		for {
			select {
			case <-ticker.C:
				sc.checkExpiredSubscriptions()
			case <-sc.stopChan:
				log.Println("Subscription checker stopped")
				return
			}
		}
	}()

	log.Printf("Subscription checker started with interval: %s", sc.interval)
}

// Stop stops the subscription checker
func (sc *SubscriptionChecker) Stop() {
	close(sc.stopChan)
}

// checkExpiredSubscriptions checks for expired subscriptions and updates their status
func (sc *SubscriptionChecker) checkExpiredSubscriptions() {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	// Try to expire subscriptions, but handle database schema issues gracefully
	count, err := sc.paymentService.ExpireSubscriptions(ctx)
	if err != nil {
		// Log the error but don't panic - the service might be starting up
		// or migrations might not have completed yet
		log.Printf("Error checking expired subscriptions: %v", err)

		// If it's a database schema error, we'll try again later
		// No need to take additional action now
		return
	}

	if count > 0 {
		log.Printf("Expired %d subscriptions", count)
	}
}
