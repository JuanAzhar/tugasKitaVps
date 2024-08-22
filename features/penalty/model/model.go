package model

import (
	"time"

	"github.com/google/uuid"
)

type Penalty struct {
	Id          uuid.UUID `gorm:"type:varchar(50);primaryKey;not null" json:"id"`
	UserId      string
	Point       int
	Description string
	Date        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}
