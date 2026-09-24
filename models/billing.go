package models

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// BillingAccount represents an independent legal, tax, and invoicing identity.
type BillingAccount struct {
	ID            uuid.UUID       `json:"id" db:"id"`
	Name          string          `json:"name" db:"name"`
	OwnerUserID   string          `json:"owner_user_id" db:"owner_user_id"`
	Currency      string          `json:"currency" db:"currency"`
	TaxID         *string         `json:"tax_id,omitempty" db:"tax_id"`
	BillingEmail  string          `json:"billing_email" db:"billing_email"`
	Address       json.RawMessage `json:"address,omitempty" db:"address"`
	PricingPlanID *uuid.UUID      `json:"pricing_plan_id,omitempty" db:"pricing_plan_id"`
	IsActive      bool            `json:"is_active" db:"is_active"`
	CreatedAt     time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time       `json:"updated_at" db:"updated_at"`
}

// BillingAccountMember maps user permissions within a specific billing account.
type BillingAccountMember struct {
	ID               uuid.UUID `json:"id" db:"id"`
	BillingAccountID uuid.UUID `json:"billing_account_id" db:"billing_account_id"`
	UserID           string    `json:"user_id" db:"user_id"`
	Role             string    `json:"role" db:"role"` // 'admin', 'viewer'
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

// ProjectBillingBinding links a project to exactly one billing account (1:N topology).
type ProjectBillingBinding struct {
	ProjectID        uuid.UUID `json:"project_id" db:"project_id"`
	BillingAccountID uuid.UUID `json:"billing_account_id" db:"billing_account_id"`
	BoundAt          time.Time `json:"bound_at" db:"bound_at"`
	BoundByUserID    string    `json:"bound_by_user_id" db:"bound_by_user_id"`
}

// PricingPlan defines subscription tiers and included resource allowances.
type PricingPlan struct {
	ID                      uuid.UUID `json:"id" db:"id"`
	Name                    string    `json:"name" db:"name"`
	Slug                    string    `json:"slug" db:"slug"`
	Description             *string   `json:"description,omitempty" db:"description"`
	MonthlyFeeMicros        int64     `json:"monthly_fee_micros" db:"monthly_fee_micros"`
	IncludedVCPUHours       int       `json:"included_vcpu_hours" db:"included_vcpu_hours"`
	IncludedRAMGBHours      int       `json:"included_ram_gb_hours" db:"included_ram_gb_hours"`
	IncludedStorageGBMonths int       `json:"included_storage_gb_months" db:"included_storage_gb_months"`
	IncludedEgressGB        int       `json:"included_egress_gb" db:"included_egress_gb"`
	IsDefault               bool      `json:"is_default" db:"is_default"`
	IsActive                bool      `json:"is_active" db:"is_active"`
	CreatedAt               time.Time `json:"created_at" db:"created_at"`
}

// RateCard defines unit prices for billable resource overages in USD micros.
type RateCard struct {
	SKU                string `json:"sku" db:"sku"`
	Name               string `json:"name" db:"name"`
	Unit               string `json:"unit" db:"unit"`
	PricePerUnitMicros int64  `json:"price_per_unit_micros" db:"price_per_unit_micros"`
	IsActive           bool   `json:"is_active" db:"is_active"`
}

// ExchangeRate holds base-to-target currency conversion rates with platform margins.
type ExchangeRate struct {
	BaseCurrency   string    `json:"base_currency" db:"base_currency"`
	TargetCurrency string    `json:"target_currency" db:"target_currency"`
	Rate           float64   `json:"rate" db:"rate"`
	MarginPercent  float64   `json:"margin_percent" db:"margin_percent"`
	Source         string    `json:"source" db:"source"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// ExchangeRateOverride contains administrator-set fixed rate overrides.
type ExchangeRateOverride struct {
	Currency      string    `json:"currency" db:"currency"`
	Rate          float64   `json:"rate" db:"rate"`
	MarginPercent float64   `json:"margin_percent" db:"margin_percent"`
	OverrideBy    string    `json:"override_by" db:"override_by"`
	Reason        *string   `json:"reason,omitempty" db:"reason"`
	UpdatedAt     time.Time `json:"updated_at" db:"updated_at"`
}

// UsageEvent stores raw, deduplicated resource lifecycle events.
type UsageEvent struct {
	ID               uuid.UUID       `json:"id" db:"id"`
	IdempotencyKey   string          `json:"idempotency_key" db:"idempotency_key"`
	BillingAccountID uuid.UUID       `json:"billing_account_id" db:"billing_account_id"`
	ProjectID        uuid.UUID       `json:"project_id" db:"project_id"`
	ResourceID       string          `json:"resource_id" db:"resource_id"`
	ResourceType     string          `json:"resource_type" db:"resource_type"`
	EventType        string          `json:"event_type" db:"event_type"`
	MetricValue      int64           `json:"metric_value" db:"metric_value"`
	Unit             string          `json:"unit" db:"unit"`
	Timestamp        time.Time       `json:"timestamp" db:"timestamp"`
	RawPayload       json.RawMessage `json:"raw_payload,omitempty" db:"raw_payload"`
}

// UsageRecord aggregates resource consumption over an hourly or daily interval.
type UsageRecord struct {
	ID               uuid.UUID `json:"id" db:"id"`
	BillingAccountID uuid.UUID `json:"billing_account_id" db:"billing_account_id"`
	ProjectID        uuid.UUID `json:"project_id" db:"project_id"`
	ResourceID       string    `json:"resource_id" db:"resource_id"`
	MetricType       string    `json:"metric_type" db:"metric_type"`
	QuantityMicros   int64     `json:"quantity_micros" db:"quantity_micros"`
	CostMicros       int64     `json:"cost_micros" db:"cost_micros"`
	Currency         string    `json:"currency" db:"currency"`
	PeriodStart      time.Time `json:"period_start" db:"period_start"`
	PeriodEnd        time.Time `json:"period_end" db:"period_end"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
}

// Invoice statuses.
const (
	InvoiceStatusDraft         = "draft"
	InvoiceStatusOpen          = "open"
	InvoiceStatusPaid          = "paid"
	InvoiceStatusVoid          = "void"
	InvoiceStatusUncollectible = "uncollectible"
)

// Invoice represents a finalized, auditable billing statement.
type Invoice struct {
	ID                   uuid.UUID     `json:"id" db:"id"`
	InvoiceNumber        string        `json:"invoice_number" db:"invoice_number"`
	BillingAccountID     uuid.UUID     `json:"billing_account_id" db:"billing_account_id"`
	Status               string        `json:"status" db:"status"`
	SubtotalMicros       int64         `json:"subtotal_micros" db:"subtotal_micros"`
	TaxMicros            int64         `json:"tax_micros" db:"tax_micros"`
	TotalMicros          int64         `json:"total_micros" db:"total_micros"`
	BaseCurrency         string        `json:"base_currency" db:"base_currency"`
	TargetCurrency       string        `json:"target_currency" db:"target_currency"`
	ExchangeRate         float64       `json:"exchange_rate" db:"exchange_rate"`
	ConvertedTotalMicros int64         `json:"converted_total_micros" db:"converted_total_micros"`
	PeriodStart          time.Time     `json:"period_start" db:"period_start"`
	PeriodEnd            time.Time     `json:"period_end" db:"period_end"`
	DueDate              time.Time     `json:"due_date" db:"due_date"`
	PaidAt               *time.Time    `json:"paid_at,omitempty" db:"paid_at"`
	PDFStoragePath       *string       `json:"pdf_storage_path,omitempty" db:"pdf_storage_path"`
	CreatedAt            time.Time     `json:"created_at" db:"created_at"`
	Items                []InvoiceItem `json:"items,omitempty" db:"-"`
}

// InvoiceItem represents an individual line item on an invoice.
type InvoiceItem struct {
	ID              uuid.UUID  `json:"id" db:"id"`
	InvoiceID       uuid.UUID  `json:"invoice_id" db:"invoice_id"`
	ProjectID       *uuid.UUID `json:"project_id,omitempty" db:"project_id"`
	SKU             string     `json:"sku" db:"sku"`
	Description     string     `json:"description" db:"description"`
	Quantity        int64      `json:"quantity" db:"quantity"`
	UnitPriceMicros int64      `json:"unit_price_micros" db:"unit_price_micros"`
	TotalMicros     int64      `json:"total_micros" db:"total_micros"`
}

// PaymentMethod stores tokenized customer payment instruments.
type PaymentMethod struct {
	ID                     uuid.UUID `json:"id" db:"id"`
	BillingAccountID       uuid.UUID `json:"billing_account_id" db:"billing_account_id"`
	Gateway                string    `json:"gateway" db:"gateway"` // 'stripe', 'paystack'
	GatewayCustomerID      string    `json:"gateway_customer_id" db:"gateway_customer_id"`
	GatewayPaymentMethodID string    `json:"gateway_payment_method_id" db:"gateway_payment_method_id"`
	Brand                  *string   `json:"brand,omitempty" db:"brand"`
	Last4                  *string   `json:"last4,omitempty" db:"last4"`
	ExpMonth               *int      `json:"exp_month,omitempty" db:"exp_month"`
	ExpYear                *int      `json:"exp_year,omitempty" db:"exp_year"`
	IsDefault              bool      `json:"is_default" db:"is_default"`
	CreatedAt              time.Time `json:"created_at" db:"created_at"`
}

// ProcessedWebhookEvent tracks processed webhooks for cryptographic replay prevention.
type ProcessedWebhookEvent struct {
	ID          uuid.UUID       `json:"id" db:"id"`
	Gateway     string          `json:"gateway" db:"gateway"`
	EventID     string          `json:"event_id" db:"event_id"`
	EventType   string          `json:"event_type" db:"event_type"`
	Payload     json.RawMessage `json:"payload" db:"payload"`
	Status      string          `json:"status" db:"status"`
	ProcessedAt time.Time       `json:"processed_at" db:"processed_at"`
}

// Budget represents monthly spending limits and proactive threshold alert configs.
type Budget struct {
	ID                 uuid.UUID  `json:"id" db:"id"`
	BillingAccountID   uuid.UUID  `json:"billing_account_id" db:"billing_account_id"`
	ProjectID          *uuid.UUID `json:"project_id,omitempty" db:"project_id"`
	MonthlyLimitMicros int64      `json:"monthly_limit_micros" db:"monthly_limit_micros"`
	Currency           string     `json:"currency" db:"currency"`
	AlertEmails        []string   `json:"alert_emails" db:"alert_emails"`
	NotifySlack        bool       `json:"notify_slack" db:"notify_slack"`
	WebhookURL         *string    `json:"webhook_url,omitempty" db:"webhook_url"`
	AutoScaleDown      bool       `json:"auto_scale_down" db:"auto_scale_down"`
	CreatedAt          time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at" db:"updated_at"`
}

// BudgetAlert records triggered threshold warnings.
type BudgetAlert struct {
	ID                   uuid.UUID `json:"id" db:"id"`
	BudgetID             uuid.UUID `json:"budget_id" db:"budget_id"`
	ThresholdPercent     int       `json:"threshold_percent" db:"threshold_percent"` // 50, 80, 100, 120
	SpendAtTriggerMicros int64     `json:"spend_at_trigger_micros" db:"spend_at_trigger_micros"`
	TriggeredAt          time.Time `json:"triggered_at" db:"triggered_at"`
	NotificationStatus   string    `json:"notification_status" db:"notification_status"`
}

// BillingAccountAuditLog tracks administrative account and binding mutations.
type BillingAccountAuditLog struct {
	ID               uuid.UUID       `json:"id" db:"id"`
	BillingAccountID uuid.UUID       `json:"billing_account_id" db:"billing_account_id"`
	ActorUserID      string          `json:"actor_user_id" db:"actor_user_id"`
	Action           string          `json:"action" db:"action"`
	Details          json.RawMessage `json:"details" db:"details"`
	CreatedAt        time.Time       `json:"created_at" db:"created_at"`
}

// --- NATS JetStream Event Payloads ---

// WorkloadUsageEventPayload is published on workload.started and workload.stopped.
type WorkloadUsageEventPayload struct {
	WorkloadID    uuid.UUID `json:"workload_id"`
	ProjectID     uuid.UUID `json:"project_id"`
	CPULimitM     int       `json:"cpu_limit_m"`
	MemoryLimitMi int       `json:"memory_limit_mi"`
	Timestamp     time.Time `json:"timestamp"`
}

// StorageUsageEventPayload is published on storage.allocated and storage.deleted.
type StorageUsageEventPayload struct {
	StorageID string    `json:"storage_id"`
	ProjectID uuid.UUID `json:"project_id"`
	SizeGB    int       `json:"size_gb"`
	Timestamp time.Time `json:"timestamp"`
}

// BudgetThresholdReachedPayload is published on billing.budget.threshold_reached.
type BudgetThresholdReachedPayload struct {
	BudgetID         uuid.UUID  `json:"budget_id"`
	BillingAccountID uuid.UUID  `json:"billing_account_id"`
	ProjectID        *uuid.UUID `json:"project_id,omitempty"`
	ThresholdPercent int        `json:"threshold_percent"`
	SpendMicros      int64      `json:"spend_micros"`
	LimitMicros      int64      `json:"limit_micros"`
	Timestamp        time.Time  `json:"timestamp"`
}

// BudgetScaleDownPayload is published on billing.budget.action.scale_down.
type BudgetScaleDownPayload struct {
	BillingAccountID uuid.UUID `json:"billing_account_id"`
	ProjectID        uuid.UUID `json:"project_id"`
	Reason           string    `json:"reason"`
	Timestamp        time.Time `json:"timestamp"`
}

// PaymentSucceededPayload is published on billing.payment.succeeded.
type PaymentSucceededPayload struct {
	InvoiceID        uuid.UUID `json:"invoice_id"`
	BillingAccountID uuid.UUID `json:"billing_account_id"`
	AmountMicros     int64     `json:"amount_micros"`
	Currency         string    `json:"currency"`
	Gateway          string    `json:"gateway"`
	PaymentIntentID  string    `json:"payment_intent_id"`
	Timestamp        time.Time `json:"timestamp"`
}
