package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
	"github.com/blaqkube/mysql-operator/agent/service"
	"github.com/blaqkube/mysql-operator/agent/service/backup"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the MySQL agent",
	Long:  `Start the MySQL agent and serve the OpenAPI for database, user and backup`,
	Run: func(cmd *cobra.Command, args []string) {
		port, err := cmd.Flags().GetInt("port")
		if err != nil {
			log.Panic(err)
		}

		log.Printf("Create exporter user")
		err = backup.CheckDb("root@tcp(localhost:3306)/", 20)
		if err != nil {
			fmt.Printf("Error checking database: %v\n", err)
			os.Exit(1)
		}
		err = backup.CreateExporter("root@tcp(localhost:3306)/")
		if err != nil {
			fmt.Printf("Error create user: %v\n", err)
			os.Exit(1)
		}

		log.Printf("Server started")
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), openapi.NewRouter(service.NewMysqlApiController())))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().IntP("port", "p", 8080, "agent api port")
}
