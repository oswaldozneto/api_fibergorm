package database

import (
	"api_fibergorm/internal/models"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

// Seed executa a carga inicial de dados necessários
func Seed(db *gorm.DB, log *logrus.Logger) error {
	log.Info("Executando seed de dados iniciais")

	// Cria categoria padrão se não existir
	if err := seedCategoriaPadrao(db, log); err != nil {
		return err
	}

	// Atualiza produtos órfãos (sem categoria) para usar a categoria padrão
	if err := fixProdutosSemCategoria(db, log); err != nil {
		return err
	}

	log.Info("Seed de dados iniciais concluído")
	return nil
}

// seedCategoriaPadrao cria a categoria padrão se não existir
func seedCategoriaPadrao(db *gorm.DB, log *logrus.Logger) error {
	var count int64
	db.Model(&models.Categoria{}).Where("nome = ?", "Geral").Count(&count)

	if count == 0 {
		categoriaPadrao := &models.Categoria{
			Nome:      "Geral",
			Descricao: "Categoria padrão para produtos sem categoria definida",
			Ativo:     true,
		}

		if err := db.Create(categoriaPadrao).Error; err != nil {
			log.WithError(err).Error("Falha ao criar categoria padrão")
			return err
		}

		log.WithField("id", categoriaPadrao.ID).Info("Categoria padrão 'Geral' criada com sucesso")
	} else {
		log.Debug("Categoria padrão 'Geral' já existe")
	}

	return nil
}

// fixProdutosSemCategoria atualiza produtos sem categoria para usar a categoria padrão
func fixProdutosSemCategoria(db *gorm.DB, log *logrus.Logger) error {
	// Busca a categoria padrão
	var categoriaPadrao models.Categoria
	if err := db.Where("nome = ?", "Geral").First(&categoriaPadrao).Error; err != nil {
		log.WithError(err).Error("Categoria padrão não encontrada")
		return err
	}

	// Conta quantos produtos estão sem categoria (categoria_id = 0 ou NULL)
	var count int64
	db.Model(&models.Produto{}).Where("categoria_id = 0 OR categoria_id IS NULL").Count(&count)

	if count > 0 {
		// Atualiza todos os produtos sem categoria
		result := db.Model(&models.Produto{}).
			Where("categoria_id = 0 OR categoria_id IS NULL").
			Update("categoria_id", categoriaPadrao.ID)

		if result.Error != nil {
			log.WithError(result.Error).Error("Falha ao atualizar produtos sem categoria")
			return result.Error
		}

		log.WithFields(logrus.Fields{
			"quantidade":   result.RowsAffected,
			"categoria_id": categoriaPadrao.ID,
		}).Info("Produtos atualizados para categoria padrão")
	} else {
		log.Debug("Nenhum produto sem categoria encontrado")
	}

	return nil
}
