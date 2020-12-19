package cmd

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/blaqkube/mysql-operator/agent/backend"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
)

var cfgFile string

var rootCmd = &cobra.Command{
	Use:   "mysql-agent",
	Short: "an agent to work with blaqkube/mysql-operator",
	Long: `an agent to perform number of tasks for blaqkube/mysql-operator:
   - recovering a backup before we start the instance
   - create a secured user to secure the root user
   - connect to the database to create backups, users and database`,
}

// Backend is a type used to store backend resources
type Backend struct {
	Backup   backend.Backup
	DB       *sql.DB
	Instance backend.Instance
	Storage  backend.Storage
}

var resources *Backend

// Execute start the agent with the various attributes
func Execute(
	backup backend.Backup,
	db *sql.DB,
	instance backend.Instance,
	storage backend.Storage,
) {
	resources = &Backend{
		Backup:   backup,
		DB:       db,
		Instance: instance,
		Storage:  storage,
	}

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".agent" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".agent")
	}

	viper.SetEnvPrefix("agt")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}
}
