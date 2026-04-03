package models

import "time"

type Booking struct {
	ID      int       `gorm:"column:id;primaryKey;autoIncrement"`
	UserID  int       `gorm:"column:user_id"`
	PlaceID int       `gorm:"column:place_id"`
	From    time.Time `gorm:"column:time_from"`
	To      time.Time `gorm:"column:time_to"`
}

func (Booking) TableName() string {
	return "bookings"
}
