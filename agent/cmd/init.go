package cmd

import (
	"fmt"
	"log"
	"os"
	"strings"

	openapi "github.com/blaqkube/mysql-operator/agent/go"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialization steps and stops",
	Long: `initialization consists in
   - recovery a backup when specified
   - create a user for the api commands`,
	Run: func(cmd *cobra.Command, args []string) {
		restore, _ := cmd.Flags().GetBool("restore")
		if !restore {
			fmt.Printf("nothing to restore\n")
			return
		}

		log.Printf("Restore database...")
		storage, err := cmd.Flags().GetString("type")
		if err != nil || storage == "" {
			storage = viper.GetString("type")
			if storage == "" {
				storage = "s3"
			}
		}
		location, err := cmd.Flags().GetString("location")
		if err != nil || location == "" {
			location = viper.GetString("location")
		}
		bucket, err := cmd.Flags().GetString("bucket")
		if err != nil || bucket == "" {
			bucket = viper.GetString("bucket")
		}
		if bucket == "" {
			fmt.Printf("Missing BUCKET, value: %s", bucket)
			os.Exit(1)
		}
		if location == "" {
			fmt.Printf("Missing LOCATION, value: %s", location)
			os.Exit(1)
		}
		fpath := strings.Split(location, string(os.PathSeparator))
		localfile := fpath[len(fpath)-1]
		_, err = os.Stat(localfile)
		if err == nil {
			log.Printf("file %s already loaded", localfile)
			os.Exit(0)
		}
		if !os.IsNotExist(err) {
			log.Printf("file %s stat error: %v", localfile, err)
			os.Exit(1)
		}
		payload := &openapi.BackupRequest{
			Backend:  storage,
			Bucket:   bucket,
			Location: location,
		}
		err = resources.Storages[payload.Backend].Pull(payload, localfile)
		if err != nil {
			log.Printf("error pulling %s: %v", localfile, err)
			os.Exit(1)
		}
		return
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolP("restore", "r", false, "restore a dump file")
	initCmd.Flags().StringP("location", "l", "", "file name on bucket")
	initCmd.Flags().StringP("bucket", "b", "", "dump file bucket")
	initCmd.Flags().StringP("type", "t", "", "type of backend (s3, blackhole)")
}
