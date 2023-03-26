package cmd

import (
	"os"

	"github.com/Cyy92/HeterogeneousVirtualization/faas-cli/cmd/function"
	"github.com/Cyy92/HeterogeneousVirtualization/faas-cli/cmd/runtime"
	"github.com/Cyy92/HeterogeneousVirtualization/faas-cli/cmd/virt"
	"github.com/spf13/cobra"
)

var faasCmd = &cobra.Command{
	Use:   "faas-cli",
	Short: "Manage FaaS",
	Long: `
	Manage FaaS functions and virtual machines from the command line interface
	`,
	Run: func(cmd *cobra.Command, args []string) {
		cmd.Usage()
	},
}

func init() {

	faasCmd.SetUsageTemplate(usageTemplate)
	faasCmd.SetHelpTemplate(helpTemplate)

	faasCmd.AddCommand(versionCmd)
	faasCmd.AddCommand(function.FunctionCmd)
	faasCmd.AddCommand(virt.VirtCmd)
	faasCmd.AddCommand(runtime.RuntimeCmd)
}

func Execute() {
	if err := faasCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var usageTemplate = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if .IsAvailableCommand }}
  {{rpad .NameAndAliases 20}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}

`

var helpTemplate = `{{with (or .Long .Short)}}{{. | trimTrailingWhitespaces}}
{{end}}{{if or .Runnable .HasSubCommands}}{{.UsageString}}{{end}}`
