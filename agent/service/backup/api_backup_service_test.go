package backup

import (
	"testing"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
	_ "github.com/go-sql-driver/mysql"

	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type BackupSuite struct {
	suite.Suite
	testService MysqlBackupServicer
}

func (s *BackupSuite) SetupSuite() {
	s.testService = NewMysqlBackupService(&mockBackupPrimitive{})
}

func (s *BackupSuite) Test_Delete() {
	_, err := s.testService.DeleteBackup("1", "2")
	require.NoError(s.T(), err)
}

func (s *BackupSuite) Test_GetByName() {
	_, _, err := s.testService.GetBackupByName("2006-01-02T15:04:05Z", "2")
	require.NoError(s.T(), err)
}

func (s *BackupSuite) Test_GetByNameWithZeroTime() {
	_, _, err := s.testService.GetBackupByName("0001-01-01T00:00:00Z", "2")
	require.NoError(s.T(), err)
}

func (s *BackupSuite) Test_GetByNameWithError() {
	_, _, err := s.testService.GetBackupByName("123", "2")
	require.Error(s.T(), err)
}

func (s *BackupSuite) Test_CreateBackup() {
	_, err := s.testService.CreateBackup(openapi.Backup{}, "2")
	require.NoError(s.T(), err)
}

func (s *BackupSuite) Test_CreateBackupWithError() {
	_, err := s.testService.CreateBackup(openapi.Backup{Location: "/123"}, "2")
	require.Error(s.T(), err)
}

/*
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

func TestBackupSuite(t *testing.T) {
	suite.Run(t, &BackupSuite{})
}
