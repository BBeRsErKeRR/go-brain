package inventory

import (
	"encoding/json"
	"fmt"
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
		fmt.Println("DEBUG")
		return nil, err
	}

	var hostVarsMeta meta
	var headGroups []string
	groups := make(map[string]*AnsibleGroup)
	hosts := make(map[string]*AnsibleHost)

	// first them all unmarshal hostvars
	err = json.Unmarshal(rawInventory["_meta"], &hostVarsMeta)
	if err != nil {
		return nil, err
	}

	for key, value := range rawInventory {
		switch {
		case key == "_meta":
			continue
		case key == "all":
			var allGrp jgroup
			err = json.Unmarshal(value, &allGrp)
			if err != nil {
				return nil, err
			}
			group := &AnsibleGroup{
				name: key,
				Vars: allGrp.Vars,
			}
			groups[key] = group
			headGroups = allGrp.Children
		default:
			var jsonGrp jgroup
			err = json.Unmarshal(value, &jsonGrp)
			if err != nil {
				return nil, err
			}

			group, ok := groups[key]
			if !ok {
				group = &AnsibleGroup{
					name: key,
					Vars: jsonGrp.Vars,
				}
				groups[key] = group
			}

			// Append childs
			for _, child := range jsonGrp.Children {
				childGroup, ok := groups[child]
				if !ok {
					childGroup = &AnsibleGroup{
						name: child,
					}
					groups[key] = childGroup
				}
				group.AddChild(childGroup)
			}

			// Append hosts
			// TODO: add pointer into host.hostvars
			for _, hostName := range jsonGrp.Hosts {
				host, ok := hosts[hostName]
				if !ok {
					host = &AnsibleHost{
						name: hostName,
					}
					hostVars, ok := hostVarsMeta.Hostvars[hostName]
					if ok {
						host.SetVars(hostVars)
					}
				}
				group.AddHost(host)
			}

		}
	}

	return &AnsibleInventory{
		headGroups: headGroups,
		groups:     groups,
	}, nil

}
