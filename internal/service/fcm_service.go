package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/barberku/backend-barber/internal/entity"
	"github.com/barberku/backend-barber/internal/repository"
)

type FCMService interface {
	RegisterToken(ctx context.Context, customerID, token, platform string) error
	SendNotification(ctx context.Context, customerID, title, body string) error
}

type FCMServiceImpl struct {
	fcmRepo    repository.FCMTokenRepository
	serverKey  string
	httpClient *http.Client
}

func NewFCMService(fcmRepo repository.FCMTokenRepository, serverKey string) *FCMServiceImpl {
	return &FCMServiceImpl{
		fcmRepo:    fcmRepo,
		serverKey:  serverKey,
		httpClient: &http.Client{Timeout: 10 * time.Second},
	}
}

func (s *FCMServiceImpl) RegisterToken(ctx context.Context, customerID, token, platform string) error {
	fcmToken := &entity.FCMToken{
		CustomerID: customerID,
		Token:      token,
		Platform:   platform,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	return s.fcmRepo.Save(ctx, fcmToken)
}

func (s *FCMServiceImpl) SendNotification(ctx context.Context, customerID, title, body string) error {
	tokens, err := s.fcmRepo.GetByCustomerID(ctx, customerID)
	if err != nil {
		return fmt.Errorf("failed to get FCM tokens: %w", err)
	}

	if len(tokens) == 0 {
		slog.Warn("no FCM tokens found for customer", "customer_id", customerID)
		return nil
	}

	for _, t := range tokens {
		if err := s.sendToFCM(ctx, t.Token, title, body); err != nil {
			slog.Error("failed to send FCM notification", "token_id", t.ID, "error", err)
		}
	}

	return nil
}

func (s *FCMServiceImpl) sendToFCM(ctx context.Context, token, title, body string) error {
	payload := map[string]interface{}{
		"to": token,
		"notification": map[string]string{
			"title": title,
			"body":  body,
		},
		"data": map[string]string{
			"type": "queue_called",
		},
	}

	bodyBytes, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://fcm.googleapis.com/fcm/send", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "key="+s.serverKey)

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("FCM returned status %d: %s", resp.StatusCode, string(respBody))
	}

	return nil
}
