package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/spf13/viper"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "mint",
	Short: "Mint jenkins jobs",
	Long:  `This is a tool to simplify the creation of jenkins jobs`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	// Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.jenkins-job-mint.yaml)")

	var user, token, url string
	rootCmd.PersistentFlags().StringVar(&user, "user", "", "Jenkins username")
	rootCmd.PersistentFlags().StringVar(&token, "token", "", "Jenkins auth token")
	rootCmd.PersistentFlags().StringVar(&url, "url", "http://localhost:8080", "Jenkins url")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "verbose mode")

	for _, f := range []string{"user", "token", "url", "verbose"} {
		viper.BindPFlag(f, rootCmd.PersistentFlags().Lookup(f))
	}

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	// allow us to use log without timestamps
	log.SetFlags(0)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := os.UserHomeDir()
		cobra.CheckErr(err)

		// Search config in home directory with name ".jenkins-job-mint" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigType("yaml")
		viper.SetConfigName(".jenkins-job-mint")
	}

	viper.SetEnvPrefix("mint")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		if viper.GetBool("verbose") {
			fmt.Fprintln(os.Stderr, "Using config file:", viper.ConfigFileUsed())
		}
	}
}

func standardValidation(cmd *cobra.Command, args []string) error {
	if err := cobra.MinimumNArgs(1)(cmd, args); err != nil {
		return err
	}
	for _, f := range []string{"user", "token", "url"} {
		if viper.GetString(f) == "" {
			return fmt.Errorf("%s must be provided", f)
		}
	}
	return nil
}
