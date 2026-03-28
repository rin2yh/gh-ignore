package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/rin2yh/gh-ignore/internal/github"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "gh ignore",
	Short: "Fetch GitHub gitignore templates",
}

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available gitignore templates",
	RunE: func(cmd *cobra.Command, args []string) error {
		templates, err := github.ListTemplates()
		if err != nil {
			return err
		}
		for _, t := range templates {
			fmt.Println(t)
		}
		return nil
	},
}

var outputFile string

var getCmd = &cobra.Command{
	Use:   "get <name>",
	Short: "Fetch a gitignore template and append it to a file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		input := args[0]

		templates, err := github.ListTemplates()
		if err != nil {
			return err
		}

		var matched string
		for _, t := range templates {
			if strings.EqualFold(t, input) {
				matched = t
				break
			}
		}
		if matched == "" {
			return fmt.Errorf("template %q not found", input)
		}

		source, err := github.GetTemplate(matched)
		if err != nil {
			return err
		}

		f, err := os.OpenFile(outputFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("failed to open file: %w", err)
		}
		defer f.Close() //nolint:errcheck // error is non-actionable in defer context

		content := fmt.Sprintf("\n### %s ###\n%s", matched, source)
		if _, err := f.WriteString(content); err != nil {
			return fmt.Errorf("failed to write: %w", err)
		}

		fmt.Fprintf(os.Stderr, "appended %s to %s\n", matched, outputFile)
		return nil
	},
}

func init() {
	getCmd.Flags().StringVarP(&outputFile, "output", "o", ".gitignore", "output file")
	rootCmd.AddCommand(listCmd, getCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
