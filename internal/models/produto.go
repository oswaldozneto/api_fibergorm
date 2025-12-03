package models

import (
	"time"

	"gorm.io/gorm"
)

// Produto representa a entidade de produto no banco de dados
type Produto struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Codigo    string         `gorm:"type:varchar(50);uniqueIndex;not null" json:"codigo"`
	Descricao string         `gorm:"type:varchar(255);not null" json:"descricao"`
	Preco     float64        `gorm:"type:decimal(10,2);not null" json:"preco"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`
}

// TableName define o nome da tabela no banco de dados
func (Produto) TableName() string {
	return "produtos"
}
