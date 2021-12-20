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
	Short: "Find jobs with config that matches the xpath query",
	Long:  `Provide xpath queries to find configs that match.`,
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
				match, err := xpathMatches(config, xpathLookup)
				if err != nil {
					log.Printf("Error matching %s: %s\n", j.Name, err)
				}
				if !match {
					continue
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

func xpathMatches(config string, xpathLookup string) (bool, error) {
	p := parser.New(parser.XMLParseNoWarning)
	doc, err := p.ParseReader(strings.NewReader(config))
	if err != nil {
		return false, err
	}
	defer doc.Free()
	root, err := doc.DocumentElement()
	if err != nil {
		return false, err
	}

	ctx, err := xpath.NewContext(root)
	if err != nil {
		return false, err
	}
	defer ctx.Free()
	x, err := ctx.Find(xpathLookup)
	if err != nil {
		return false, err
	}

	defer x.Free()
	nl := x.NodeList()
	if len(nl) == 0 {
		return false, nil
	}
	return true, nil
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
}
