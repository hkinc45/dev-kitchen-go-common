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
	PricingPlanID                 *uuid.UUID      `json:"pricing_plan_id,omitempty" db:"pricing_plan_id"`
	PlanEnrolledAt                *time.Time      `json:"plan_enrolled_at,omitempty" db:"plan_enrolled_at"`
	PlanExpiresAt                 *time.Time      `json:"plan_expires_at,omitempty" db:"plan_expires_at"`
	FreeTierEligibility           string          `json:"free_tier_eligibility,omitempty" db:"free_tier_eligibility"`
	SubscriptionCancelAtPeriodEnd bool            `json:"subscription_cancel_at_period_end" db:"subscription_cancel_at_period_end"`
	SubscriptionCanceledAt        *time.Time      `json:"subscription_canceled_at,omitempty" db:"subscription_canceled_at"`
	SubscriptionStatus            string          `json:"subscription_status,omitempty" db:"subscription_status"`
	BillingCycle                  string          `json:"billing_cycle,omitempty" db:"billing_cycle"`
	BillingCycleAnchorDay         int             `json:"billing_cycle_anchor_day,omitempty" db:"billing_cycle_anchor_day"`
	PendingBillingCycle           *string         `json:"pending_billing_cycle,omitempty" db:"pending_billing_cycle"`
	PendingBillingCycleAnchorDay  *int            `json:"pending_billing_cycle_anchor_day,omitempty" db:"pending_billing_cycle_anchor_day"`
	PendingCycleEffectiveAt      *time.Time      `json:"pending_cycle_effective_at,omitempty" db:"pending_cycle_effective_at"`
	AutoApplyCredits              bool            `json:"auto_apply_credits" db:"auto_apply_credits"`
	PricingMode                   string          `json:"pricing_mode,omitempty" db:"pricing_mode"`
	IsActive                      bool            `json:"is_active" db:"is_active"`
	HasPaymentMethods             bool            `json:"has_payment_methods" db:"-"`
	CreatedAt                     time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt                     time.Time       `json:"updated_at" db:"updated_at"`
}

// Standard Billing Account and Project Billing Roles
const (
	BillingRoleOwner          = "owner"
	BillingRoleBillingManager = "billing-manager"
	BillingRoleBillingMember  = "billing-member"
	BillingRoleBillingViewer  = "billing-viewer"
)

// MemberProjectSource represents a linked project that grants billing access to a member.
type MemberProjectSource struct {
	ProjectID   uuid.UUID `json:"project_id"`
	ProjectName string    `json:"project_name,omitempty"`
	Role        string    `json:"role,omitempty"`
}

// BillingAccountMember maps user permissions within a specific billing account.
type BillingAccountMember struct {
	ID               uuid.UUID             `json:"id" db:"id"`
	BillingAccountID uuid.UUID             `json:"billing_account_id" db:"billing_account_id"`
	UserID           string                `json:"user_id" db:"user_id"`
	Role             string                `json:"role" db:"role"` // 'owner', 'billing-manager', 'billing-member', 'billing-viewer'
	Username         string                `json:"username,omitempty" db:"-"`
	FirstName        string                `json:"first_name,omitempty" db:"-"`
	LastName         string                `json:"last_name,omitempty" db:"-"`
	Email            string                `json:"email,omitempty" db:"-"`
	AvatarURL        string                `json:"avatar_url,omitempty" db:"-"`
	SourceProjects   []MemberProjectSource `json:"source_projects,omitempty" db:"-"`
	CreatedAt        time.Time             `json:"created_at" db:"created_at"`
}

// ProjectBillingBinding links a project to exactly one billing account (1:N topology).
type ProjectBillingBinding struct {
	ProjectID        uuid.UUID `json:"project_id" db:"project_id"`
	BillingAccountID uuid.UUID `json:"billing_account_id" db:"billing_account_id"`
	BoundAt          time.Time `json:"bound_at" db:"bound_at"`
	BoundByUserID    string    `json:"bound_by_user_id" db:"bound_by_user_id"`
}

// PlanEntitlement represents a feature, resource quota, SLA, or capability granted by a pricing plan.
type PlanEntitlement struct {
	Key   string      `json:"key"`
	Label string      `json:"label"`
	Value interface{} `json:"value"` // can be number, boolean, or string (e.g. "99.9%")
	Unit  string      `json:"unit,omitempty"`
	Type  string      `json:"type"` // "quota", "boolean", "badge", "limit"
	Icon  string      `json:"icon,omitempty"`
}

// ProductFamily groups plans into upgrade/downgrade hierarchies within a product domain.
type ProductFamily struct {
	ID            uuid.UUID  `json:"id" db:"id"`
	ProductTypeID *uuid.UUID `json:"product_type_id,omitempty" db:"product_type_id"`
	Code          string     `json:"code" db:"code"`
	Name          string     `json:"name" db:"name"`
	Description   *string    `json:"description,omitempty" db:"description"`
	IsStandalone  bool       `json:"is_standalone" db:"is_standalone"`
	IsActive      bool       `json:"is_active" db:"is_active"`
	CreatedAt     time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at" db:"updated_at"`
}

// ProductType defines high-level billing catalog domains and their default billing modes.
type ProductType struct {
	ID                    uuid.UUID       `json:"id" db:"id"`
	Code                  string          `json:"code" db:"code"`
	Name                  string          `json:"name,omitempty" db:"name"`
	DisplayName           string          `json:"display_name" db:"display_name"`
	Description           *string         `json:"description,omitempty" db:"description"`
	BillingMode           string          `json:"billing_mode" db:"billing_mode"` // 'metered', 'subscription', 'hybrid'
	IsSubscriptionEnabled bool            `json:"is_subscription_enabled" db:"is_subscription_enabled"`
	SKUPrefix             string          `json:"sku_prefix" db:"sku_prefix"`
	Icon                  string          `json:"icon" db:"icon"`
	Metadata              json.RawMessage `json:"metadata,omitempty" db:"metadata"`
	IsActive              bool            `json:"is_active" db:"is_active"`
	CreatedAt             time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time       `json:"updated_at" db:"updated_at"`
}

// PricingPlan defines subscription tiers and included resource allowances.
type PricingPlan struct {
	ID                      uuid.UUID         `json:"id" db:"id"`
	ProductTypeID           *uuid.UUID        `json:"product_type_id,omitempty" db:"product_type_id"`
	ProductFamilyID         *uuid.UUID        `json:"product_family_id,omitempty" db:"product_family_id"`
	TierRank                int               `json:"tier_rank" db:"tier_rank"`
	Name                    string            `json:"name" db:"name"`
	Slug                    string            `json:"slug" db:"slug"`
	Description             *string           `json:"description,omitempty" db:"description"`
	MonthlyFeeMicros        int64             `json:"monthly_fee_micros" db:"monthly_fee_micros"`
	IncludedVCPUHours       int               `json:"included_vcpu_hours" db:"included_vcpu_hours"`
	IncludedRAMGBHours      int               `json:"included_ram_gb_hours" db:"included_ram_gb_hours"`
	IncludedStorageGBMonths int               `json:"included_storage_gb_months" db:"included_storage_gb_months"`
	IncludedEgressGB        int               `json:"included_egress_gb" db:"included_egress_gb"`
	TrialDays               int               `json:"trial_days" db:"trial_days"` // 0 = unlimited / no trial expiration
	Category                string            `json:"category" db:"category"`     // 'compute', 'support', 'monitoring', 'security_scanning', 'disaster_recovery', 'all_in_one'
	PlanType                string            `json:"plan_type" db:"plan_type"`   // 'recurring', 'one_time'
	BillingInterval         string            `json:"billing_interval" db:"billing_interval"` // 'month', 'year', 'quarter', 'one_time'
	BillingIntervalCount    int               `json:"billing_interval_count" db:"billing_interval_count"`
	Features                []string          `json:"features" db:"features"`
	Entitlements            []PlanEntitlement `json:"entitlements" db:"entitlements"`
	IsDefault               bool              `json:"is_default" db:"is_default"`
	IsActive                bool              `json:"is_active" db:"is_active"`
	CreatedAt               time.Time         `json:"created_at" db:"created_at"`
}

// PaymentGatewayConfig defines the configuration and operational status of a payment gateway.
type PaymentGatewayConfig struct {
	ID                     string          `json:"id" db:"id"`
	Name                   string          `json:"name" db:"name"`
	Description            string          `json:"description" db:"description"`
	IsActive               bool            `json:"is_active" db:"is_active"`
	IsDefault              bool            `json:"is_default" db:"is_default"`
	SupportedCurrencies    []string        `json:"supported_currencies" db:"supported_currencies"`
	PublicKey              string          `json:"public_key" db:"public_key"`
	SecretKeyVaultPath     string          `json:"secret_key_vault_path,omitempty" db:"secret_key_vault_path"`
	WebhookSecretVaultPath string          `json:"webhook_secret_vault_path,omitempty" db:"webhook_secret_vault_path"`
	HasSecretKey           bool            `json:"has_secret_key" db:"has_secret_key"`
	HasWebhookSecret       bool            `json:"has_webhook_secret" db:"has_webhook_secret"`
	Metadata               json.RawMessage `json:"metadata,omitempty" db:"metadata"`
	CreatedAt              time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt              time.Time       `json:"updated_at" db:"updated_at"`
}

// PlanAccountEnrollment represents a billing account's active or past enrollment in a pricing plan.
type PlanAccountEnrollment struct {
	AccountID            uuid.UUID  `json:"account_id"`
	AccountName          string     `json:"account_name"`
	OwnerUserID          string     `json:"owner_user_id"`
	BillingEmail         string     `json:"billing_email"`
	Currency             string     `json:"currency"`
	PricingPlanID        *uuid.UUID `json:"pricing_plan_id,omitempty"`
	PlanName             string     `json:"plan_name"`
	IsDefault            bool       `json:"is_default"`
	PlanEnrolledAt       *time.Time `json:"plan_enrolled_at,omitempty"`
	PlanExpiresAt        *time.Time `json:"plan_expires_at,omitempty"`
	FreeTierEligibility string     `json:"free_tier_eligibility"`
	Status               string     `json:"status"` // 'active', 'expiring_soon', 'expired', 'unlimited'
	DaysRemaining        *int       `json:"days_remaining,omitempty"`
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
	InvoiceStatusRefunded      = "refunded"
)

// Invoice types.
const (
	InvoiceTypeCycle        = "cycle"
	InvoiceTypeSubscription = "subscription"
	InvoiceTypeManual       = "manual"
	InvoiceTypeUpfront      = "subscription" // Backward compatibility alias
)

// Invoice represents a finalized, auditable billing statement.
type Invoice struct {
	ID                   uuid.UUID     `json:"id" db:"id"`
	InvoiceNumber        string        `json:"invoice_number" db:"invoice_number"`
	BillingAccountID     uuid.UUID     `json:"billing_account_id" db:"billing_account_id"`
	InvoiceType          string        `json:"invoice_type" db:"invoice_type"`
	Status               string        `json:"status" db:"status"`
	SubtotalMicros       int64         `json:"subtotal_micros" db:"subtotal_micros"`
	TaxMicros            int64         `json:"tax_micros" db:"tax_micros"`
	TotalMicros          int64         `json:"total_micros" db:"total_micros"`
	CreditsAppliedMicros int64         `json:"credits_applied_micros" db:"credits_applied_micros"`
	AmountDueMicros      int64         `json:"amount_due_micros" db:"amount_due_micros"`
	BaseCurrency         string        `json:"base_currency" db:"base_currency"`
	TargetCurrency       string        `json:"target_currency" db:"target_currency"`
	ExchangeRate         float64       `json:"exchange_rate" db:"exchange_rate"`
	ConvertedTotalMicros int64         `json:"converted_total_micros" db:"converted_total_micros"`
	PeriodStart          time.Time     `json:"period_start" db:"period_start"`
	PeriodEnd            time.Time     `json:"period_end" db:"period_end"`
	DueDate              time.Time     `json:"due_date" db:"due_date"`
	PaidAt               *time.Time    `json:"paid_at,omitempty" db:"paid_at"`
	RefundedAt           *time.Time    `json:"refunded_at,omitempty" db:"refunded_at"`
	RefundReason         *string       `json:"refund_reason,omitempty" db:"refund_reason"`
	RefundAmountMicros   int64         `json:"refund_amount_micros,omitempty" db:"refund_amount_micros"`
	DeletedAt            *time.Time    `json:"deleted_at,omitempty" db:"deleted_at"`
	SubscriptionPlanID   *uuid.UUID    `json:"subscription_plan_id,omitempty" db:"subscription_plan_id"`
	PDFStoragePath       *string       `json:"pdf_storage_path,omitempty" db:"pdf_storage_path"`
	CreatedAt            time.Time     `json:"created_at" db:"created_at"`
	Items                []InvoiceItem `json:"items,omitempty" db:"-"`
}

// BillingPaymentTransaction records individual incoming payments, charges, and gateway transactions.
type BillingPaymentTransaction struct {
	ID                uuid.UUID       `json:"id" db:"id"`
	BillingAccountID  uuid.UUID       `json:"billing_account_id" db:"billing_account_id"`
	InvoiceID         *uuid.UUID      `json:"invoice_id,omitempty" db:"invoice_id"`
	InvoiceNumber     string          `json:"invoice_number,omitempty" db:"-"`
	AmountMicros      int64           `json:"amount_micros" db:"amount_micros"`
	Currency          string          `json:"currency" db:"currency"`
	PaymentMethodType string          `json:"payment_method_type" db:"payment_method_type"` // 'saved_card', 'online_gateway', 'wallet_credits', 'manual_admin', 'wallet_topup'
	Gateway           string          `json:"gateway" db:"gateway"`                         // 'paystack', 'stripe', 'wallet', 'admin'
	GatewayReference  string          `json:"gateway_reference" db:"gateway_reference"`
	Status            string          `json:"status" db:"status"`                           // 'succeeded', 'failed', 'pending'
	Metadata          json.RawMessage `json:"metadata" db:"metadata"`
	CreatedAt         time.Time       `json:"created_at" db:"created_at"`
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
	ProductType   string    `json:"product_type,omitempty"`
	Region        string    `json:"region,omitempty"`
	UserID        string    `json:"user_id,omitempty"`
	Timestamp     time.Time `json:"timestamp"`
}

// StorageUsageEventPayload is published on storage.allocated and storage.deleted.
type StorageUsageEventPayload struct {
	StorageID   string    `json:"storage_id"`
	ProjectID   uuid.UUID `json:"project_id"`
	SizeGB      int       `json:"size_gb"`
	ProductType string    `json:"product_type,omitempty"`
	Region      string    `json:"region,omitempty"`
	UserID      string    `json:"user_id,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
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

// --- Virtual Wallets & Cloud Credits ---

// BillingWallet represents a virtual prepaid wallet belonging to a billing account.
type BillingWallet struct {
	ID               uuid.UUID `json:"id" db:"id"`
	BillingAccountID uuid.UUID `json:"billing_account_id" db:"billing_account_id"`
	Currency         string    `json:"currency" db:"currency"`
	BalanceMicros    int64     `json:"balance_micros" db:"balance_micros"`
	CreatedAt        time.Time `json:"created_at" db:"created_at"`
	UpdatedAt        time.Time `json:"updated_at" db:"updated_at"`
}

// Wallet transaction types.
const (
	WalletTxTypeGrant      = "grant"
	WalletTxTypePurchase   = "purchase"
	WalletTxTypeDeduction  = "deduction"
	WalletTxTypeRefund     = "refund"
	WalletTxTypeExpiration = "expiration"
)

// Wallet reference types.
const (
	WalletRefTypeInvoice          = "invoice"
	WalletRefTypeAdminGrant       = "admin_grant"
	WalletRefTypeStripeCheckout   = "stripe_checkout"
	WalletRefTypePaystackCheckout = "paystack_checkout"
	WalletRefTypeExpirationJob    = "expiration_job"
)

// BillingWalletTransaction provides an immutable double-entry audit record for any balance change.
type BillingWalletTransaction struct {
	ID                 uuid.UUID       `json:"id" db:"id"`
	WalletID           uuid.UUID       `json:"wallet_id" db:"wallet_id"`
	BillingAccountID   uuid.UUID       `json:"billing_account_id" db:"billing_account_id"`
	Type               string          `json:"type" db:"type"` // 'grant', 'purchase', 'deduction', 'refund', 'expiration'
	AmountMicros       int64           `json:"amount_micros" db:"amount_micros"`
	BalanceAfterMicros int64           `json:"balance_after_micros" db:"balance_after_micros"`
	Currency           string          `json:"currency" db:"currency"`
	ReferenceType      string          `json:"reference_type" db:"reference_type"` // 'invoice', 'admin_grant', 'stripe_checkout', 'paystack_checkout', 'expiration_job'
	ReferenceID        *string         `json:"reference_id,omitempty" db:"reference_id"`
	Description        string          `json:"description" db:"description"`
	ExpiresAt          *time.Time      `json:"expires_at,omitempty" db:"expires_at"`
	Metadata           json.RawMessage `json:"metadata,omitempty" db:"metadata"`
	CreatedAt          time.Time       `json:"created_at" db:"created_at"`
}

// BillingWalletGrant tracks promotional grants with expiration dates for FIFO consumption.
type BillingWalletGrant struct {
	ID                    uuid.UUID `json:"id" db:"id"`
	WalletID              uuid.UUID `json:"wallet_id" db:"wallet_id"`
	BillingAccountID      uuid.UUID `json:"billing_account_id" db:"billing_account_id"`
	InitialAmountMicros   int64     `json:"initial_amount_micros" db:"initial_amount_micros"`
	RemainingAmountMicros int64     `json:"remaining_amount_micros" db:"remaining_amount_micros"`
	Currency              string    `json:"currency" db:"currency"`
	ExpiresAt             *time.Time `json:"expires_at" db:"expires_at"`
	IsExpired             bool      `json:"is_expired" db:"is_expired"`
	ApplicableProductTypes string   `json:"applicable_product_types" db:"applicable_product_types"` // 'all' or comma-separated e.g. 'storage_s3,compute_vm'
	CreatedAt             time.Time `json:"created_at" db:"created_at"`
	UpdatedAt             time.Time `json:"updated_at" db:"updated_at"`
}

// AccountRateCardOverride defines negotiated per-account unit prices for specific SKUs.
type AccountRateCardOverride struct {
	ID                       uuid.UUID `json:"id" db:"id"`
	BillingAccountID         uuid.UUID `json:"billing_account_id" db:"billing_account_id"`
	SKU                      string    `json:"sku" db:"sku"`
	CustomPricePerUnitMicros int64     `json:"custom_price_per_unit_micros" db:"custom_price_per_unit_micros"`
	Notes                    *string   `json:"notes,omitempty" db:"notes"`
	CreatedBy                string    `json:"created_by" db:"created_by"`
	CreatedAt                time.Time `json:"created_at" db:"created_at"`
	UpdatedAt                time.Time `json:"updated_at" db:"updated_at"`
}

// --- Infrastructure Primitives Telemetry Payloads ---

// VMLifecycleEventPayload is published on vm.started and vm.stopped.
type VMLifecycleEventPayload struct {
	VMID        uuid.UUID `json:"vm_id"`
	ProjectID   uuid.UUID `json:"project_id"`
	VCPUs       int       `json:"vcpus"`
	MemoryGB    int       `json:"memory_gb"`
	ProductType string    `json:"product_type,omitempty"`
	Region      string    `json:"region,omitempty"`
	UserID      string    `json:"user_id,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
}

// DBLifecycleEventPayload is published on db.started and db.stopped.
type DBLifecycleEventPayload struct {
	ClusterID   uuid.UUID `json:"cluster_id"`
	ProjectID   uuid.UUID `json:"project_id"`
	Engine      string    `json:"engine"`
	NodesCount  int       `json:"nodes_count"`
	ProductType string    `json:"product_type,omitempty"`
	Region      string    `json:"region,omitempty"`
	UserID      string    `json:"user_id,omitempty"`
	Timestamp   time.Time `json:"timestamp"`
}

// S3UsageEventPayload is published on storage.s3.usage_polled.
type S3UsageEventPayload struct {
	BucketID     uuid.UUID `json:"bucket_id"`
	ProjectID    uuid.UUID `json:"project_id"`
	BucketName   string    `json:"bucket_name"`
	SizeBytes    int64     `json:"size_bytes"`
	ObjectsCount int64     `json:"objects_count"`
	ProductType  string    `json:"product_type,omitempty"`
	Region       string    `json:"region,omitempty"`
	UserID       string    `json:"user_id,omitempty"`
	Timestamp    time.Time `json:"timestamp"`
}

// PlatformControls holds global operational billing parameters and debt ceilings.
type PlatformControls struct {
	ID                     string    `json:"id" db:"id"`
	GracePeriodDays        int       `json:"grace_period_days" db:"grace_period_days"`
	MaxOverdueDebtMicros   int64     `json:"max_overdue_debt_micros" db:"max_overdue_debt_micros"`
	AutoRetryIntervalHours int       `json:"auto_retry_interval_hours" db:"auto_retry_interval_hours"`
	AutoChargeEnabled      bool      `json:"auto_charge_enabled" db:"auto_charge_enabled"`
	DefaultTaxPercent      float64   `json:"default_tax_percent" db:"default_tax_percent"`
	PlatformBaseCurrency   string    `json:"platform_base_currency" db:"platform_base_currency"`
	DefaultPaymentGateway  string    `json:"default_payment_gateway" db:"default_payment_gateway"`
	DefaultPricingMode     string    `json:"default_pricing_mode" db:"default_pricing_mode"`
	UpdatedAt              time.Time `json:"updated_at" db:"updated_at"`
}

// PricingMode represents a platform-wide or product-level billing strategy (metered, subscription, hybrid).
type PricingMode struct {
	ID             string                 `json:"id" db:"id"`
	Name           string                 `json:"name" db:"name"`
	DisplayName    string                 `json:"display_name" db:"display_name"`
	Description    string                 `json:"description,omitempty" db:"description"`
	IsEnabled      bool                   `json:"is_enabled" db:"is_enabled"`
	IsDefault      bool                   `json:"is_default" db:"is_default"`
	IsFoundational bool                   `json:"is_foundational" db:"is_foundational"`
	SortOrder      int                    `json:"sort_order" db:"sort_order"`
	Metadata       map[string]interface{} `json:"metadata,omitempty" db:"metadata"`
	CreatedAt      time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time              `json:"updated_at" db:"updated_at"`
}

// MeterDefinition represents a plug-and-play schema for the usage aggregation and rating engine.
type MeterDefinition struct {
	ID                   uuid.UUID              `json:"id" db:"id"`
	ProductTypeID        uuid.UUID              `json:"product_type_id" db:"product_type_id"`
	Code                 string                 `json:"code" db:"code"`
	Name                 string                 `json:"name" db:"name"`
	Description          string                 `json:"description,omitempty" db:"description"`
	EventType            string                 `json:"event_type" db:"event_type"`
	ResourceType         string                 `json:"resource_type" db:"resource_type"`
	MetricField          string                 `json:"metric_field" db:"metric_field"`
	FilterCriteria       map[string]interface{} `json:"filter_criteria,omitempty" db:"filter_criteria"`
	AggregationType      string                 `json:"aggregation_type" db:"aggregation_type"` // 'gauge', 'avg', 'sum', 'duration_seconds'
	SourceUnit           string                 `json:"source_unit" db:"source_unit"`
	TargetUnit           string                 `json:"target_unit" db:"target_unit"`
	UnitConversionFactor float64                `json:"unit_conversion_factor" db:"unit_conversion_factor"`
	TargetSKU            string                 `json:"target_sku" db:"target_sku"`
	PlanQuotaKey         *string                `json:"plan_quota_key,omitempty" db:"plan_quota_key"`
	IsActive             bool                   `json:"is_active" db:"is_active"`
	CreatedAt            time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time              `json:"updated_at" db:"updated_at"`
}

