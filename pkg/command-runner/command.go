package command_runner

import (
	"bytes"
	"cmd/internal/models"
	"cmd/pkg/json-check"
	"gorm.io/gorm"
	"os/exec"
	"regexp"
)

func runCommand(cmd string, t models.Task, db *gorm.DB) string {
	// See https://regexr.com/4154h for custom regex to parse commands
	// Inspired by https://gist.github.com/danesparza/a651ac923d6313b9d1b7563c9245743b
	pattern := `(--[^\s]+="[^"]+")|"([^"]+)"|'([^']+)'|([^\s]+)`
	parts := regexp.MustCompile(pattern).FindAllString(cmd, -1)

	//	The first part is the command, the rest are the args:
	head := parts[0]
	arguments := parts[1:len(parts)]

	//	Format the command
	command := exec.Command(head, arguments...)

	//	Sanity check -- capture stdout and stderr:
	var out bytes.Buffer
	var stderr bytes.Buffer
	command.Stdout = &out    // Standard out: out.String()
	command.Stderr = &stderr // Standard errors: stderr.String()

	//	Run the command
	command.Run()

	t.Status = "Completed"

	// Add results to db if in JSON format
	if json_check.IsJSON(out.String()) {
		t.Response = out.String()
	}

	db.Save(&t)

	return out.String()

}
