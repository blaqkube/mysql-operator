package cmd

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/blaqkube/mysql-operator/agent/backend/mysql"
	openapi "github.com/blaqkube/mysql-operator/agent/go"
	"github.com/blaqkube/mysql-operator/agent/service"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const numberOfDBChecks = 24

func createExporterUser(expUsername, expPassword string) error {
	if expUsername == "" || expPassword == "" {
		log.Printf("Skipping exporter creation")
	}

	log.Printf("Creating exporter user %s, starting", expUsername)
	i := mysql.NewInstance(resources.DB)
	err := i.Check(numberOfDBChecks)
	if err != nil {
		log.Printf("Could not connect to the database after %d retries", numberOfDBChecks)
		return err
	}

	_, err = resources.DB.Exec(
		fmt.Sprintf(
			`CREATE USER '%s'@'%%' IDENTIFIED BY '%s' WITH MAX_USER_CONNECTIONS 3`,
			expUsername,
			expPassword,
		),
	)
	if err != nil {
		log.Printf("Error creating exporter user: %v", err)
		return err
	}
	_, err = resources.DB.Exec(
		fmt.Sprintf(
			`GRANT PROCESS, REPLICATION CLIENT, SELECT ON *.* TO '%s'@'%%'`,
			expUsername,
		),
	)
	if err != nil {
		log.Printf("Error granting privileges to exporter: %v", err)
		return err
	}
	log.Printf("Creating exporter user %s, done", expUsername)
	return nil
}

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
			if workdir == "" {
				workdir = viper.GetString("workdir")
			}
			if workdir != "" {
				log.Printf("Moving WORKDIR to %s", workdir)
				os.Chdir(workdir)
			}
		}

		expUsername := viper.GetString("exporter_username")
		expPassword := viper.GetString("exporter_password")
		createExporterUser(expUsername, expPassword)

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
