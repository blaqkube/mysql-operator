package backup

/*
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	openapi "github.com/blaqkube/mysql-operator/agent/go"
	_ "github.com/go-sql-driver/mysql"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type Suite struct {
	suite.Suite
	DB          *sql.DB
	mock        sqlmock.Sqlmock
	testService MysqlBackupServicer
}

func (s *Suite) SetupSuite() {
	var err error
	s.DB, s.mock, err = sqlmock.New()
	require.NoError(s.T(), err)

	s.testService = &MysqlBackupService{DB: s.DB}
}

func (s *Suite) Test_Delete() {
	var (
		id     = int32(1)
		tenant = "test1"
	)

	s.mock.ExpectBegin()

	s.mock.ExpectExec(regexp.QuoteMeta(
		"UPDATE `NGN_FACTS` SET `deleted_at`=? WHERE tenant = ? and id=?",
	)).
		WithArgs(sqlmock.AnyArg(), tenant, id).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	err := s.testService.Delete(id, "test1")
	require.NoError(s.T(), err)
}

func (s *Suite) Test_List() {
	var (
		id     = int32(1)
		tenant = "test1"
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `NGN_FACTS` WHERE tenant = ? AND `NGN_FACTS`.`deleted_at` IS NULL",
	)).
		WithArgs(tenant).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "tenant", "inventory_id", "product_id", "repository_id", "indicator_id", "value"}).
			AddRow(id, time.Now(), time.Now(), nil, tenant, uint(1), uint(1), uint(1), uint(1), float64(5)))

	v, err := s.testService.List("test1")
	require.NoError(s.T(), err)
	switch i := v.(type) {
	case openapi.ListFactsResponse:
		require.Equal(s.T(), int32(1), i.Size, "Return an array size 1")
	default:
		panic("wrong type")
	}

}

func (s *Suite) Test_GetByID() {
	var (
		id     = int32(1)
		tenant = "test1"
	)

	s.mock.ExpectQuery(regexp.QuoteMeta(
		"SELECT * FROM `NGN_FACTS` WHERE tenant = ? and id=? AND `NGN_FACTS`.`deleted_at` IS NULL ORDER BY `NGN_FACTS`.`id` LIMIT 1",
	)).
		WithArgs(tenant, uint(id)).
		WillReturnRows(sqlmock.NewRows([]string{"id", "created_at", "updated_at", "deleted_at", "tenant", "inventory_id", "product_id", "repository_id", "indicator_id", "value"}).
			AddRow(id, time.Now(), time.Now(), nil, tenant, uint(1), uint(1), uint(1), uint(1), float64(5)))

	v, err := s.testService.GetByID(id, "test1")
	require.NoError(s.T(), err)
	switch i := v.(type) {
	case openapi.Fact:
		require.Equal(s.T(), int32(1), i.Id, "Return Fact with Id 1")
	default:
		panic("wrong type")
	}

}

func (s *Suite) Test_UpdateByID() {
	var (
		id     = int32(1)
		tenant = "test1"
	)

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(
		"INSERT INTO `NGN_FACTS` (`created_at`,`updated_at`,`deleted_at`,`tenant`,`inventory_id`,`product_id`,`repository_id`,`indicator_id`,`value`) VALUES (?,?,?,?,?,?,?,?,?)",
	)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), tenant, uint(1), uint(1), uint(1), uint(1), float64(5)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	input := openapi.Fact{
		Value: float64(5),
		Product: openapi.Product{
			Id: int32(1),
			Inventory: openapi.Inventory{
				Id: int32(1),
			},
		},
		Repository: openapi.Repository{
			Id: int32(1),
		},
		Indicator: openapi.Indicator{
			Id: int32(1),
		},
	}
	v, err := s.testService.UpdateByID(input, id, "test1")
	t, err := json.Marshal(v)
	fmt.Println(string(t))
	require.NoError(s.T(), err)
	switch i := v.(type) {
	case openapi.Fact:
		require.Equal(s.T(), int32(1), i.Id, "Return Fact Id 1")
	default:
		panic("wrong type")
	}

}

func (s *Suite) Test_Create() {
	var (
		id     = int32(1)
		tenant = "test1"
	)

	s.mock.ExpectBegin()
	s.mock.ExpectExec(regexp.QuoteMeta(
		"INSERT INTO `NGN_FACTS` (`created_at`,`updated_at`,`deleted_at`,`tenant`,`inventory_id`,`product_id`,`repository_id`,`indicator_id`,`value`) VALUES (?,?,?,?,?,?,?,?,?)",
	)).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), tenant, uint(1), uint(1), uint(1), uint(1), float64(5)).
		WillReturnResult(sqlmock.NewResult(1, 1))
	s.mock.ExpectCommit()
	input := openapi.Fact{
		Value: float64(5),
		Product: openapi.Product{
			Id: int32(1),
			Inventory: openapi.Inventory{
				Id: int32(1),
			},
		},
		Repository: openapi.Repository{
			Id: int32(1),
		},
		Indicator: openapi.Indicator{
			Id: int32(1),
		},
	}

	v, err := s.testService.Create(input, "test1")
	t, err := json.Marshal(v)
	fmt.Println(string(t))
	require.NoError(s.T(), err)
	switch i := v.(type) {
	case openapi.Fact:
		require.Equal(s.T(), id, i.Id, "Return an array size 1")
	default:
		panic("wrong type")
	}

}

func TestSuite(t *testing.T) {
	suite.Run(t, &Suite{})
}
*/
