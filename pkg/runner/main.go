package runner

import (
	"bytes"
	"cmd/internal/models"
	"cmd/pkg/json-manager"
	logger "cmd/pkg/zap-logger"
	"gorm.io/gorm"
	"os/exec"
	"regexp"
	"strconv"
)

type Runner struct {
	DB *gorm.DB
}

func (r Runner) RunCommand(cmd string, t models.Task) string {
	pattern := `(--[^\s]+="[^"]+")|"([^"]+)"|'([^']+)'|([^\s]+)`
	parts := regexp.MustCompile(pattern).FindAllString(cmd, -1)

	//	The first part is the command, the rest are the args:
	head := parts[0]
	arguments := parts[1:]

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
		logger.GetLogger().Fatal(strconv.Itoa(int(t.ID)) + " command execution fail")
	}

	t.Status = "Completed"

	// Add results to db if in JSON format
	if json_manager.IsJSON(out.String()) {
		t.Response = out.String()
	}

	r.DB.Save(&t)

	return out.String()
}
