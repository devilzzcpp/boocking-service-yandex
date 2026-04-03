package booking

import (
	"context"
	"errors"
	"time"

	"booking_service/internal/models"
)

var ErrConflict = errors.New("booking conflict")

type Service interface {
	Create(ctx context.Context, userID, placeID int, from, to time.Time) error
	List(ctx context.Context, userID *int, placeID *int) ([]models.Booking, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, userID, placeID int, from, to time.Time) error {
	conflict, err := s.repo.HasConflict(ctx, placeID, from, to)
	if err != nil {
		return err
	}
	if conflict {
		return ErrConflict
	}
	return s.repo.Create(ctx, userID, placeID, from, to)
}

func (s *service) List(ctx context.Context, userID *int, placeID *int) ([]models.Booking, error) {
	if userID != nil {
		return s.repo.ListByUser(ctx, *userID)
	}
	return s.repo.ListByPlace(ctx, *placeID)
}
