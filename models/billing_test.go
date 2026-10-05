package models

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBillingWallet_JSONSerialization(t *testing.T) {
	walletID := uuid.New()
	accountID := uuid.New()
	now := time.Now().Truncate(time.Second)

	tests := []struct {
		name   string
		wallet BillingWallet
	}{
		{
			name: "standard wallet",
			wallet: BillingWallet{
				ID:               walletID,
				BillingAccountID: accountID,
				Currency:         "NGN",
				BalanceMicros:    50000000,
				CreatedAt:        now,
				UpdatedAt:        now,
			},
		},
		{
			name: "zero balance wallet",
			wallet: BillingWallet{
				ID:               walletID,
				BillingAccountID: accountID,
				Currency:         "USD",
				BalanceMicros:    0,
				CreatedAt:        now,
				UpdatedAt:        now,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data, err := json.Marshal(tt.wallet)
			require.NoError(t, err)

			var deserialized BillingWallet
			err = json.Unmarshal(data, &deserialized)
			require.NoError(t, err)
			assert.Equal(t, tt.wallet.ID, deserialized.ID)
			assert.Equal(t, tt.wallet.BillingAccountID, deserialized.BillingAccountID)
			assert.Equal(t, tt.wallet.Currency, deserialized.Currency)
			assert.Equal(t, tt.wallet.BalanceMicros, deserialized.BalanceMicros)
		})
	}
}

func TestBillingWalletTransaction_JSONSerialization(t *testing.T) {
	refID := "inv-987"
	expires := time.Now().Add(30 * 24 * time.Hour).Truncate(time.Second)
	meta := json.RawMessage(`{"admin":"ops@devkitchen.com"}`)

	tx := BillingWalletTransaction{
		ID:                 uuid.New(),
		WalletID:           uuid.New(),
		BillingAccountID:   uuid.New(),
		Type:               WalletTxTypeGrant,
		AmountMicros:       25000000,
		BalanceAfterMicros: 25000000,
		Currency:           "NGN",
		ReferenceType:      WalletRefTypeAdminGrant,
		ReferenceID:        &refID,
		Description:        "Promotional starter grant",
		ExpiresAt:          &expires,
		Metadata:           meta,
		CreatedAt:          time.Now().Truncate(time.Second),
	}

	data, err := json.Marshal(tx)
	require.NoError(t, err)

	var deserialized BillingWalletTransaction
	err = json.Unmarshal(data, &deserialized)
	require.NoError(t, err)
	assert.Equal(t, tx.ID, deserialized.ID)
	assert.Equal(t, tx.Type, deserialized.Type)
	assert.Equal(t, tx.AmountMicros, deserialized.AmountMicros)
	assert.Equal(t, *tx.ReferenceID, *deserialized.ReferenceID)
}

func TestBillingWalletGrant_JSONSerialization(t *testing.T) {
	grantExp := time.Now().Add(60 * 24 * time.Hour).Truncate(time.Second)
	grant := BillingWalletGrant{
		ID:                    uuid.New(),
		WalletID:              uuid.New(),
		BillingAccountID:      uuid.New(),
		InitialAmountMicros:   50000000,
		RemainingAmountMicros: 20000000,
		Currency:              "NGN",
		ExpiresAt:             &grantExp,
		IsExpired:             false,
		CreatedAt:             time.Now().Truncate(time.Second),
		UpdatedAt:             time.Now().Truncate(time.Second),
	}

	data, err := json.Marshal(grant)
	require.NoError(t, err)

	var deserialized BillingWalletGrant
	err = json.Unmarshal(data, &deserialized)
	require.NoError(t, err)
	assert.Equal(t, grant.RemainingAmountMicros, deserialized.RemainingAmountMicros)
	assert.False(t, deserialized.IsExpired)
}

func TestAccountRateCardOverride_JSONSerialization(t *testing.T) {
	notes := "Enterprise negotiated vCPU rate"
	override := AccountRateCardOverride{
		ID:                       uuid.New(),
		BillingAccountID:         uuid.New(),
		SKU:                      "compute.vm.hour",
		CustomPricePerUnitMicros: 75000000,
		Notes:                    &notes,
		CreatedBy:                "enterprise-team",
		CreatedAt:                time.Now().Truncate(time.Second),
		UpdatedAt:                time.Now().Truncate(time.Second),
	}

	data, err := json.Marshal(override)
	require.NoError(t, err)

	var deserialized AccountRateCardOverride
	err = json.Unmarshal(data, &deserialized)
	require.NoError(t, err)
	assert.Equal(t, override.SKU, deserialized.SKU)
	assert.Equal(t, override.CustomPricePerUnitMicros, deserialized.CustomPricePerUnitMicros)
}

func TestTelemetryPayloads_JSONSerialization(t *testing.T) {
	vmPayload := VMLifecycleEventPayload{
		VMID:      uuid.New(),
		ProjectID: uuid.New(),
		VCPUs:     4,
		MemoryGB:  16,
		Timestamp: time.Now().Truncate(time.Second),
	}
	vmData, err := json.Marshal(vmPayload)
	require.NoError(t, err)
	var vmOut VMLifecycleEventPayload
	require.NoError(t, json.Unmarshal(vmData, &vmOut))
	assert.Equal(t, vmPayload.VCPUs, vmOut.VCPUs)

	dbPayload := DBLifecycleEventPayload{
		ClusterID:  uuid.New(),
		ProjectID:  uuid.New(),
		Engine:     "postgres",
		NodesCount: 3,
		Timestamp:  time.Now().Truncate(time.Second),
	}
	dbData, err := json.Marshal(dbPayload)
	require.NoError(t, err)
	var dbOut DBLifecycleEventPayload
	require.NoError(t, json.Unmarshal(dbData, &dbOut))
	assert.Equal(t, dbPayload.NodesCount, dbOut.NodesCount)

	s3Payload := S3UsageEventPayload{
		BucketID:     uuid.New(),
		ProjectID:    uuid.New(),
		BucketName:   "dk-123-assets",
		SizeBytes:    10737418240, // 10 GB
		ObjectsCount: 5000,
		Timestamp:    time.Now().Truncate(time.Second),
	}
	s3Data, err := json.Marshal(s3Payload)
	require.NoError(t, err)
	var s3Out S3UsageEventPayload
	require.NoError(t, json.Unmarshal(s3Data, &s3Out))
	assert.Equal(t, s3Payload.SizeBytes, s3Out.SizeBytes)
}

func TestInvoice_CreditFields_JSONSerialization(t *testing.T) {
	inv := Invoice{
		ID:                   uuid.New(),
		InvoiceNumber:        "INV-2026-001",
		BillingAccountID:     uuid.New(),
		Status:               InvoiceStatusPaid,
		SubtotalMicros:       100000000,
		TaxMicros:            7500000,
		TotalMicros:          107500000,
		CreditsAppliedMicros: 107500000,
		AmountDueMicros:      0,
		BaseCurrency:         "NGN",
		TargetCurrency:       "NGN",
		ExchangeRate:         1.0,
		ConvertedTotalMicros: 107500000,
		PeriodStart:          time.Now().Add(-30 * 24 * time.Hour),
		PeriodEnd:            time.Now(),
		DueDate:              time.Now().Add(14 * 24 * time.Hour),
		CreatedAt:            time.Now(),
	}

	data, err := json.Marshal(inv)
	require.NoError(t, err)

	var deserialized Invoice
	err = json.Unmarshal(data, &deserialized)
	require.NoError(t, err)
	assert.Equal(t, int64(107500000), deserialized.CreditsAppliedMicros)
	assert.Equal(t, int64(0), deserialized.AmountDueMicros)
}
