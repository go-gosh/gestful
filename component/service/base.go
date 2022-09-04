package service

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-gosh/gestful/component/mapper"
	"gorm.io/gorm"
)

type PageRequest interface {
	MakePage() mapper.Paginator
	MakeWrapper() func(*gorm.DB) *gorm.DB
}

type BasePageRequest struct {
	mapper.Paginator
}

const DefaultPageLimit = 10
const MaxPageLimit = 500

func (b BasePageRequest) MakePage() mapper.Paginator {
	if b.Limit <= 0 {
		b.Limit = DefaultPageLimit
	}
	if b.Limit > MaxPageLimit {
		b.Limit = MaxPageLimit
	}

	return b.Paginator
}

func (b BasePageRequest) MakeWrapper() func(*gorm.DB) *gorm.DB {
	return mapper.EmptyWrapperFunc
}

type CreateRequest[T any] interface {
	MakeCreate() (*T, error)
}

type BaseCreateRequest[T any] struct {
	Data T `json:"data"`
}

func (b BaseCreateRequest[T]) MakeCreate() (*T, error) {
	return &b.Data, nil
}

type UpdateRequest interface {
	MakeUpdate() (map[string]interface{}, error)
}

type BaseUpdateRequest struct {
	Data map[string]interface{} `json:"data"`
}

func (b BaseUpdateRequest) MakeUpdate() (map[string]interface{}, error) {
	if _, ok := b.Data["id"]; ok {
		delete(b.Data, "id")
	}
	return b.Data, nil
}

type BaseRestfulService[T any] interface {
	RegisterGroupRoute(group *gin.RouterGroup, source string)
	Create(ctx *gin.Context) error
	Paginate(ctx *gin.Context) (*mapper.PageRes[T], error)
	Retrieve(ctx *gin.Context) (*T, error)
	Update(ctx *gin.Context) error
	Delete(ctx *gin.Context) error
}

// NewBaseService new base restful service
func NewBaseService[T any, U CreateRequest[T], V PageRequest, W UpdateRequest](mapper mapper.BaseMapper[T]) BaseRestfulService[T] {
	return &baseService[T, U, V, W]{mapper: mapper}
}

type baseService[T any, U CreateRequest[T], V PageRequest, W UpdateRequest] struct {
	mapper mapper.BaseMapper[T]
}

func handleErrorAdapter(handler func(*gin.Context) error) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		err := handler(ctx)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, "not found")
			return
		}
		if err != nil {
			ctx.AbortWithStatusJSON(500, err.Error())
			return
		}
		ctx.JSON(200, "success")
	}
}

func (s baseService[T, U, V, W]) RegisterGroupRoute(group *gin.RouterGroup, source string) {
	group.GET(fmt.Sprintf("/%s", source), func(ctx *gin.Context) {
		res, err := s.Paginate(ctx)
		if err != nil {
			ctx.AbortWithStatusJSON(500, err.Error())
			return
		}
		ctx.JSON(200, res)
	})
	group.POST(fmt.Sprintf("/%s", source), handleErrorAdapter(s.Create))
	group.GET(fmt.Sprintf("/%s/:id", source), func(ctx *gin.Context) {
		res, err := s.Retrieve(ctx)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.AbortWithStatusJSON(http.StatusNotFound, "not found")
			return
		}
		if err != nil {
			ctx.AbortWithStatusJSON(500, err.Error())
			return
		}
		ctx.JSON(200, res)
	})
	group.PUT(fmt.Sprintf("/%s/:id", source), handleErrorAdapter(s.Update))
	group.DELETE(fmt.Sprintf("/%s/:id", source), handleErrorAdapter(s.Delete))
}

func (s baseService[T, U, V, W]) Create(ctx *gin.Context) error {
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

func (s baseService[T, U, V, W]) Paginate(ctx *gin.Context) (*mapper.PageRes[T], error) {
	var req V
	err := ctx.ShouldBindQuery(&req)
	if err != nil {
		return nil, err
	}

	res, err := s.mapper.Paginate(ctx, req.MakePage(), req.MakeWrapper())
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (s baseService[T, U, V, W]) Retrieve(ctx *gin.Context) (*T, error) {
	id := struct {
		ID uint `uri:"id"`
	}{}
	if err := ctx.ShouldBindUri(&id); err != nil {
		return nil, err
	}

	return s.mapper.OneById(ctx, id.ID)
}

func (s baseService[T, U, V, W]) Update(ctx *gin.Context) error {
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

func (s baseService[T, U, V, W]) Delete(ctx *gin.Context) error {
	id := struct {
		ID uint `uri:"id"`
	}{}
	if err := ctx.ShouldBindUri(&id); err != nil {
		return err
	}

	return s.mapper.DeleteById(ctx, id.ID)
}
