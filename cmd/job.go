package cmd

import (
	"context"
	"log"

	"github.com/bndr/gojenkins"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// jobCmd represents the job command
var jobCmd = &cobra.Command{
	Args:  cobra.ExactValidArgs(1),
	Use:   "job new-job-name",
	Short: "Create a job",
	Long: `Create a new job in Jenkins

Provide info to replace in the configs.
Jenkins configs support templating.
`,
	Run: func(cmd *cobra.Command, args []string) {
		// FIXME: ensure we also allow creation in folders
		user := viper.GetString("user")
		url := viper.GetString("url")
		token := viper.GetString("token")
		jenkins := gojenkins.CreateJenkins(nil, url, user, token)
		ctx := context.Background()
		if _, err := jenkins.Init(ctx); err != nil {
			log.Fatal("Failed to connect to Jenkins", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(jobCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// jobCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// jobCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
