package inventory

import (
	"encoding/json"
)

type meta struct {
	Hostvars map[string]map[string]any `json:"hostvars"`
}

type jgroup struct {
	Vars     map[string]any `json:"vars,omitempty"`
	Children []string       `json:"children,omitempty"`
	Hosts    []string       `json:"hosts,omitempty"`
}

func NewInventoryFromJson(jsonData []byte) (*AnsibleInventory, error) {
	var rawInventory map[string]json.RawMessage
	err := json.Unmarshal([]byte(jsonData), &rawInventory)
	if err != nil {
		return nil, err
	}

	var hostvars meta
	var headGroups []string
	groups := make(map[string]*AnsibleGroup)
	hosts := make(map[string]*AnsibleHost)
	for key, value := range rawInventory {
		switch {
		case "_meta" == key:
			err = json.Unmarshal(value, &hostvars)
			if err != nil {
				return nil, err
			}
		case "all" == key:
			err = json.Unmarshal(value, &headGroups)
			if err != nil {
				return nil, err
			}
		default:
			var jgrp jgroup
			err = json.Unmarshal(value, &jgrp)
			if err != nil {
				return nil, err
			}

			group, ok := groups[key]
			if !ok {
				group = &AnsibleGroup{
					name: key,
					vars: jgrp.Vars,
				}
				groups[key] = group
			}

			// Append childs

			// Append hosts

		}
	}

}
