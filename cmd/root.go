/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"go_tools/log"
	"os"

	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "go_tools",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		log.Info("this is first cobra example")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.go_tools.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	copyCmd.PersistentFlags().String("s", "", "source file/directory path")
	copyCmd.PersistentFlags().String("t", "", "target file/directory path")
	copyCmd.PersistentFlags().String("p", "", "copy parallelism")
	rootCmd.AddCommand(copyCmd)

	compressCmd.PersistentFlags().String("s", "", "source file/directory path")
	compressCmd.PersistentFlags().String("t", "", "target file/directory path")
	compressCmd.PersistentFlags().String("p", "", "compress parallelism")
	rootCmd.AddCommand(compressCmd)

	decompressCmd.PersistentFlags().String("s", "", "source file/directory path")
	decompressCmd.PersistentFlags().String("t", "", "target file/directory path")
	decompressCmd.PersistentFlags().String("p", "", "decompress parallelism")
	rootCmd.AddCommand(decompressCmd)
}
