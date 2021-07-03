package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// getCmd represents the get command
var getCmd = &cobra.Command{
	Args:  standardValidation,
	Use:   "get job [job2 ...]",
	Short: "Retrieve jenkins config",
	Long:  `Download the jenkins config xml.`,
	Run: func(cmd *cobra.Command, args []string) {
		user := viper.Get("user")
		fmt.Printf("get called user:%s jobs: %v\n", user, args)
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
