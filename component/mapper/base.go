package mapper

import (
	"context"

	"gorm.io/gorm"
)

type Paginator struct {
	StartId uint `json:"start_id" form:"start_id"`
	Limit   int  `json:"limit" form:"limit"`
}

type PageRes[T any] struct {
	Paginator
	More bool `json:"more"`
	Data []T  `json:"data"`
}

type BaseMapper[T any] interface {
	IMapper[T, Paginator, PageRes[T]]
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
func NewBaseMapper[T any](db *gorm.DB) BaseMapper[T] {
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
	err := db.Order("id asc").
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
	d := wrapper(m.db.WithContext(ctx)).Delete(&t)
	if d.Error != nil {
		return d.Error
	}
	if d.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (m baseMapper[T]) Create(ctx context.Context, entity *T) error {
	return m.db.WithContext(ctx).Create(entity).Error
}

func (m baseMapper[T]) Update(ctx context.Context, wrapper func(*gorm.DB) *gorm.DB, updated map[string]interface{}) error {
	var t T
	updates := wrapper(m.db.WithContext(ctx).Model(&t)).Updates(updated)
	if updates.Error != nil {
		return updates.Error
	}
	if updates.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (m baseMapper[T]) Count(ctx context.Context, wrapper func(*gorm.DB) *gorm.DB) (int, error) {
	var c int64
	var t T
	err := wrapper(m.db.WithContext(ctx).Model(&t)).Count(&c).Error
	return int(c), err
}
