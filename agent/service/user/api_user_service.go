package user

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
)

// MysqlUserService is a service that implents the logic for the MysqlUserServicer
// This service should implement the business logic for every endpoint for the MysqlUser API.
// Include any external packages or services that will be required by this service.
type MysqlUserService struct {
	DB *sql.DB
}

// NewMysqlUserService creates a default api service
func NewMysqlUserService(db *sql.DB) MysqlUserServicer {
	return &MysqlUserService{
		DB: db,
	}
}

// CreateUser - create an on-demand user
func (s *MysqlUserService) CreateUser(user openapi.User, apiKey string) (interface{}, error) {
	// TODO - update CreateUser with the required logic for this service method.
	// Add api_mysql_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.
	fmt.Printf("Connect to database\n")
	var name string
	err := s.DB.QueryRow("SELECT user FROM mysql.user where user=?", user.Username).Scan(&name)
	if err == nil {
		return user, nil
	} else if err != sql.ErrNoRows {
		return nil, err
	}
	sql := fmt.Sprintf(
		"create user '%s'@'%%' identified by '%s'",
		user.Username,
		user.Password,
	)
	fmt.Println(sql)
	_, err = s.DB.Exec(sql)
	if err != nil {
		fmt.Printf("Error %v\n", err)
		return nil, err
	}
	return user, nil
}

// DeleteUser - Deletes a user
func (s *MysqlUserService) DeleteUser(user string, apiKey string) (interface{}, error) {
	// TODO - update DeleteUser with the required logic for this service method.
	// Add api_mysql_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.
	return nil, errors.New("service method 'DeleteUser' not implemented")
}

// GetUserByName - Get user properties
func (s *MysqlUserService) GetUserByName(user string, apiKey string) (interface{}, error) {
	// TODO - update GetUserByName with the required logic for this service method.
	// Add api_mysql_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.
	var name string
	err := s.DB.QueryRow("SELECT User FROM mysql.user where User=?", user).Scan(&name)
	if err != nil {
		fmt.Printf("Error %v\n", err)
		return nil, err
	}
	return openapi.User{Username: name}, nil
}

// GetUsers - list all users
func (s *MysqlUserService) GetUsers(apiKey string) (interface{}, error) {
	results, err := s.DB.Query("SELECT User FROM mysql.user WHERE Host='%'")
	if err != nil {
		return openapi.Message{Code: int32(http.StatusInternalServerError), Message: fmt.Sprintf("%v", err)}, err
	}
	users := []openapi.User{}
	count := int32(0)
	for results.Next() {
		var name string
		err = results.Scan(&name)
		user := openapi.User{Username: name}
		users = append(users, user)
		count++
	}
	return openapi.ListUsers{
		Size:  count,
		Items: users,
	}, nil
}
