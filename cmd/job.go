package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"log"
	"os"
	"path"
	"strings"
	"text/template"

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

The templates use text/template syntax and have job plus other variables
provided via the command line passed in.
`,
	Run: func(cmd *cobra.Command, args []string) {
		job := args[0]

		variablesString := viper.GetString("variables")
		vars := map[string]interface{}{}
		vars["job"] = job
		if variablesString != "" {
			if err := json.Unmarshal([]byte(variablesString), &vars); err != nil {
				log.Fatalf("Unable to decode the variables: %s", err)
			}
		}

		user := viper.GetString("user")
		url := viper.GetString("url")
		token := viper.GetString("token")

		jenkins := gojenkins.CreateJenkins(nil, url, user, token)
		ctx := context.Background()
		if viper.GetBool("debug") {
			//nolint:staticcheck
			ctx = context.WithValue(ctx, "debug", true)
		}
		if _, err := jenkins.Init(ctx); err != nil {
			log.Fatal("Failed to connect to Jenkins", err)
		}

		templateName := viper.GetString("template")
		config, err := os.ReadFile(templateName)
		if err != nil {
			log.Fatalf("Failed to load template %s: %s", templateName, err)
		}

		tmpl, err := template.New("job-config").Parse(string(config))
		if err != nil {
			log.Fatalf("Failed to parse template %s: %s", templateName, err)
		}
		var buf bytes.Buffer
		if err = tmpl.Execute(&buf, vars); err != nil {
			log.Fatalf("Failed to process template %s: %s", templateName, err)
		}

		folder := viper.GetString("folder")
		if folder != "" {
			folders := splitFolders(folder)
			if _, err := jenkins.CreateJobInFolder(ctx, buf.String(), job, folders...); err != nil {
				log.Fatalf("Failed to create job %s: %s", job, err)
			}
		} else {
			if _, err := jenkins.CreateJob(ctx, buf.String(), job); err != nil {
				if strings.HasPrefix(err.Error(), "A job already exists with the name") {
					// this could be neater, no error handling.
					jenkins.UpdateJob(ctx, job, buf.String())
				} else {
					log.Fatalf("Failed to create job %s: %s", job, err)
				}
			}
		}

	},
}

func init() {
	rootCmd.AddCommand(jobCmd)

	flags := []struct {
		Name        string
		Description string
	}{
		{"template", "Template file containing jenkins config (required)"},
		{"variables", "Variables to use in the template (json)"},
		{"folder", "Jenkins folder to place job in (optional).  Example: a/b"},
	}
	for _, f := range flags {
		jobCmd.PersistentFlags().String(f.Name, "", f.Description)
		if err := viper.BindPFlag(f.Name, jobCmd.PersistentFlags().Lookup(f.Name)); err != nil {
			log.Fatal("Programmer error:", err)
		}
	}

	jobCmd.PersistentFlags().Bool("debug", false, "Emit debug from the jenkins library")
	if err := viper.BindPFlag("debug", jobCmd.PersistentFlags().Lookup("debug")); err != nil {
		log.Fatal("Programmer error:", err)
	}
}

func splitFolders(folder string) []string {
	folders := []string{}
	dir, file := path.Split(folder)
	for ; dir != ""; dir, file = path.Split(folder) {
		folders = append(folders, file)
		folder = path.Clean(dir)
	}
	folders = append(folders, file)
	for left, right := 0, len(folders)-1; left < right; left, right = left+1, right-1 {
		folders[left], folders[right] = folders[right], folders[left]
	}

	return folders
}
