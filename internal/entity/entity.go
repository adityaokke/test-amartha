package entity

import (
	"time"
)

type BaseTimeStruct struct {
	// default value to prevent gorm automatic set to current time on code level
	CreatedAt time.Time  `json:"createdAt" gorm:"type:TIMESTAMP;default:CURRENT_TIMESTAMP"`
	UpdatedAt *time.Time `json:"updatedAt" gorm:"type:TIMESTAMP;default:NULL"`
}
