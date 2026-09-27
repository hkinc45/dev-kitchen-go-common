package gateways

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

// CustomerInput contains parameters to register a customer on the payment gateway.
type CustomerInput struct {
	Email    string            `json:"email"`
	Name     string            `json:"name"`
	Metadata map[string]string `json:"metadata,omitempty"`
}

// CustomerOutput represents the created customer reference on the gateway.
type CustomerOutput struct {
	CustomerID string          `json:"customer_id"`
	Raw        json.RawMessage `json:"raw,omitempty"`
}

// PaymentIntentInput parameters for initializing a payment or checkout session.
type PaymentIntentInput struct {
	AmountMicros int64             `json:"amount_micros"`
	Currency     string            `json:"currency"`
	CustomerID   string            `json:"customer_id,omitempty"`
	Description  string            `json:"description,omitempty"`
	Metadata     map[string]string `json:"metadata,omitempty"`
}

// PaymentIntentOutput returns intent IDs and client secrets for checkout completion.
type PaymentIntentOutput struct {
	PaymentIntentID string          `json:"payment_intent_id"`
	ClientSecret    string          `json:"client_secret,omitempty"`
	Status          string          `json:"status"`
	CheckoutURL     string          `json:"checkout_url,omitempty"` // For gateways redirecting to hosted checkout (e.g. Paystack)
	Raw             json.RawMessage `json:"raw,omitempty"`
}

// ChargeInput parameters for direct charging off-session payment instruments.
type ChargeInput struct {
	AmountMicros    int64             `json:"amount_micros"`
	Currency        string            `json:"currency"`
	CustomerID      string            `json:"customer_id"`
	PaymentMethodID string            `json:"payment_method_id,omitempty"`
	Description     string            `json:"description,omitempty"`
	Metadata        map[string]string `json:"metadata,omitempty"`
}

// ChargeOutput details of an executed charge.
type ChargeOutput struct {
	ChargeID   string          `json:"charge_id"`
	Status     string          `json:"status"`
	ReceiptURL string          `json:"receipt_url,omitempty"`
	Raw        json.RawMessage `json:"raw,omitempty"`
}

// WebhookEvent represents a normalized payment event from incoming provider webhooks.
type WebhookEvent struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	Gateway   string          `json:"gateway"`
	Data      json.RawMessage `json:"data"`
	Timestamp time.Time       `json:"timestamp"`
	Raw       []byte          `json:"-"`
}

// PaymentGateway defines the unified provider-agnostic payment gateway contract.
type PaymentGateway interface {
	Name() string
	CreateCustomer(ctx context.Context, input CustomerInput) (*CustomerOutput, error)
	CreatePaymentIntent(ctx context.Context, input PaymentIntentInput) (*PaymentIntentOutput, error)
	Charge(ctx context.Context, input ChargeInput) (*ChargeOutput, error)
	VerifyWebhookSignature(req *http.Request, secret string) ([]byte, error)
	ParseWebhookEvent(payload []byte) (*WebhookEvent, error)
}
