package grant

import (
	"database/sql"
	"fmt"
	"regexp"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
)

// MysqlGrantService is a service that implements the logic for the MysqlGrantServicer
// This service should implement the business logic for every endpoint for the MysqlGrant API.
// Include any external packages or services that will be required by this service.
type MysqlGrantService struct {
	DB *sql.DB
}

// NewMysqlGrantService creates a default api service
func NewMysqlGrantService(db *sql.DB) MysqlGrantServicer {
	return &MysqlGrantService{
		DB: db,
	}
}

// CreateGrantByUserDatabase - create an on-demand user
func (s *MysqlGrantService) CreateGrantByUserDatabase(grant openapi.Grant, user, database, apiKey string) (interface{}, error) {
	fmt.Printf("Connect to database\n")
	sql := fmt.Sprintf("GRANT ALL PRIVILEGES ON %s.* TO '%s'@'%%'", database, user)
	if grant.AccessMode == ReadOnlyAccessMode {
		sql = fmt.Sprintf("GRANT SELECT ON %s.* TO '%s'@'%%'", database, user)
	}
	_, err := s.DB.Exec(sql)
	if err != nil {
		fmt.Printf("Error granting privileges; %v\n", err)
		return nil, err
	}
	return &grant, nil
}

// GetGrantByUserDatabase - Get user properties
func (s *MysqlGrantService) GetGrantByUserDatabase(user, database, apiKey string) (interface{}, error) {
	rows, err := s.DB.Query(fmt.Sprintf("SHOW GRANTS FOR '%s'@'%%'", user))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	mode := "none"
	rwmode, _ := regexp.Compile(regexp.QuoteMeta(
		fmt.Sprintf("GRANT SELECT ON `%s`.* TO `%s`@`%%`", database, user),
	))
	romode, _ := regexp.Compile(regexp.QuoteMeta(
		fmt.Sprintf("GRANT ALL PRIVILEGES ON `%s`.* TO `%s`@`%%`", database, user),
	))
	for rows.Next() {
		grant := ""
		err = rows.Scan(&grant)
		if err != nil {
			return nil, err
		}
		if romode.MatchString(grant) && mode != ReadWriteAccessMode {
			mode = ReadOnlyAccessMode
		}
		if rwmode.MatchString(grant) {
			mode = ReadWriteAccessMode
		}
	}
	return &openapi.Grant{AccessMode: mode}, nil
}
