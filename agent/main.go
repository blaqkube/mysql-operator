/*
 * blaqkube MySQL agent
 *
 * Agent used by [blaqkube MySQL operator](http://github.com/blaqkube/mysql-operator) to manage MySQL backup/restore 
 *
 * API version: 0.0.1
 * Contact: contact@blaqkube.io
 * Generated by: OpenAPI Generator (https://openapi-generator.tech)
 */

package main

import (
	"log"
	"net/http"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
)

func main() {
	log.Printf("Server started")

	MysqlApiService := openapi.NewMysqlApiService()
	MysqlApiController := openapi.NewMysqlApiController(MysqlApiService)

	PetApiService := openapi.NewPetApiService()
	PetApiController := openapi.NewPetApiController(PetApiService)

	router := openapi.NewRouter(MysqlApiController, PetApiController)

	log.Fatal(http.ListenAndServe(":8080", router))
}
