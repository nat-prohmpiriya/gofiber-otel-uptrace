// internal/domain/base.go
package domain

import (
	"time"
)

// Base contains common fields that should be embedded in all domain entities
type Base struct {
	ID        string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt time.Time `json:"deleted_at,omitempty" gorm:"index"`
}

// IsDeleted checks if the entity is soft deleted
func (b *Base) IsDeleted() bool {
	return !b.DeletedAt.IsZero()
}

// MarkAsDeleted marks the entity as deleted
func (b *Base) MarkAsDeleted() {
	b.DeletedAt = time.Now()
}

// Timestamps returns created and updated times
func (b *Base) Timestamps() (created, updated time.Time) {
	return b.CreatedAt, b.UpdatedAt
}
