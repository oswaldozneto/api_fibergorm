package models

import (
	"api_fibergorm/pkg/arquitetura/entity"
)

// Categoria representa a entidade de categoria no banco de dados
type Categoria struct {
	entity.BaseEntity
	Nome      string `gorm:"type:varchar(100);uniqueIndex;not null" json:"nome"`
	Descricao string `gorm:"type:varchar(255)" json:"descricao"`
	Ativo     bool   `gorm:"default:true" json:"ativo"`

	// Relacionamento: Uma categoria tem muitos produtos
	Produtos []Produto `gorm:"foreignKey:CategoriaID" json:"produtos,omitempty"`
}

// TableName define o nome da tabela no banco de dados
func (Categoria) TableName() string {
	return "categorias"
}
