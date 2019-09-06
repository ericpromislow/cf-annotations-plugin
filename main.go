package main

import (
	"code.cloudfoundry.org/cli/plugin"
	"fmt"
	"os"
	"github.com/zrob/annotations/scripts/annotations"
)

type AnnotationsPlugin struct{}

func (c *AnnotationsPlugin) Run(cliConnection plugin.CliConnection, args []string) {
	var err error
	if args[0] == "set-annotation" {
		if len(args) < 4 {
			fmt.Printf("Error: usage: %s resource-type resource-name annotations... \n", c.GetMetadata().Commands[0].UsageDetails)
			os.Exit(1)
		}
		err = annotations.SetAnnotations(cliConnection, args[1:])
	} else if args[0] == "unset-annotation" {
		if len(args) < 4 {
			fmt.Printf("Error: usage: %s resource-type resource-name annotations... \n", c.GetMetadata().Commands[0].UsageDetails)
			os.Exit(1)
		}
		err = annotations.UnsetAnnotations(cliConnection, args[1:])
	} else if args[0] == "annotations" {
		if len(args) < 3 {
			fmt.Printf("Error: usage: %s resource-type resource-name \n", c.GetMetadata().Commands[0].UsageDetails)
			os.Exit(1)
		}
		err = annotations.ViewAnnotations(cliConnection, args[1:])
	} else {
		return
	}
	if err != nil {
		fmt.Fprintf(os.Stderr,"Annotations plugin error: %s\n", err)
	}
}

// GetMetadata must be implemented as part of the plugin interface
// defined by the core CLI.
//
// GetMetadata() returns a PluginMetadata struct. The first field, Name,
// determines the name of the plugin which should generally be without spaces.
// If there are spaces in the name a user will need to properly quote the name
// during uninstall otherwise the name will be treated as seperate arguments.
// The second value is a slice of Command structs. Our slice only contains one
// Command Struct, but could contain any number of them. The first field Name
// defines the command `cf basic-plugin-command` once installed into the CLI. The
// second field, HelpText, is used by the core CLI to display help information
// to the user in the core commands `cf help`, `cf`, or `cf -h`.
func (c *AnnotationsPlugin) GetMetadata() plugin.PluginMetadata {
	return plugin.PluginMetadata{
		Name: "Annotations",
		Version: plugin.VersionType{
			Major: 1,
			Minor: 0,
			Build: 0,
		},
		MinCliVersion: plugin.VersionType{
			Major: 6,
			Minor: 7,
			Build: 0,
		},
		Commands: []plugin.Command{
			{
				Name:     "set-annotation",
				HelpText: "Manipulate annotations",
				UsageDetails: plugin.Usage{
					Usage: "cf set-annotation resource-type resource-name annotations...\ncf unset-annotation resource-type resource-name annotation-names...\ncf annotations resource-type resource-name",
				},
			},
			{
				Name:     "unset-annotation",
				HelpText: "Manipulate annotations",
				UsageDetails: plugin.Usage{
					Usage: "cf set-annotation resource-type resource-name annotations...\ncf unset-annotation resource-type resource-name annotation-names...\ncf annotations resource-type resource-name",
				},
			},
			{
				Name:     "annotations",
				HelpText: "Manipulate annotations",
				UsageDetails: plugin.Usage{
					Usage: "cf set-annotation resource-type resource-name annotations...\ncf unset-annotation resource-type resource-name annotation-names...\ncf annotations resource-type resource-name",
				},
			},
		},
	}
}

// Unlike most Go programs, the `Main()` function will not be used to run all of the
// commands provided in your plugin. Main will be used to initialize the plugin
// process, as well as any dependencies you might require for your
// plugin.
func main() {
	// Any initialization for your plugin can be handled here
	//
	// Note: to run the plugin.Start method, we pass in a pointer to the struct
	// implementing the interface defined at "code.cloudfoundry.org/cli/plugin/plugin.go"
	//
	// Note: The plugin's main() method is invoked at install time to collect
	// metadata. The plugin will exit 0 and the Run([]string) method will not be
	// invoked.
	plugin.Start(new(AnnotationsPlugin))
	// Plugin code should be written in the Run([]string) method,
	// ensuring the plugin environment is bootstrapped.
}
