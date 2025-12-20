//go:build !solution

package ciletters

import (
	"bytes"
	_ "embed"
	"fmt"
	"strings"
	"text/template"
)

//go:embed template.txt
var pipelineTemplate string

func fetchLast10Lines(text string) []string {
	lines := strings.Split(text, "\n")
	if len(lines) < 10 {
		return lines
	}
	return lines[len(lines)-10:]
}

var funcMap template.FuncMap = template.FuncMap{
	"fetchLast10Lines": fetchLast10Lines,
}

func MakeLetter(n *Notification) (string, error) {
	var buffer bytes.Buffer

	tmpl := template.New("pipelineTemplate").Funcs(funcMap)

	tmpl, err := tmpl.Parse(pipelineTemplate)
	if err != nil {
		return "", fmt.Errorf("template.Parse: %w", err)
	}

	err = tmpl.Execute(&buffer, n)
	if err != nil {
		return "", fmt.Errorf("fillNotification: %w", err)
	}

	return buffer.String(), nil
}
