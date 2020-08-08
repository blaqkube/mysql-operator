package cmd

import (
	"fmt"
	"log"
	"os"

	openapi "github.com/blaqkube/mysql-operator/agent/go"
	"github.com/blaqkube/mysql-operator/agent/mysql"
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
		filename, err := cmd.Flags().GetString("filename")
		if err != nil || filename == "" {
			filename = viper.GetString("filename")
		}
		bucket, err := cmd.Flags().GetString("bucket")
		if err != nil || bucket == "" {
			bucket = viper.GetString("bucket")
		}
		filePath, err := cmd.Flags().GetString("path")
		if err != nil || filePath == "" {
			filePath = viper.GetString("path")
		}
		if filePath == "" || bucket == "" || filename == "" {
			fmt.Println("Missing parameter, check FILENAME, BUCKET and FILEPATH are set")
			os.Exit(1)
		}
		my := mysql.NewS3MysqlBackup()
		b := &openapi.Backup{
			Location: filePath,
			S3access: openapi.S3Info{
				Bucket: bucket,
				Path:   filePath,
			},
		}
		err = my.PullS3File(b, filePath, filename)
		if err != nil {
			fmt.Printf("Error while reading s3://%s%s: %v\n", bucket, filePath, err)
			os.Exit(1)
		}
		return
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolP("restore", "r", false, "restore a dump file")
	initCmd.Flags().StringP("filename", "f", "", "dump file name")
	initCmd.Flags().StringP("bucket", "b", "", "dump file bucket")
	initCmd.Flags().StringP("path", "p", "", "dump file remote path")
}
