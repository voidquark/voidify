package yaml_config

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

func ReadYAMLConfig(configFile string) (map[string]map[string]map[string]interface{}, error) {
	yamlFile, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}

	var data map[string]map[string]map[string]interface{}

	err = yaml.Unmarshal(yamlFile, &data)
	if err != nil {
		return nil, err
	}

	return data, nil
}

func GenerateSSHConfig(data map[string]map[string]map[string]interface{}, sshConfigFile string) error {
	var sshConfigBuilder strings.Builder

	// Loop through the YAML data and build SSH config entries
	for _, hosts := range data {
		for host, properties := range hosts["hosts"] {
			hostEntry := fmt.Sprintf("Host %s\n", host)
			sshConfigBuilder.WriteString(hostEntry)

			// Write host properties
			for key, value := range properties.(map[string]interface{}) {
				key = strings.TrimSpace(key)
				valueStr := strings.TrimSpace(fmt.Sprintf("%v", value))
				propertyEntry := fmt.Sprintf("\t%s %s\n", key, valueStr)
				sshConfigBuilder.WriteString(propertyEntry)
			}

			sshConfigBuilder.WriteString("\n")
		}
	}

	// Write the generated SSH config to the specified file
	err := os.WriteFile(sshConfigFile, []byte(sshConfigBuilder.String()), 0600)
	if err != nil {
		return err
	}

	return nil
}
