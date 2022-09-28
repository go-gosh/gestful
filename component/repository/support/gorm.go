package support

import (
	"github.com/go-gosh/gestful/component/domain"
	"gorm.io/gorm"
)

type GormJpaRepository[T, ID any] struct {
	*gorm.DB
}

func (g GormJpaRepository[T, ID]) Save(entity *T) (*T, error) {
	err := g.DB.Create(entity).Error
	return entity, err
}

func (g GormJpaRepository[T, ID]) SaveAll(entity ...*T) ([]*T, error) {
	err := g.DB.Create(&entity).Error
	return entity, err
}

func (g GormJpaRepository[T, ID]) FindById(id ID) (*T, error) {
	var entity T
	err := g.DB.Where("id=?", id).First(&entity).Error
	return &entity, err
}

func (g GormJpaRepository[T, ID]) ExistsById(id ID) (bool, error) {
	_, err := g.FindById(id)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (g GormJpaRepository[T, ID]) FindAll() ([]T, error) {
	res := make([]T, 0)
	err := g.DB.Find(&res).Error
	return res, err
}

func (g GormJpaRepository[T, ID]) FindAllById(id ...ID) ([]T, error) {
	res := make([]T, 0)
	err := g.DB.Where("id in ?", id).Find(&res).Error
	return res, err
}

func (g GormJpaRepository[T, ID]) Count() (int, error) {
	var c int64
	err := g.DB.Count(&c).Error
	return int(c), err
}

func (g GormJpaRepository[T, ID]) DeleteById(id ID) error {
	var t T
	return g.DB.Model(&t).Where("id=?", id).Delete(&t).Error
}

func (g GormJpaRepository[T, ID]) Delete(entity T) error {
	return g.DB.Delete(&entity).Error
}

func (g GormJpaRepository[T, ID]) DeleteAllById(id ...ID) error {
	var t T
	return g.DB.Model(&t).Where("id in ?", id).Delete(&t).Error
}

func (g GormJpaRepository[T, ID]) DeleteAll(entity ...T) error {
	return g.DB.Transaction(func(tx *gorm.DB) error {
		for i := 0; i < len(entity); i++ {
			err := tx.Delete(&entity[i]).Error
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (g GormJpaRepository[T, ID]) FindAllBySort(sort domain.Sort) ([]T, error) {
	return g.FindAll()
}

func (g GormJpaRepository[T, ID]) FindAllByPage(page domain.Pageable) (domain.Page[T], error) {
	r, err := g.FindAllBySort(page.GetSort())
	if err != nil {
		return nil, err
	}
	total, err := g.Count()
	if err != nil {
		return nil, err
	}
	return domain.NewPage(total, page, r), nil
}
