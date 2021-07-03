package cmd

import (
	"context"
	"errors"
	"log"
	"os"

	"github.com/bndr/gojenkins"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// jobCmd represents the job command
var jobCmd = &cobra.Command{
	Args: func(cmd *cobra.Command, args []string) error {
		if err := standardValidation(cmd, args); err != nil {
			return err
		}
		if err := cobra.ExactValidArgs(1)(cmd, args); err != nil {
			return err
		}
		// check we have template
		if !viper.IsSet("template") {
			return errors.New("template must be specified")
		}
		return nil
	},
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

		template := viper.GetString("template")
		config, err := os.ReadFile(template)
		if err != nil {
			log.Fatalf("Failed to load template %s: %s", template, err)
		}

		if _, err := jenkins.CreateJob(ctx, string(config), args[0]); err != nil {
			log.Fatalf("Failed to create job %s: %s", args[0], err)
		}

	},
}

func init() {
	rootCmd.AddCommand(jobCmd)

	jobCmd.PersistentFlags().String("template", "", "Template file containing jenkins config")
	if err := viper.BindPFlag("template", jobCmd.PersistentFlags().Lookup("template")); err != nil {
		log.Fatal("Programmer error:", err)
	}
}
