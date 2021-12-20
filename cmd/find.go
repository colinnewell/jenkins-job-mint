package cmd

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/bndr/gojenkins"
	"github.com/lestrrat-go/libxml2/parser"
	"github.com/lestrrat-go/libxml2/xpath"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// findCmd represents the find command
var findCmd = &cobra.Command{
	Args:  hasCredentials,
	Use:   "find",
	Short: "Find jobs",
	Long:  `Look up by type, color or match a config using an xpath query.`,
	Run: func(cmd *cobra.Command, args []string) {
		user := viper.GetString("user")
		url := viper.GetString("url")
		token := viper.GetString("token")
		jenkins := gojenkins.CreateJenkins(nil, url, user, token)
		ctx := context.Background()
		if _, err := jenkins.Init(ctx); err != nil {
			log.Fatal("Failed to connect to Jenkins", err)
		}

		// loop through the jobs.
		jobs, err := jenkins.GetAllJobNames(ctx)
		if err != nil {
			log.Fatal("Error getting jobs", err)
		}

		jobType := viper.GetString("type")
		color := viper.GetString("color")
		xpathLookup := viper.GetString("xpath")
		newValue := viper.GetString("new-value")
		setValue := viper.IsSet("new-value")
		for _, j := range jobs {
			if jobType != "" && j.Class != jobType {
				continue
			}
			if color != "" && j.Color != color {
				continue
			}
			if xpathLookup != "" {
				// download config and do stuff
				job, err := jenkins.GetJob(ctx, j.Name)
				if err != nil {
					log.Printf("Failed to get job for %s: %s\n", j.Name, err)
					continue
				}
				config, err := job.GetConfig(ctx)
				if err != nil {
					log.Printf("Failed to get config for %s: %s\n", j.Name, err)
				}
				match, updated, err := xpathMatches(config, xpathLookup, setValue, newValue)
				if err != nil {
					log.Printf("Error matching %s: %s\n", j.Name, err)
				}
				if !match {
					continue
				}
				if setValue {
					job.UpdateConfig(ctx, updated)
				}
			}
			if viper.GetBool("verbose") {
				fmt.Printf("%s (type: %s) - %s - %s\n", j.Name, j.Class, j.Color, j.Url)
			} else {
				fmt.Println(j.Name)
			}
		}
	},
}

func xpathMatches(config string, xpathLookup string, setValue bool, newValue string) (bool, string, error) {
	p := parser.New(parser.XMLParseNoWarning)
	doc, err := p.ParseReader(strings.NewReader(config))
	if err != nil {
		return false, "", err
	}
	defer doc.Free()
	root, err := doc.DocumentElement()
	if err != nil {
		return false, "", err
	}

	ctx, err := xpath.NewContext(root)
	if err != nil {
		return false, "", err
	}
	defer ctx.Free()
	x, err := ctx.Find(xpathLookup)
	if err != nil {
		return false, "", err
	}

	defer x.Free()
	nl := x.NodeList()
	if len(nl) == 0 {
		return false, "", nil
	}
	if !setValue {
		return true, "", nil
	}

	for i := range nl {
		nl[i].SetNodeValue(newValue)
	}
	return true, doc.String(), nil
}

func init() {
	rootCmd.AddCommand(findCmd)

	findCmd.PersistentFlags().String("type", "", "Limit to a job type (e.g. hudson.model.FreeStyleProject)")
	if err := viper.BindPFlag("type", findCmd.PersistentFlags().Lookup("type")); err != nil {
		log.Fatal("Programmer error:", err)
	}

	findCmd.PersistentFlags().String("color", "", "Limit to a job color (e.g. blue)")
	if err := viper.BindPFlag("color", findCmd.PersistentFlags().Lookup("color")); err != nil {
		log.Fatal("Programmer error:", err)
	}

	findCmd.PersistentFlags().String("xpath", "", "Limit to jobs with config matching the xpath query")
	if err := viper.BindPFlag("xpath", findCmd.PersistentFlags().Lookup("xpath")); err != nil {
		log.Fatal("Programmer error:", err)
	}

	findCmd.PersistentFlags().String("new-value", "", "Change string value in config matched by xpath")
	if err := viper.BindPFlag("new-value", findCmd.PersistentFlags().Lookup("new-value")); err != nil {
		log.Fatal("Programmer error:", err)
	}
}
