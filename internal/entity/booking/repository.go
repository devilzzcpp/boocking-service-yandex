package booking

import (
	"context"
	"time"

	"booking_service/internal/models"

	"gorm.io/gorm"
)

type Repository interface {
	HasConflict(ctx context.Context, placeID int, from, to time.Time) (bool, error)
	Create(ctx context.Context, userID, placeID int, from, to time.Time) error
	ListByUser(ctx context.Context, userID int) ([]models.Booking, error)
	ListByPlace(ctx context.Context, placeID int) ([]models.Booking, error)
}

type repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &repository{db: db}
}

func (r *repository) HasConflict(ctx context.Context, placeID int, from, to time.Time) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&models.Booking{}).
		Where("place_id = ? AND time_from < ? AND ? < time_to", placeID, to.UTC(), from.UTC()).
		Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *repository) Create(ctx context.Context, userID, placeID int, from, to time.Time) error {
	item := models.Booking{
		UserID:  userID,
		PlaceID: placeID,
		From:    from.UTC(),
		To:      to.UTC(),
	}
	return r.db.WithContext(ctx).Create(&item).Error
}

func (r *repository) ListByUser(ctx context.Context, userID int) ([]models.Booking, error) {
	result := make([]models.Booking, 0)
	err := r.db.WithContext(ctx).
		Where("user_id = ?", userID).
		Order("time_from ASC, id ASC").
		Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

func (r *repository) ListByPlace(ctx context.Context, placeID int) ([]models.Booking, error) {
	result := make([]models.Booking, 0)
	err := r.db.WithContext(ctx).
		Where("place_id = ?", placeID).
		Order("time_from ASC, id ASC").
		Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}
