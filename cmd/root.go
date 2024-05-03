package cmd

import (
	"fmt"
	"os"
	"os/user"

	"github.com/spf13/cobra"
	menu "github.com/voidquark/voidify/interactive_menu"
	web_generator "github.com/voidquark/voidify/web"
	yaml_config "github.com/voidquark/voidify/yaml"
)

const CLIVersion string = "1.0.3"

var configFile string
var sshConfigFile string
var htmlFile string

var rootCmd = &cobra.Command{
	Use:     "voidify",
	Short:   "Simplify and Fastify your SSH Management",
	Version: CLIVersion,
	Long: `
┓┏  • ┓•┏
┃┃┏┓┓┏┫┓╋┓┏
┗┛┗┛┗┗┻┗┛┗┫
	  ┛

Voidify simplifies and accelerates SSH management, eliminating the need to deal with bash auto-completions.
With Voidify, you don't have to worry about remembering all the server details.
Instead, just run Voidify, use your arrow keys in the terminal to navigate through environment selections, and choose the server name you want to connect to.
You can even start typing to filter hosts while making your selection.
It takes inspiration from Ansible's YAML-based inventory to simplify configuration, which is automatically translated into SSH config.`,
	Run: func(cmd *cobra.Command, args []string) {
	},
}

func Execute() {
	rootCmd.Flags().StringVarP(&configFile, "config-file", "c", "", "specify the path to the YAML inventory file (required)")
	rootCmd.Flags().StringVarP(&sshConfigFile, "out-ssh-config-file", "o", "", "specify the path to the SSH config file (default: $HOME/.ssh/config)")
	rootCmd.Flags().StringVarP(&htmlFile, "web-html-file", "w", "", "optionally generate a static HTML website, specify the file path, including the file name, e.g., /tmp/voidify.html (Not generated by default)")
	rootCmd.MarkFlagRequired("config-file")

	if err := rootCmd.Execute(); err != nil {
		os.Exit(0)
	}

	if rootCmd.Flags().Changed("help") {
		os.Exit(0)
	}

	if rootCmd.Flags().Changed("version") {
		os.Exit(0)
	}

	if sshConfigFile == "" {
		usr, err := user.Current()
		if err != nil {
			fmt.Println("Error getting home directory:", err)
		}
		sshConfigFile = usr.HomeDir + "/.ssh/config"
	}

	// Read YAML
	data, err := yaml_config.ReadYAMLConfig(configFile)
	if err != nil {
		fmt.Println("Error reading YAML config:", err)
		os.Exit(1)
	}

	// Generate SSH Config
	err = yaml_config.GenerateSSHConfig(data, sshConfigFile)
	if err != nil {
		fmt.Println("Error generating SSH config:", err)
		os.Exit(1)
	}

	// Generate HTML Website
	if rootCmd.Flags().Changed("web-html-file") {
		if err = web_generator.Html(data, htmlFile); err != nil {
			fmt.Println("Error generating website:", err)
		}
	}

	// Call InteractiveHostSelection to interactively select a host
	if err = menu.InteractiveHostSelection(data); err != nil {
		fmt.Println("Error selecting host:", err)
	}

}
