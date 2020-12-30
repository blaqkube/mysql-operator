package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

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
		workdir, err := cmd.Flags().GetString("workdir")
		if err == nil {
			log.Printf("Moving WORKDIR to %s", workdir)
			os.Chdir(workdir)
		}

		// TODO: recreate the exporter user

		log.Fatal(
			http.ListenAndServe(
				fmt.Sprintf(":%d", port),
				openapi.NewRouter(service.NewMysqlAPIController(resources.DB, resources.Backup, resources.Storages)),
			),
		)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().IntP("port", "p", 8080, "agent api port")
	serveCmd.Flags().StringP("workdir", "w", "", "working directory")
}
