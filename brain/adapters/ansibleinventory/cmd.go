package ansibleinventory

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/bberserkerr/go-brain/brain/domain/entity/aggregate/inventory"
	brainerrors "github.com/bberserkerr/go-brain/brain/domain/errors"
	"github.com/bberserkerr/go-brain/brain/domain/interfaceadapters"
	"github.com/bberserkerr/go-brain/pkg/command"
)

var parsedInventoryPattern = regexp.MustCompile(`(?mi)^((.*\n)+?)(\s*(Parsed|\[WARNING)(.*\n)+.*)$`)

func parseCmd(logger interfaceadapters.Logger, stdOut, stdErr []byte) (*inventory.AnsibleInventory, error) {
	stderr := string(stdErr)
	output := string(stdOut)
	bracketPos := strings.Index(output, "\n{")
	if bracketPos == -1 {
		return nil, &brainerrors.AnsibleInventoryException{
			Message: "Unable to extract json from ansible-inventory output",
		}
	}
	jsonSubstring := output[bracketPos:]
	matches := parsedInventoryPattern.FindStringSubmatch(stderr)
	var infoMessage, warningMessage string
	if len(matches) > 0 {
		warningMessage = strings.Replace(matches[3], "\\n", "\n", -1)
		infoMessage = strings.Replace(output[:bracketPos-1]+"\n"+matches[1], "\\n", "\n", -1)
	} else {
		infoMessage = output[:bracketPos-1]
	}

	if infoMessage != "" {
		logger.Debug(fmt.Sprintf("Command Info messages:\n %s", infoMessage))
	}

	if warningMessage != "" && (strings.Contains(warningMessage, "warning") || strings.Contains(warningMessage, "WARNING")) {
		logger.Warn(fmt.Sprintf("Command Warning messages:\n %s", warningMessage))
	}

	return inventory.NewInventoryFromJson(jsonSubstring)
}

// TODO: resolve short/user path; trim paths
func NewCMDAnsibleInventory(ctx context.Context, logger interfaceadapters.Logger, paths []string, limits string, chd string) (*inventory.AnsibleInventory, error) {
	cmd := "ansible-inventory"
	args := make([]string, 0, 2*len(paths)+4)
	for _, p := range paths {
		args = append(args, "-i", p)
	}
	args = append(args, "--list", "-vvvv")
	if limits != "" {
		args = append(args, "-l", limits)
	}
	envVars := map[string]string{
		"ANSIBLE_INVENTORY_UNPARSED_FAILED": "true",
	}
	stdOut, stdErr, err := command.WithContext(ctx, cmd, args, envVars, chd)
	if err != nil {
		return nil, err
	}

	ansibleInventory, err := parseCmd(logger, stdOut, stdErr)
	if err != nil {
		return nil, err
	}

	return ansibleInventory, nil
}
