package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/go-gosh/gestful/component/mapper"
	"gorm.io/gorm"
)

type PageRequest interface {
	MakePage() mapper.Paginator
	MakeWrapper() func(*gorm.DB) *gorm.DB
}

type CreateRequest[T any] interface {
	MakeCreate() (*T, error)
}

type UpdateRequest interface {
	MakeUpdate() (map[string]interface{}, error)
}

type BaseService[T any, U CreateRequest[T], V PageRequest, W UpdateRequest] struct {
	mapper mapper.Mapper[T]
}

func handleErrorAdapter(handler func(*gin.Context) error) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := handler(ctx)
		if err != nil {
			ctx.JSON(500, err)
			return
		}
		ctx.JSON(200, "success")
	}
}

func (s BaseService[T, U, V, W]) RegisterGroupRoute(group *gin.RouterGroup, source string) {
	group.GET(fmt.Sprintf("/%s", source), func(ctx *gin.Context) {
		res, err := s.Paginate(ctx)
		if err != nil {
			ctx.JSON(500, err)
			return
		}
		ctx.JSON(200, res)
	})
	group.POST(fmt.Sprintf("/%s", source), handleErrorAdapter(s.Create))
	group.GET(fmt.Sprintf("/%s/:id", source), func(ctx *gin.Context) {
		res, err := s.Retrieve(ctx)
		if err != nil {
			ctx.JSON(500, err)
			return
		}
		ctx.JSON(200, res)
	})
	group.PUT(fmt.Sprintf("/%s/:id", source), handleErrorAdapter(s.Update))
	group.DELETE(fmt.Sprintf("/%s/:id", source), handleErrorAdapter(s.Delete))
}

func (s BaseService[T, U, V, W]) Create(ctx *gin.Context) error {
	var req U
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return err
	}
	create, err := req.MakeCreate()
	if err != nil {
		return err
	}
	if err := s.mapper.Create(ctx, create); err != nil {
		return err
	}

	return nil
}

func (s BaseService[T, U, V, W]) Paginate(ctx *gin.Context) (*mapper.PageRes[T], error) {
	var req V
	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		return nil, err
	}

	res, err := s.mapper.Paginate(ctx, req.MakePage(), req.MakeWrapper())
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s BaseService[T, U, V, W]) Retrieve(ctx *gin.Context) (*T, error) {
	id := struct {
		ID uint `uri:"id"`
	}{}
	if err := ctx.ShouldBindUri(&id); err != nil {
		return nil, err
	}

	return s.mapper.OneById(ctx, id.ID)
}

func (s BaseService[T, U, V, W]) Update(ctx *gin.Context) error {
	id := struct {
		ID uint `uri:"id"`
	}{}
	if err := ctx.ShouldBindUri(&id); err != nil {
		return err
	}
	var req W
	if err := ctx.ShouldBindJSON(&req); err != nil {
		return err
	}

	updated, err := req.MakeUpdate()
	if err != nil {
		return err
	}

	return s.mapper.UpdateById(ctx, id.ID, updated)
}

func (s BaseService[T, U, V, W]) Delete(ctx *gin.Context) error {
	id := struct {
		ID uint `uri:"id"`
	}{}
	if err := ctx.ShouldBindUri(&id); err != nil {
		return err
	}

	return s.mapper.DeleteById(ctx, id.ID)
}
