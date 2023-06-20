package factory

import (
	"math/rand"

	"gitlab.spesolution.net/bni-merchant-management-system/go-sekeleton/internal/repository/mysql/entity"

	"github.com/google/uuid"
)

func StubbedWallets() []*entity.Wallet {
	return []*entity.Wallet{
		{
			UserName: uuid.NewString(),
			Balance:  rand.Int63(),
			ID:       rand.Int63(),
		},
	}
}

func StubbedWallet() *entity.Wallet {
	return StubbedWallets()[0]
}

// Example
// ID:             rand.Int63(),
// UserID:         rand.Int63(),
// Address:        uuid.NewString(),
// Latitude:       Float64Ptr(rand.Float64()),
// Longitude:      Float64Ptr(rand.Float64()),
// PostalCode:     rand.Int31(),
// DistrictID:     rand.Int31(),
// Source:         entity.SourceSmartseller,
// SourceID:       StrPtr(uuid.NewString()),
// Phone:          uuid.NewString(),
// RecipientName:  uuid.NewString(),
// AdditionalNote: StrPtr(uuid.NewString()),
// MainAddress:    false,
// CreatedAt:      time.Now(),
// UpdatedAt:      time.Now(),
