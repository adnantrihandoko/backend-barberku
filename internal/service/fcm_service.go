package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"firebase.google.com/go/v4/messaging"

	"github.com/barberku/backend-barber/internal/entity"
	"github.com/barberku/backend-barber/internal/repository"
)

type FCMService interface {
	RegisterToken(ctx context.Context, customerID, token, platform string) error
	SendNotification(ctx context.Context, customerID, title, body string) error
}

type FCMServiceImpl struct {
	fcmRepo  repository.FCMTokenRepository
	client   *messaging.Client
	enabled  bool
}

func NewFCMService(fcmRepo repository.FCMTokenRepository, client *messaging.Client) *FCMServiceImpl {
	enabled := client != nil
	if !enabled {
		slog.Warn("FCM client is nil, push notifications disabled")
	}
	return &FCMServiceImpl{
		fcmRepo: fcmRepo,
		client:  client,
		enabled: enabled,
	}
}

func (s *FCMServiceImpl) RegisterToken(ctx context.Context, customerID, token, platform string) error {
	fcmToken := &entity.FCMToken{
		ID:         fmt.Sprintf("%s-%s", customerID, token),
		CustomerID: customerID,
		Token:      token,
		Platform:   platform,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	return s.fcmRepo.Save(ctx, fcmToken)
}

func (s *FCMServiceImpl) SendNotification(ctx context.Context, customerID, title, body string) error {
	if !s.enabled {
		slog.Warn("FCM not enabled, skipping notification", "customer_id", customerID)
		return nil
	}

	tokens, err := s.fcmRepo.GetByCustomerID(ctx, customerID)
	if err != nil {
		return fmt.Errorf("failed to get FCM tokens: %w", err)
	}

	if len(tokens) == 0 {
		slog.Warn("no FCM tokens found for customer", "customer_id", customerID)
		return nil
	}

	tokenStrings := make([]string, 0, len(tokens))
	for _, t := range tokens {
		tokenStrings = append(tokenStrings, t.Token)
	}

	message := &messaging.MulticastMessage{
		Tokens: tokenStrings,
		Notification: &messaging.Notification{
			Title: title,
			Body:  body,
		},
		Data: map[string]string{
			"type": "queue_called",
		},
	}

	response, err := s.client.SendEachForMulticast(ctx, message)
	if err != nil {
		return fmt.Errorf("failed to send FCM notification: %w", err)
	}

	slog.Info("FCM notification sent",
		"customer_id", customerID,
		"total", len(tokens),
		"success", response.SuccessCount,
		"failure", response.FailureCount,
	)

	if response.FailureCount > 0 {
		for i, result := range response.Responses {
			if result.Error != nil {
				slog.Error("FCM send failed for token",
					"index", i,
					"error", result.Error,
				)
			}
		}
	}

	return nil
}
