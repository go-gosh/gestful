package mapper

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type _testCRUDMapper struct {
	suite.Suite
	db     *gorm.DB
	mapper CRUDMapper[_testFoo]
}

func (t *_testCRUDMapper) SetupTest() {
	var err error
	t.db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	t.Require().NoError(err)
	t.db = t.db.Debug()
	t.Require().NoError(t.db.AutoMigrate(&_testFoo{}))
	t.mapper = NewCRUDMapper[_testFoo](t.db)
}

func (t *_testCRUDMapper) TearDownTest() {
	db, err := t.db.DB()
	t.Require().NoError(err)
	t.Require().NoError(db.Close())
}

func (t *_testCRUDMapper) Test_Paginate_AllData() {
	data := t.addData(30)
	ctx := context.TODO()
	res, err := t.mapper.Paginate(ctx, CRUDPaginator{}, EmptyWrapperFunc)
	t.NoError(err)
	t.NotNil(res.Data)
	t.EqualValues(0, res.Page)
	t.EqualValues(0, res.PageSize)
	t.EqualValues(30, res.Total)
	t.EqualValues(0, res.TotalPage)
	t.Len(res.Data, 30)
	for i, v := range res.Data {
		as, err := json.Marshal(v)
		t.NoError(err)
		es, err := json.Marshal(data[i])
		t.NoError(err)
		t.EqualValues(es, as)
	}
}

func (t *_testCRUDMapper) Test_Paginate_FirstPageNoData() {
	ctx := context.TODO()
	res, err := t.mapper.Paginate(ctx, CRUDPaginator{
		PageSize: 10,
	}, EmptyWrapperFunc)
	t.NoError(err)
	t.NotNil(res.Data)
	t.EqualValues(1, res.Page)
	t.EqualValues(10, res.PageSize)
	t.EqualValues(0, res.Total)
	t.EqualValues(0, res.TotalPage)
	t.Len(res.Data, 0)
}

func (t *_testCRUDMapper) Test_Paginate_OverPageHasData() {
	t.addData(101)
	ctx := context.TODO()
	res, err := t.mapper.Paginate(ctx, CRUDPaginator{
		Page:     12,
		PageSize: 10,
	}, EmptyWrapperFunc)
	t.NoError(err)
	t.NotNil(res.Data)
	t.EqualValues(12, res.Page)
	t.EqualValues(10, res.PageSize)
	t.EqualValues(101, res.Total)
	t.EqualValues(11, res.TotalPage)
	t.Len(res.Data, 0)
}

func (t *_testCRUDMapper) Test_Paginate_ZeroPageNoMoreData() {
	data := t.addData(10)
	ctx := context.TODO()
	res, err := t.mapper.Paginate(ctx, CRUDPaginator{
		PageSize: 10,
	}, EmptyWrapperFunc)
	t.NoError(err)
	t.NotNil(res.Data)
	t.EqualValues(1, res.Page)
	t.EqualValues(10, res.PageSize)
	t.EqualValues(10, res.Total)
	t.EqualValues(1, res.TotalPage)
	t.Len(res.Data, 10)
	for i, v := range res.Data {
		as, err := json.Marshal(v)
		t.NoError(err)
		es, err := json.Marshal(data[i])
		t.NoError(err)
		t.EqualValues(es, as)
	}
}

func (t *_testCRUDMapper) Test_Paginate_FirstPageMoreData() {
	data := t.addData(30)
	ctx := context.TODO()
	res, err := t.mapper.Paginate(ctx, CRUDPaginator{
		Page:     1,
		PageSize: 10,
	}, EmptyWrapperFunc)
	t.NoError(err)
	t.NotNil(res.Data)
	t.EqualValues(1, res.Page)
	t.EqualValues(10, res.PageSize)
	t.EqualValues(30, res.Total)
	t.EqualValues(3, res.TotalPage)
	t.Len(res.Data, 10)
	for i, v := range res.Data {
		as, err := json.Marshal(v)
		t.NoError(err)
		es, err := json.Marshal(data[i])
		t.NoError(err)
		t.EqualValues(es, as)
	}
}

func (t *_testCRUDMapper) Test_Paginate_LastPageMoreData() {
	data := t.addData(32)
	ctx := context.TODO()
	res, err := t.mapper.Paginate(ctx, CRUDPaginator{
		Page:     4,
		PageSize: 10,
	}, EmptyWrapperFunc)
	t.NoError(err)
	t.NotNil(res.Data)
	t.EqualValues(4, res.Page)
	t.EqualValues(10, res.PageSize)
	t.EqualValues(32, res.Total)
	t.EqualValues(4, res.TotalPage)
	t.Len(res.Data, 2)
	for i, v := range res.Data {
		as, err := json.Marshal(v)
		t.NoError(err)
		es, err := json.Marshal(data[i+30])
		t.NoError(err)
		t.EqualValues(es, as)
	}
}

func (t *_testCRUDMapper) addData(num int) []_testFoo {
	res := make([]_testFoo, 0, num)
	for i := 0; i < num; i++ {
		foo := _testFoo{}
		t.Require().NoError(t.db.Create(&foo).Error)
		res = append(res, foo)
	}
	return res
}

func TestCRUDMapper(t *testing.T) {
	suite.Run(t, &_testCRUDMapper{})
}
