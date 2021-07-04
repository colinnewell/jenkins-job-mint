package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/bndr/gojenkins"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Args:  standardValidation,
	Use:   "get job [job2 ...]",
	Short: "Retrieve jenkins config",
	Long:  `Download the jenkins config xml and store in files.`,
	Run: func(cmd *cobra.Command, args []string) {
		user := viper.GetString("user")
		url := viper.GetString("url")
		token := viper.GetString("token")
		jenkins := gojenkins.CreateJenkins(nil, url, user, token)
		ctx := context.Background()
		if _, err := jenkins.Init(ctx); err != nil {
			log.Fatal("Failed to connect to Jenkins", err)
		}

		folder := viper.GetString("output-folder")

		for _, jobName := range args {
			// write to files.
			job, err := jenkins.GetJob(ctx, jobName)
			if err != nil {
				log.Printf("Failed to get config for %s: %s\n", jobName, err)
				continue
			}
			config, err := job.GetConfig(ctx)
			if err != nil {
				log.Printf("Failed to get config for %s: %s\n", jobName, err)
			}
			filename := fmt.Sprintf("job-%s-config-%d.xml", jobName, time.Now().Unix())
			filename = filepath.Join(folder, filename)
			if err := os.WriteFile(filename, []byte(config), 0644); err != nil {
				fmt.Printf("Failed to write config for %s: %s\n", jobName, err)
				continue
			}
			if viper.GetBool("verbose") {
				fmt.Printf("%s: wrote %s\n", jobName, filename)
			}
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	getCmd.PersistentFlags().String("output-folder", "", "Output folder for config files to be written")
	if err := viper.BindPFlag("output-folder", getCmd.PersistentFlags().Lookup("output-folder")); err != nil {
		log.Fatal("Programmer error:", err)
	}
}
