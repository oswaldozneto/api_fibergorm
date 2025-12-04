package repository

import (
	"errors"
	"reflect"

	"api_fibergorm/pkg/arquitetura/entity"
	arqerrors "api_fibergorm/pkg/arquitetura/errors"

	"gorm.io/gorm"
)

// BaseRepositoryImpl é a implementação base do repositório genérico
// E é o tipo ponteiro da entidade que implementa entity.Entity (ex: *models.Categoria)
type BaseRepositoryImpl[E entity.Entity] struct {
	db           *gorm.DB
	preloads     []string
	defaultOrder string
}

// NewBaseRepository cria uma nova instância do repositório base
func NewBaseRepository[E entity.Entity](db *gorm.DB) *BaseRepositoryImpl[E] {
	return &BaseRepositoryImpl[E]{
		db:           db,
		preloads:     []string{},
		defaultOrder: "id ASC",
	}
}

// WithPreloads configura os preloads padrão (retorna o próprio repositório para chaining)
func (r *BaseRepositoryImpl[E]) WithPreloads(preloads ...string) *BaseRepositoryImpl[E] {
	r.preloads = preloads
	return r
}

// WithDefaultOrder configura a ordenação padrão (retorna o próprio repositório para chaining)
func (r *BaseRepositoryImpl[E]) WithDefaultOrder(order string) *BaseRepositoryImpl[E] {
	r.defaultOrder = order
	return r
}

// GetDB retorna a instância do banco de dados
func (r *BaseRepositoryImpl[E]) GetDB() *gorm.DB {
	return r.db
}

// newEntity cria uma nova instância da entidade usando reflection
// Como E é um tipo ponteiro (ex: *models.Categoria), precisamos criar a struct subjacente
func (r *BaseRepositoryImpl[E]) newEntity() E {
	var zero E
	t := reflect.TypeOf(zero)
	if t.Kind() == reflect.Ptr {
		return reflect.New(t.Elem()).Interface().(E)
	}
	return zero
}

// Create insere uma nova entidade no banco de dados
func (r *BaseRepositoryImpl[E]) Create(entity E) error {
	return r.db.Create(entity).Error
}

// FindByID busca uma entidade pelo ID
func (r *BaseRepositoryImpl[E]) FindByID(id uint) (E, error) {
	entity := r.newEntity()
	query := r.db

	// Aplica preloads se configurados
	for _, preload := range r.preloads {
		query = query.Preload(preload)
	}

	err := query.First(entity, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, arqerrors.ErrNotFound
		}
		return entity, err
	}
	return entity, nil
}

// FindByIDWithPreloads busca uma entidade pelo ID com preloads específicos
func (r *BaseRepositoryImpl[E]) FindByIDWithPreloads(id uint, preloads ...string) (E, error) {
	entity := r.newEntity()
	query := r.db

	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	err := query.First(entity, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, arqerrors.ErrNotFound
		}
		return entity, err
	}
	return entity, nil
}

// FindAll retorna todas as entidades com paginação
func (r *BaseRepositoryImpl[E]) FindAll(page, pageSize int, orderBy string) ([]E, int64, error) {
	var entities []E
	var total int64

	// Conta o total de registros
	if err := r.db.Model(r.newEntity()).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calcula o offset
	offset := (page - 1) * pageSize

	// Define a ordenação
	order := orderBy
	if order == "" {
		order = r.defaultOrder
	}

	// Aplica preloads se configurados
	query := r.db
	for _, preload := range r.preloads {
		query = query.Preload(preload)
	}

	// Busca com paginação
	err := query.Offset(offset).Limit(pageSize).Order(order).Find(&entities).Error
	if err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

// FindAllWithPreloads retorna todas as entidades com preloads específicos
func (r *BaseRepositoryImpl[E]) FindAllWithPreloads(page, pageSize int, orderBy string, preloads ...string) ([]E, int64, error) {
	var entities []E
	var total int64

	if err := r.db.Model(r.newEntity()).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	order := orderBy
	if order == "" {
		order = r.defaultOrder
	}

	query := r.db
	for _, preload := range preloads {
		query = query.Preload(preload)
	}

	err := query.Offset(offset).Limit(pageSize).Order(order).Find(&entities).Error
	if err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

// FindAllWhere busca entidades com condições
func (r *BaseRepositoryImpl[E]) FindAllWhere(page, pageSize int, orderBy string, condition interface{}, args ...interface{}) ([]E, int64, error) {
	var entities []E
	var total int64

	if err := r.db.Model(r.newEntity()).Where(condition, args...).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * pageSize

	order := orderBy
	if order == "" {
		order = r.defaultOrder
	}

	query := r.db.Where(condition, args...)
	for _, preload := range r.preloads {
		query = query.Preload(preload)
	}

	err := query.Offset(offset).Limit(pageSize).Order(order).Find(&entities).Error
	if err != nil {
		return nil, 0, err
	}

	return entities, total, nil
}

// FindOneWhere busca uma entidade com condição
func (r *BaseRepositoryImpl[E]) FindOneWhere(condition interface{}, args ...interface{}) (E, error) {
	entity := r.newEntity()
	query := r.db.Where(condition, args...)

	for _, preload := range r.preloads {
		query = query.Preload(preload)
	}

	err := query.First(entity).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, arqerrors.ErrNotFound
		}
		return entity, err
	}
	return entity, nil
}

// Update atualiza uma entidade existente
func (r *BaseRepositoryImpl[E]) Update(entity E) error {
	return r.db.Save(entity).Error
}

// Delete remove uma entidade pelo ID (soft delete se configurado)
func (r *BaseRepositoryImpl[E]) Delete(id uint) error {
	result := r.db.Delete(r.newEntity(), id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return arqerrors.ErrNotFound
	}
	return nil
}

// ExistsByID verifica se existe uma entidade com o ID informado
func (r *BaseRepositoryImpl[E]) ExistsByID(id uint) (bool, error) {
	var count int64
	err := r.db.Model(r.newEntity()).Where("id = ?", id).Count(&count).Error
	return count > 0, err
}

// ExistsWhere verifica se existe uma entidade com a condição informada
func (r *BaseRepositoryImpl[E]) ExistsWhere(condition interface{}, args ...interface{}) (bool, error) {
	var count int64
	err := r.db.Model(r.newEntity()).Where(condition, args...).Count(&count).Error
	return count > 0, err
}

// ExistsWhereExcludingID verifica se existe outra entidade com a condição (excluindo o ID)
func (r *BaseRepositoryImpl[E]) ExistsWhereExcludingID(id uint, condition string, args ...interface{}) (bool, error) {
	var count int64
	fullCondition := condition + " AND id != ?"
	fullArgs := append(args, id)
	err := r.db.Model(r.newEntity()).Where(fullCondition, fullArgs...).Count(&count).Error
	return count > 0, err
}

// CountWhere conta registros com condição
func (r *BaseRepositoryImpl[E]) CountWhere(condition interface{}, args ...interface{}) (int64, error) {
	var count int64
	err := r.db.Model(r.newEntity()).Where(condition, args...).Count(&count).Error
	return count, err
}
