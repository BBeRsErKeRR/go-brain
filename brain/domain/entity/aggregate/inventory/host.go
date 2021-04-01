package inventory

import (
	"fmt"
)

type AnsibleHost struct {
	name string `json:"-"`
	vars map[string]any
}

func NewAnsibleHost(name string) *AnsibleHost {
	return &AnsibleHost{
		name: name,
	}
}

func (h *AnsibleHost) GetName() string {
	return h.name
}

func (h *AnsibleHost) GetVars() map[string]any {
	return h.vars
}

func (h *AnsibleHost) SetVar(key string, value any) {
	h.vars[key] = value
}

func (h *AnsibleHost) SetVars(varsMap map[string]any) {
	h.vars = varsMap
}

func (h *AnsibleHost) String() string {
	return fmt.Sprintf("AnsibleHost(name=%s, vars=%v)", h.name, h.vars)
}
