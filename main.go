package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/c-bata/go-prompt"
	"github.com/mattn/go-shellwords"
	"github.com/smousa/devoted-db/store"
	"github.com/spf13/cobra"
)

const helpTemplate = `{{if gt (len .Short) 0}}{{.Short}}{{end}}{{if .Runnable}}

Usage:
  {{.Use}}{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}
`

func main() {
	var p = shellwords.NewParser()
	var s store.Store = store.NewDatabase()

	// define commands
	commands := []*cobra.Command{
		{
			Use:   "set [name] [value]",
			Short: "Sets the name in the database to the given value.",
			Args:  cobra.ExactArgs(2),
			Run: func(cmd *cobra.Command, args []string) {
				s.Set(args[0], args[1])
			},
			DisableFlagParsing: true,
		}, {
			Use:   "get [name]",
			Short: "Prints the value from the given name.  If the value is not in the database, prints NULL.",
			Args:  cobra.ExactArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				value, ok := s.Get(args[0])
				if ok {
					fmt.Println(value)
				} else {
					fmt.Println("NULL")
				}
			},
			DisableFlagParsing: true,
		}, {
			Use:   "delete [name]",
			Short: "Deletes the name from the database.",
			Args:  cobra.ExactArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				s.Delete(args[0])
			},
		}, {
			Use:   "count [value]",
			Short: "Returns the number of names that have the given value assigned to them.",
			Args:  cobra.ExactArgs(1),
			Run: func(cmd *cobra.Command, args []string) {
				fmt.Println(s.Count(args[0]))
			},
			DisableFlagParsing: true,
		}, {
			Use:   "end",
			Short: "Exits the database.",
			Args:  cobra.NoArgs,
			Run: func(cmd *cobra.Command, args []string) {
				os.Exit(0)
			},
			DisableFlagParsing: true,
		}, {
			Use:   "begin",
			Short: "Begins a new transaction.",
			Args:  cobra.NoArgs,
			Run: func(cmd *cobra.Command, args []string) {
				s = s.Begin()
			},
			DisableFlagParsing: true,
		}, {
			Use:   "rollback",
			Short: "Rolls back the most recent transaction",
			Args:  cobra.NoArgs,
			Run: func(cmd *cobra.Command, args []string) {
				post, err := s.Rollback()
				if err != nil {
					fmt.Println(err)
					return
				}
				s = post
			},
			DisableFlagParsing: true,
		}, {
			Use:   "commit",
			Short: "Commits all of the open transactions",
			Args:  cobra.NoArgs,
			Run: func(cmd *cobra.Command, args []string) {
				s = s.Commit()
			},
			DisableFlagParsing: true,
		},
	}

	var rootCmd = &cobra.Command{
		Use:                "devoted-db [command]",
		SilenceUsage:       true,
		SilenceErrors:      true,
		DisableFlagParsing: true,
	}
	rootCmd.AddCommand(commands...)
	rootCmd.SetHelpTemplate(helpTemplate)

	fmt.Println("Welcome to devoted-db :)")
	fmt.Println("Type \"help\" for more information.")
	fmt.Println("Type \"help [command]\" for information about a specific command.")
	fmt.Println("Type \"end\" to exit.")

	// main loop
	for {
		t := prompt.Input(">> ", func(d prompt.Document) []prompt.Suggest { return nil })

		args, err := p.Parse(t)
		if err != nil {
			fmt.Println(err)
		}
		if len(args) == 0 {
			continue
		}
		args[0] = strings.ToLower(args[0])

		rootCmd.SetArgs(args)
		if err := rootCmd.Execute(); err != nil {
			fmt.Println(err)
		}
		rootCmd.ResetFlags()
	}
}
