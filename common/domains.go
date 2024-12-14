package common

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type UUIDModel struct {
	// UUID
	ID        uuid.UUID `gorm:"type:uuid;default:null"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *UUIDModel) BeforeCreate(tx *gorm.DB) (err error) {
	u.ID = uuid.New()
	return
}
