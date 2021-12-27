package command_runner

import (
	"bytes"
	"cmd/internal/models"
	"cmd/pkg/json-check"
	zap_logger "cmd/pkg/zap-logger"
	"gorm.io/gorm"
	"os/exec"
	"regexp"
	"strconv"
)

func RunCommand(cmd string, t models.Task, db *gorm.DB) string {
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
	err := command.Run()

	if err != nil {
		zap_logger.GetLogger().Fatal(strconv.Itoa(int(t.ID)) + " command execution fail")
	}

	t.Status = "Completed"

	// Add results to db if in JSON format
	if json_check.IsJSON(out.String()) {
		t.Response = out.String()
	}

	db.Save(&t)

	return out.String()
}
