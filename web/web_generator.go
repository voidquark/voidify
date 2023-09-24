package web_generator

import (
	_ "embed"
	"encoding/base64"
	"fmt"
	"os"
	"sort"
)

func Html(data map[string]map[string]map[string]interface{}, htmlFile string) error {
	htmlContent := generateHTML(data)
	filePath := htmlFile
	err := writeHTMLToFile(filePath, htmlContent)
	if err != nil {
		fmt.Printf("Error writing HTML file: %v\n", err)
		return err
	}
	return nil

}

func writeHTMLToFile(filePath, htmlContent string) error {
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(htmlContent)
	if err != nil {
		return err
	}

	if err := file.Chmod(0600); err != nil {
		return err
	}

	return nil
}

//go:embed voidify.png
var voidifyLogo []byte

func encodeLogo() (string, error) {
	base64Encoded := base64.StdEncoding.EncodeToString(voidifyLogo)

	return base64Encoded, nil
}

func generateHTML(data map[string]map[string]map[string]interface{}) string {
	base64Image, err := encodeLogo()
	if err != nil {
		fmt.Println("Error encoding image:", err)
	}
	html := `
<!DOCTYPE html>
<html>
<head>
    <meta name="viewport" content="width=device-width, initial-scale=1.0">`
	html += fmt.Sprintf(`
    <link rel="shortcut icon" href="data:image/png;base64,%s" type="image/png">
	`, base64Image)
	html += `
	<title>Voidify</title>
    <style>
		body {
			font-family: Inter,ui-sans-serif,system-ui,-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Helvetica Neue,Arial,Noto Sans,sans-serif,Apple Color Emoji,Segoe UI Emoji,Segoe UI Symbol,Noto Color Emoji;
			background-color: #111;
			color: #ffffff;
			padding: 20px;
		}
		h1 {
			color: #9518e2;
		}
        .header {
            display: flex;
            align-items: center;
            margin-bottom: 20px; /* Add space between logo and h1 */
        }
        .logo {
            max-width: 60px;
            margin-right: 20px;
        }
		.host-box {
			max-width: 80%;
			background-color: #3D244B;
			color: #ffffff;
			padding: 10px;
			margin: 10px 0;
			border-radius: 5px;
		}
		.code {
			background-color: #333333;
			color: #ffffff;
			padding: 5px;
			margin: 5px;
			border-radius: 5px;
			display: inline-block;
			font-family: 'Courier New', monospace;
		}
		.ssh-command {
			background-color: #3D244B;
			margin-right: 10px;
			padding: 5px;
			border-radius: 5px;
			display: inline-block;
			font-family: 'Courier New', monospace;
		}
		.ssh-code {
			background-color: #1e1e1e;
			color: #ffffff;
			padding: 5px;
			margin: 5px;
			border-radius: 5px;
			display: inline-block;
		}
		.copy-button {
			background-color: #9518e2;
			border: none;
			color: #ffffff;
			padding: 10px 20px;
			font-weight: bold;
			text-align: center;
			text-decoration: none;
			transition: background-color 0.3s ease, opacity 0.3s ease;
			border-radius: 5px;
			cursor: pointer;
			font-family: 'Courier New', monospace;
		}
		.copy-button:hover,
		.copy-button:focus {
			background-color: #555;
			opacity: 0.9;
		}
        .env-toggle {
			font-weight: bold;
            cursor: pointer;
            display: flex;
            align-items: center;
        }
        .arrow {
            margin-right: 5px;
            transition: transform 0.3s ease;
        }
        .env-content {
			font-weight: bold;
            display: none;
            margin-left: 20px;
			font-weight: bold;
        }
        .env-toggle.expanded + .env-content {
            display: block;
			font-weight: bold;
        }
		table {
			max-width: 80%; /* Adjust the maximum width as needed */
			overflow-x: auto;
			border-collapse: collapse;
			border-spacing: 0;
			margin-top: 20px;
		}
		th, td {
			font-weight: bold;
			padding: 8px 16px;
			text-align: left;
			border-bottom: 1px solid #ddd;
		}
		th {
			background-color: #333;
			color: #ffffff;
		}
    </style>
</head>
<body>
`
	// Sort environment names
	var envNames []string
	for envName := range data {
		envNames = append(envNames, envName)
	}
	sort.Strings(envNames)

	html += fmt.Sprintf(`
	<div class="header">
	<img class="logo" src="data:image/png;base64,%s" alt="Logo">
	<h1>Voidify</h1>
	</div>
	`, base64Image)

	// Generate the sorted table
	html += `
    <table>
        <tr>
            <th>Environment</th>
            <th>Number of Hosts</th>
        </tr>`
	for _, envName := range envNames {
		envData := data[envName]
		numHosts := len(envData["hosts"])
		html += fmt.Sprintf(`
        <tr>
            <td>%s</td>
            <td>%d</td>
        </tr>`, envName, numHosts)
	}

	html += `
	<br>
    </table>
	<br>`

	// Loop through each environment
	for _, envName := range envNames {
		envData := data[envName]

		// Sort host names within each environment
		var hostNames []string
		for hostName := range envData["hosts"] {
			hostNames = append(hostNames, hostName)
		}
		sort.Strings(hostNames)

		html += fmt.Sprintf(`<div class="env-toggle" onclick="toggleEnv(this)"><span class="arrow">➡️</span> %s</div>`, envName)
		html += fmt.Sprintf(`<div class="env-content">`)

		// Loop through each host in the environment
		for _, hostName := range hostNames {
			hostProps := envData["hosts"][hostName]

			html += `<div class="host-box">`
			html += fmt.Sprintf("<h3>%s</h3>", hostName)
			html += fmt.Sprintf(`<button class="copy-button" onclick="copyToClipboard('` + hostName + `')">📋 Copy SSH Command</button>`)
			html += fmt.Sprintf(`<div class="ssh-command"><div class="ssh-code" id="ssh-command-`+hostName+`">ssh %s</div></div>`, hostName)
			for propName, propValue := range hostProps.(map[string]interface{}) {
				propStringValue := fmt.Sprintf("%v", propValue)
				html += fmt.Sprintf(`<div class="code">%s %s</div>`, propName, propStringValue)
			}
			html += `</div>`
		}
		html += `</div>`
	}

	html += `
    <script>
        function toggleEnv(envToggle) {
            var envContent = envToggle.nextElementSibling;
            var arrow = envToggle.querySelector('.arrow');
            if (envContent.style.display === 'none' || !envContent.style.display) {
                envContent.style.display = 'block';
                envToggle.classList.add('expanded');
                arrow.innerText = '⬇️';
            } else {
                envContent.style.display = 'none';
                envToggle.classList.remove('expanded');
                arrow.innerText = '➡️';
            }
        }

        function copyToClipboard(hostName) {
            var sshCommandElement = document.getElementById('ssh-command-' + hostName);
            var textArea = document.createElement("textarea");
            textArea.value = sshCommandElement.innerText;
            document.body.appendChild(textArea);
            textArea.select();
            document.execCommand('copy');
            document.body.removeChild(textArea);
        }
    </script>
</body>
</html>
`

	return html
}
