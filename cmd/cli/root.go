package cli

import (
	"fmt"
	"os"

	"github.com/open-ug/runner/cmd/builder"
	containerstart "github.com/open-ug/runner/cmd/container-start"
	containerstop "github.com/open-ug/runner/cmd/container-stop"
	gitcloner "github.com/open-ug/runner/cmd/git-cloner"
	"github.com/spf13/cobra"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "flutter-ci",
	Short: "Flutter CI",
	Long: `
Flutter CI is a continuous integration tool for Flutter applications.
	`,
	Version: "0.1.0",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Flutter CI ")
		fmt.Println("Run 'flutter-ci --help' to see available commands")
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.AddCommand(GitClonerDriverCmd)
	rootCmd.AddCommand(ContainerStartDriverCmd)
	rootCmd.AddCommand(ContainerStopDriverCmd)
	rootCmd.AddCommand(BuilderDriverCmd)

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

var GitClonerDriverCmd = &cobra.Command{
	Use:   "git-cloner",
	Short: "Start the Git Cloner Driver",
	Long:  `Start the Git Cloner Driver`,
	Run: func(cmd *cobra.Command, args []string) {
		gitcloner.Listen()
	},
}

var ContainerStartDriverCmd = &cobra.Command{
	Use:   "container-start",
	Short: "Start the Container Start Driver",
	Long:  `Start the Container Start Driver`,
	Run: func(cmd *cobra.Command, args []string) {
		containerstart.Listen()
	},
}

var ContainerStopDriverCmd = &cobra.Command{
	Use:   "container-stop",
	Short: "Start the Container Stop Driver",
	Long:  `Start the Container Stop Driver`,
	Run: func(cmd *cobra.Command, args []string) {
		containerstop.Listen()
	},
}

var BuilderDriverCmd = &cobra.Command{
	Use:   "builder",
	Short: "Start the Builder Driver",
	Long:  `Start the Builder Driver`,
	Run: func(cmd *cobra.Command, args []string) {
		builder.Listen()
	},
}
