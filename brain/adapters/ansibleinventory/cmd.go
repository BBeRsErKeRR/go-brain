package ansibleinventory

import (
	"context"

	"github.com/bberserkerr/go-brain/brain/domain/entity/aggregate/inventory"
	"github.com/bberserkerr/go-brain/brain/domain/interfaceadapters"
	"github.com/bberserkerr/go-brain/pkg/command"
)

func parseCmd(logger interfaceadapters.Logger, stdOut, stdErr []byte) (inventory.AnsibleInventory, error) {
	return inventory.AnsibleInventory{}, nil
}

// TODO: resolve short/user path; trim paths
func NewCMDAnsibleInventory(ctx context.Context, logger interfaceadapters.Logger, paths []string, chd string) (*inventory.AnsibleInventory, error) {
	cmd := "ansible-inventory"
	args := make([]string, 0, 2*len(paths)+2)
	for _, p := range paths {
		args = append(args, "-i", p)
	}
	args = append(args, "--list", "-vvvv")
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

	return &ansibleInventory, nil
}
