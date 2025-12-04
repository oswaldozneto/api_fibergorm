package entity

import (
	"time"

	"gorm.io/gorm"
)

// Entity é a interface base que todas as entidades do sistema devem implementar
type Entity interface {
	GetID() uint
	SetID(id uint)
	GetCreatedAt() string
	GetUpdatedAt() string
	TableName() string
}

// BaseEntity contém os campos comuns a todas as entidades
// Deve ser embutida em todas as entidades do sistema
// NOTA: As entidades que embutem BaseEntity devem implementar TableName()
type BaseEntity struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// GetID retorna o ID da entidade
func (e *BaseEntity) GetID() uint {
	return e.ID
}

// SetID define o ID da entidade
func (e *BaseEntity) SetID(id uint) {
	e.ID = id
}

// GetCreatedAt retorna a data de criação formatada
func (e *BaseEntity) GetCreatedAt() string {
	return e.CreatedAt.Format("2006-01-02 15:04:05")
}

// GetUpdatedAt retorna a data de atualização formatada
func (e *BaseEntity) GetUpdatedAt() string {
	return e.UpdatedAt.Format("2006-01-02 15:04:05")
}

// TableName deve ser implementado pelas entidades que embutem BaseEntity
// Este método existe apenas para documentação - cada entidade DEVE implementar seu próprio TableName()
// func (e *BaseEntity) TableName() string { panic("TableName must be implemented by embedding entity") }
