package gorm

import (
	"time"
)

type BaseModel struct {
	ID        uint64     `gorm:"column:id;type:bigint UNSIGNED;primaryKey;not null;" json:"id"`
	CreatedAt *time.Time `gorm:"column:created_at;type:TIMESTAMP;autoCreateTime;" json:"created_at"`
	UpdatedAt *time.Time `gorm:"column:updated_at;type:TIMESTAMP;autoCreateTime;" json:"updated_at"`
	//DeletedAt gorm.DeletedAt `json:"-"`
	//IsDeleted bool           `json:"-"`
}
