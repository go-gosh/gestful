package mapper

import (
	"context"

	"gorm.io/gorm"
)

type Paginator struct {
	StartId uint `json:"start_id"`
	Limit   int  `json:"limit"`
}

type PageRes[T any] struct {
	Paginator
	More bool `json:"more"`
	Data []T  `json:"data"`
}

type Mapper[T any] interface {
	One(ctx context.Context, wrapper func(*gorm.DB) *gorm.DB) (*T, error)
	OneById(ctx context.Context, id uint) (*T, error)
	All(ctx context.Context, wrapper func(*gorm.DB) *gorm.DB) ([]T, error)
	Paginate(ctx context.Context, pager Paginator, wrapper func(*gorm.DB) *gorm.DB) (*PageRes[T], error)
	Delete(ctx context.Context, wrapper func(*gorm.DB) *gorm.DB) error
	DeleteById(ctx context.Context, id uint) error
	Create(ctx context.Context, entity *T) error
	Update(ctx context.Context, wrapper func(*gorm.DB) *gorm.DB, updated map[string]interface{}) error
	UpdateById(ctx context.Context, id uint, updated map[string]interface{}) error
}

type baseMapper[T any] struct {
	db *gorm.DB
}

func WrapperFuncById(id uint) func(db *gorm.DB) *gorm.DB {
	return func(db *gorm.DB) *gorm.DB {
		return db.Where("id=?", id)
	}
}

func EmptyWrapperFunc(db *gorm.DB) *gorm.DB {
	return db
}

// NewBaseMapper base mapper
func NewBaseMapper[T any](db *gorm.DB) Mapper[T] {
	return &baseMapper[T]{db: db}
}

func (m baseMapper[T]) OneById(ctx context.Context, id uint) (*T, error) {
	return m.One(ctx, WrapperFuncById(id))
}

func (m baseMapper[T]) DeleteById(ctx context.Context, id uint) error {
	return m.Delete(ctx, WrapperFuncById(id))
}

func (m baseMapper[T]) UpdateById(ctx context.Context, id uint, updated map[string]interface{}) error {
	return m.Update(ctx, WrapperFuncById(id), updated)
}

func (m baseMapper[T]) One(ctx context.Context, wrapper func(*gorm.DB) *gorm.DB) (*T, error) {
	var res T
	err := wrapper(m.db.WithContext(ctx)).
		First(&res).Error
	return &res, err
}

func (m baseMapper[T]) All(ctx context.Context, wrapper func(*gorm.DB) *gorm.DB) ([]T, error) {
	res := make([]T, 0)
	err := wrapper(m.db.WithContext(ctx)).
		Find(&res).Error
	return res, err
}

func (m baseMapper[T]) Paginate(ctx context.Context, pager Paginator, wrapper func(*gorm.DB) *gorm.DB) (*PageRes[T], error) {
	res := make([]T, 0, pager.Limit+1)
	db := wrapper(m.db.WithContext(ctx))
	if pager.StartId > 0 {
		db = db.
			Where("id>?", pager.StartId)
	}
	err := db.
		Limit(pager.Limit + 1).
		Find(&res).Error
	if err != nil {
		return nil, err
	}
	more := len(res) > pager.Limit
	if more {
		res = res[:pager.Limit]
	}

	return &PageRes[T]{
		Paginator: pager,
		More:      more,
		Data:      res,
	}, nil
}

func (m baseMapper[T]) Delete(ctx context.Context, wrapper func(*gorm.DB) *gorm.DB) error {
	var t T
	return wrapper(m.db.WithContext(ctx)).Delete(&t).Error
}

func (m baseMapper[T]) Create(ctx context.Context, entity *T) error {
	return m.db.WithContext(ctx).Create(entity).Error
}

func (m baseMapper[T]) Update(ctx context.Context, wrapper func(*gorm.DB) *gorm.DB, updated map[string]interface{}) error {
	return wrapper(m.db.WithContext(ctx)).Updates(updated).Error
}
