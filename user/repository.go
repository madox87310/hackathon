package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Save(ctx context.Context, user *Model) error {
	if result := r.db.WithContext(ctx).Create(&user); result.Error != nil {
		return fmt.Errorf("r.db.Create: %w", result.Error)
	}
	return nil
}

func (r *Repository) Update(ctx context.Context, user *Model) error {
	if result := r.db.WithContext(ctx).Save(&user); result.Error != nil {
		return fmt.Errorf("r.db.WithContext(ctx).Save: %w", result.Error)
	}
	return nil
}

func (r *Repository) FindByID(ctx context.Context, id uuid.UUID) (*Model, error) {
	var user Model
	if result := r.db.WithContext(ctx).First(&user, id); result.Error != nil {
		return nil, fmt.Errorf("r.db.First: %w", result.Error)
	}
	return &user, nil
}

func (r *Repository) FindByPhoneNumber(ctx context.Context, phoneNumber string) (*Model, error) {
	var user Model
	result := r.db.WithContext(ctx).Where("phone_number = ?", phoneNumber).First(&user)
	if result.Error != nil {
		return nil, fmt.Errorf("r.db.WithContext(ctx).Where: %w", result.Error)
	}
	return &user, nil
}
