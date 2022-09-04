package mapper

import (
	"context"
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/suite"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type _testFoo struct {
	ID uint `gorm:"primaryKey"`
}

type _testMapper struct {
	suite.Suite
	db     *gorm.DB
	mapper BaseMapper[_testFoo]
}

func (t *_testMapper) SetupTest() {
	var err error
	t.db, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	t.Require().NoError(err)
	t.db = t.db.Debug()
	t.Require().NoError(t.db.AutoMigrate(&_testFoo{}))
	t.mapper = NewBaseMapper[_testFoo](t.db)
}

func (t *_testMapper) TearDownTest() {
	db, err := t.db.DB()
	t.Require().NoError(err)
	t.Require().NoError(db.Close())
}

func (t *_testMapper) Test_OneById_FindData() {
	t.addData(10)
	ctx := context.TODO()
	foo, err := t.mapper.OneById(ctx, 10)
	t.NoError(err)
	t.NotNil(foo)
	t.EqualValues(10, foo.ID, foo)
}

func (t *_testMapper) Test_OneById_NotFindData() {
	t.addData(10)
	ctx := context.TODO()
	_, err := t.mapper.OneById(ctx, 11)
	t.ErrorIs(err, gorm.ErrRecordNotFound)
}

func (t *_testMapper) Test_DeleteById_FindData() {
	t.addData(10)
	ctx := context.TODO()
	err := t.mapper.DeleteById(ctx, 10)
	t.NoError(err)
	var foo _testFoo
	t.ErrorIs(t.db.First(&foo, 10).Error, gorm.ErrRecordNotFound)
}

func (t *_testMapper) Test_DeleteById_NotFindData() {
	t.addData(10)
	ctx := context.TODO()
	err := t.mapper.DeleteById(ctx, 11)
	t.ErrorIs(err, gorm.ErrRecordNotFound)
}

func (t *_testMapper) Test_UpdateById_FindData() {
	t.addData(10)
	ctx := context.TODO()
	err := t.mapper.UpdateById(ctx, 10, map[string]interface{}{"id": 11})
	t.NoError(err)
	t.ErrorIs(t.db.First(&_testFoo{}, 10).Error, gorm.ErrRecordNotFound)
	var foo _testFoo
	t.NoError(t.db.First(&foo, 11).Error)
	t.NotNil(foo)
	t.EqualValues(11, foo.ID)
}

func (t *_testMapper) Test_UpdateById_NotFindData() {
	t.addData(10)
	ctx := context.TODO()
	err := t.mapper.UpdateById(ctx, 11, map[string]interface{}{"id": 12})
	t.ErrorIs(err, gorm.ErrRecordNotFound)
	t.ErrorIs(t.db.First(&_testFoo{}, 12).Error, gorm.ErrRecordNotFound)
}

func (t *_testMapper) Test_Create() {
	t.addData(10)
	ctx := context.TODO()
	err := t.mapper.Create(ctx, &_testFoo{})
	t.NoError(err)
	var foo _testFoo
	t.NoError(t.db.First(&foo, 11).Error)
	t.NotNil(foo)
	t.EqualValues(11, foo.ID)
}

func (t *_testMapper) Test_All() {
	data := t.addData(20)
	ctx := context.TODO()
	res, err := t.mapper.All(ctx, EmptyWrapperFunc)
	t.NoError(err)
	t.Len(res, 20)
	for i, v := range res {
		as, err := json.Marshal(v)
		t.NoError(err)
		es, err := json.Marshal(data[i])
		t.NoError(err)
		t.EqualValues(es, as)
	}
}

func (t *_testMapper) Test_Count() {
	t.addData(20)
	ctx := context.TODO()
	res, err := t.mapper.Count(ctx, EmptyWrapperFunc)
	t.NoError(err)
	t.EqualValues(20, res)
}

func (t *_testMapper) Test_Paginate_NoData() {
	ctx := context.TODO()
	res, err := t.mapper.Paginate(ctx, Paginator{
		StartId: 0,
		Limit:   10,
	}, EmptyWrapperFunc)
	t.NoError(err)
	t.EqualValues(0, res.StartId)
	t.EqualValues(10, res.Limit)
	t.EqualValues(false, res.More)
	t.Len(res.Data, 0)
}

func (t *_testMapper) Test_Paginate_NoMoreData() {
	data := t.addData(10)
	ctx := context.TODO()
	res, err := t.mapper.Paginate(ctx, Paginator{
		StartId: 0,
		Limit:   10,
	}, EmptyWrapperFunc)
	t.NoError(err)
	t.EqualValues(0, res.StartId)
	t.EqualValues(10, res.Limit)
	t.EqualValues(false, res.More)
	t.Len(res.Data, 10)
	for i, v := range res.Data {
		as, err := json.Marshal(v)
		t.NoError(err)
		es, err := json.Marshal(data[i])
		t.NoError(err)
		t.EqualValues(es, as)
	}
}

func (t *_testMapper) Test_Paginate_FirstPageWhenMoreData() {
	data := t.addData(11)
	ctx := context.TODO()
	res, err := t.mapper.Paginate(ctx, Paginator{
		StartId: 0,
		Limit:   10,
	}, EmptyWrapperFunc)
	t.NoError(err)
	t.EqualValues(0, res.StartId)
	t.EqualValues(10, res.Limit)
	t.EqualValues(true, res.More)
	t.Len(res.Data, 10)
	for i, v := range res.Data {
		as, err := json.Marshal(v)
		t.NoError(err)
		es, err := json.Marshal(data[i])
		t.NoError(err)
		t.EqualValues(es, as)
	}
}

func (t *_testMapper) Test_Paginate_SecondPageWhenMoreData() {
	data := t.addData(19)
	ctx := context.TODO()
	res, err := t.mapper.Paginate(ctx, Paginator{
		StartId: 10,
		Limit:   10,
	}, EmptyWrapperFunc)
	t.NoError(err)
	t.EqualValues(10, res.StartId)
	t.EqualValues(10, res.Limit)
	t.EqualValues(false, res.More)
	t.Len(res.Data, 9)
	for i, v := range res.Data {
		as, err := json.Marshal(v)
		t.NoError(err)
		es, err := json.Marshal(data[i+10])
		t.NoError(err)
		t.EqualValues(es, as)
	}
}

func (t *_testMapper) addData(num int) []_testFoo {
	res := make([]_testFoo, 0, num)
	for i := 0; i < num; i++ {
		foo := _testFoo{}
		t.Require().NoError(t.db.Create(&foo).Error)
		res = append(res, foo)
	}
	return res
}

func TestBaseMapper(t *testing.T) {
	suite.Run(t, &_testMapper{})
}
