package menu

import (
	"fmt"
	"os"
	"os/exec"
	"sort"
	"strings"

	"github.com/AlecAivazis/survey/v2"
)

func InteractiveHostSelection(data map[string]map[string]map[string]interface{}) error {
	//Extract all env
	var allEnv []string
	for env := range data {
		allEnv = append(allEnv, env)
	}
	sort.Strings(allEnv)

	// Prompt the user to select an env
	var selectedEnv string
	promptEnv := &survey.Select{
		Message: "Select Environment:",
		Options: allEnv,
	}
	if err := survey.AskOne(promptEnv, &selectedEnv, survey.WithPageSize(35), survey.WithIcons(func(icons *survey.IconSet) {
		icons.Question.Text = "--"
		icons.Question.Format = "magenta+hb"
		icons.SelectFocus.Format = "green+hb"
	})); err != nil {
		return err
	}

	// Extract all hosts
	hosts, ok := data[selectedEnv]["hosts"]
	if !ok {
		return fmt.Errorf("Environment not found: %s", selectedEnv)
	}
	var allHosts []string
	for host := range hosts {
		allHosts = append(allHosts, host)
	}
	sort.Strings(allHosts)

	// Prompt the user to select a host
	var selectedHost string
	promptHost := &survey.Select{
		Message: "Select Host:",
		Options: append(allHosts, "-->RETURN TO ENVIRONMENT SELECTION"),
	}
	if err := survey.AskOne(promptHost, &selectedHost, survey.WithPageSize(35), survey.WithIcons(func(icons *survey.IconSet) {
		icons.Question.Text = "--"
		icons.Question.Format = "magenta+hb"
		icons.SelectFocus.Format = "blue+hb"
	})); err != nil {
		return err
	}

	if strings.Contains(selectedHost, "-->RETURN TO ENVIRONMENT SELECTION") {
		if err := clearScreen(); err != nil {
			fmt.Println("Failed to clear screen", err)
		}
		if err := InteractiveHostSelection(data); err != nil {
			fmt.Println("Error selecting host:", err)
			os.Exit(1)
		}
	} else {
		sshcmd := exec.Command("ssh", selectedHost)
		sshcmd.Stdin = os.Stdin
		sshcmd.Stdout = os.Stdout
		sshcmd.Stderr = os.Stderr
		sshcmd.Run()
	}

	return nil
}

func clearScreen() error {
	cmd := exec.Command("clear")
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
