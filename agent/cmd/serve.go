package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
	"github.com/blaqkube/mysql-operator/agent/mysql"
	"github.com/blaqkube/mysql-operator/agent/service"
	_ "github.com/go-sql-driver/mysql"
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the MySQL agent",
	Long:  `Start the MySQL agent and serve the OpenAPI for database, user and backup`,
	Run: func(cmd *cobra.Command, args []string) {
		port, err := cmd.Flags().GetInt("port")
		db, err := sql.Open("mysql", "root@tcp(localhost:3306)/")
		tool := mysql.NewDBTools(db)
		if err != nil {
			log.Panic(err)
		}

		log.Printf("Create exporter user")
		err = tool.CheckDB(20)
		if err != nil {
			fmt.Printf("Error checking database: %v\n", err)
			os.Exit(1)
		}
		err = tool.CreateExporter()
		if err != nil {
			fmt.Printf("Error create user: %v\n", err)
			os.Exit(1)
		}

		log.Printf("Server started")
		my := mysql.NewS3MysqlBackup()
		log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), openapi.NewRouter(service.NewMysqlApiController(db, my))))
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
	serveCmd.Flags().IntP("port", "p", 8080, "agent api port")
}
