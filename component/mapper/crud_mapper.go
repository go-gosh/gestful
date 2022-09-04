package mapper

import (
	"context"

	"gorm.io/gorm"
)

type CRUDPaginator struct {
	Page     uint `json:"page" form:"page"`
	PageSize uint `json:"page_size" form:"page_size"`
}

type CRUDPageResult[Model any] struct {
	CRUDPaginator
	Total     uint    `json:"total"`
	TotalPage uint    `json:"total_page"`
	Data      []Model `json:"data"`
}

type CRUDMapper[Model any] interface {
	IMapper[Model, CRUDPaginator, CRUDPageResult[Model]]
}

func NewCRUDMapper[Model any](db *gorm.DB) CRUDMapper[Model] {
	return &crudMapper[Model]{
		db:     db,
		mapper: NewBaseMapper[Model](db),
	}
}

type crudMapper[Model any] struct {
	db     *gorm.DB
	mapper BaseMapper[Model]
}

func (c *crudMapper[Model]) One(ctx context.Context, wrapper func(*gorm.DB) *gorm.DB) (*Model, error) {
	return c.mapper.One(ctx, wrapper)
}

func (c *crudMapper[Model]) OneById(ctx context.Context, id uint) (*Model, error) {

	return c.mapper.OneById(ctx, id)
}

func (c *crudMapper[Model]) All(ctx context.Context, wrapper func(*gorm.DB) *gorm.DB) ([]Model, error) {
	return c.mapper.All(ctx, wrapper)
}

func (c *crudMapper[Model]) Paginate(ctx context.Context, pager CRUDPaginator, wrapper func(*gorm.DB) *gorm.DB) (*CRUDPageResult[Model], error) {
	if pager.PageSize == 0 {
		data, err := c.All(ctx, wrapper)
		if err != nil {
			return nil, err
		}
		total, err := c.Count(ctx, wrapper)
		if err != nil {
			return nil, err
		}
		return &CRUDPageResult[Model]{
			CRUDPaginator: pager,
			Total:         uint(total),
			TotalPage:     0,
			Data:          data,
		}, nil
	}
	data := make([]Model, 0, pager.PageSize)
	db := wrapper(c.db)
	if pager.Page == 0 {
		pager.Page = 1
	}
	offset := int((pager.Page - 1) * pager.PageSize)
	err := db.Offset(offset).Limit(int(pager.PageSize)).Find(&data).Error
	if err != nil {
		return nil, err
	}
	res := CRUDPageResult[Model]{
		CRUDPaginator: pager,
		Data:          data,
	}
	total, err := c.Count(ctx, wrapper)
	if err != nil {
		return nil, err
	}
	res.Total = uint(total)
	if res.Total != 0 {
		res.TotalPage = (res.Total-1)/res.PageSize + 1
	}

	return &res, nil
}

func (c *crudMapper[Model]) Count(ctx context.Context, wrapper func(*gorm.DB) *gorm.DB) (int, error) {
	return c.mapper.Count(ctx, wrapper)
}

func (c *crudMapper[Model]) Delete(ctx context.Context, wrapper func(*gorm.DB) *gorm.DB) error {

	return c.mapper.Delete(ctx, wrapper)
}

func (c *crudMapper[Model]) DeleteById(ctx context.Context, id uint) error {
	return c.mapper.DeleteById(ctx, id)
}

func (c *crudMapper[Model]) Create(ctx context.Context, entity *Model) error {
	return c.mapper.Create(ctx, entity)
}

func (c *crudMapper[Model]) Update(ctx context.Context, wrapper func(*gorm.DB) *gorm.DB, updated map[string]interface{}) error {

	return c.mapper.Update(ctx, wrapper, updated)
}

func (c *crudMapper[Model]) UpdateById(ctx context.Context, id uint, updated map[string]interface{}) error {
	return c.mapper.UpdateById(ctx, id, updated)
}
