package factory

import (
	"math/rand"
	"time"

	"github.com/rahmatrdn/go-skeleton/internal/repository/mysql/entity"

	"github.com/google/uuid"
)

func StubbedTodoLists() []*entity.TodoList {
	return []*entity.TodoList{
		{
			Title:       uuid.NewString(),
			Description: uuid.NewString(),
			DoingAt:     time.Now(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			UserID:      rand.Int63(),
			ID:          rand.Int63(),
		},
	}
}

func StubbedTodoList() *entity.TodoList {
	return StubbedTodoLists()[0]
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
