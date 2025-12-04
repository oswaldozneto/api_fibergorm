package models

import (
	"api_fibergorm/pkg/arquitetura/entity"
)

// Produto representa a entidade de produto no banco de dados
type Produto struct {
	entity.BaseEntity
	Codigo    string  `gorm:"type:varchar(50);uniqueIndex;not null" json:"codigo"`
	Descricao string  `gorm:"type:varchar(255);not null" json:"descricao"`
	Preco     float64 `gorm:"type:decimal(10,2);not null" json:"preco"`

	// Chave estrangeira para Categoria
	CategoriaID uint      `gorm:"not null" json:"categoria_id"`
	Categoria   Categoria `gorm:"foreignKey:CategoriaID" json:"categoria,omitempty"`
}

// TableName define o nome da tabela no banco de dados
func (Produto) TableName() string {
	return "produtos"
}
