package cmd

import (
	"fmt"
	"log"
	"net/http"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
	"github.com/blaqkube/mysql-operator/agent/service"
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
			port = 8080
		}
		// TODO: recreate the exporter user

		log.Fatal(
			http.ListenAndServe(
				fmt.Sprintf(":%d", port),
				openapi.NewRouter(service.NewMysqlAPIController(resources.DB, resources.Backup, resources.Storage)),
			),
		)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().IntP("port", "p", 8080, "agent api port")
}
