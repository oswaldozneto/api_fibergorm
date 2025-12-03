package models

import (
	"time"

	"gorm.io/gorm"
)

// Categoria representa a entidade de categoria no banco de dados
type Categoria struct {
	ID        uint           `gorm:"primaryKey" json:"id"`
	Nome      string         `gorm:"type:varchar(100);uniqueIndex;not null" json:"nome"`
	Descricao string         `gorm:"type:varchar(255)" json:"descricao"`
	Ativo     bool           `gorm:"default:true" json:"ativo"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at,omitempty"`

	// Relacionamento: Uma categoria tem muitos produtos
	Produtos []Produto `gorm:"foreignKey:CategoriaID" json:"produtos,omitempty"`
}

// TableName define o nome da tabela no banco de dados
func (Categoria) TableName() string {
	return "categorias"
}
