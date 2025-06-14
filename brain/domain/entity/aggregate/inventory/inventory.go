package inventory

import (
	"encoding/json"
	"fmt"
	"maps"
	"slices"
)

type AnsibleInventory struct {
	headGroups []string
	groups     map[string]*AnsibleGroup `json:"-"`
}

func (i *AnsibleInventory) String() string {
	return fmt.Sprintf("AnsibleInventory(headGroups=%s, groups=%s)", i.headGroups, i.groups)
}

func (i *AnsibleInventory) GetGroupsNames() []string {
	return slices.Collect(maps.Keys(i.groups))
}

func (i *AnsibleInventory) MarshalJSON() ([]byte, error) {
	res := make(map[string]any)
	for k, v := range i.groups {
		res[k] = v
	}
	res["_meta"] = &meta{}

	return json.Marshal(res)
}

// TODO: migrate to MarshalJSON()
func (i *AnsibleInventory) GetInventoryStruct() map[string]any {
	result := make(map[string]any)
	return result
}

func (i *AnsibleInventory) GetHeadGroups() []string {
	return i.headGroups
}

func (i *AnsibleInventory) GetGroups(groupNames []string) []*AnsibleGroup {
	res := make([]*AnsibleGroup, 0, len(groupNames))
	for _, name := range groupNames {
		g, ok := i.groups[name]
		if ok {
			res = append(res, g)
		}
	}
	return res
}

func (i *AnsibleInventory) GetGroup(groupName string) *AnsibleGroup {
	return i.groups[groupName]
}

func (i *AnsibleInventory) GetAllGroups() []*AnsibleGroup {
	return slices.Collect(maps.Values(i.groups))
}

func (i *AnsibleInventory) GetHostsByGroup(groupName string) []*AnsibleHost {
	group, ok := i.groups[groupName]
	if !ok {
		return []*AnsibleHost{}
	}
	return group.GetAllHosts()
}

// func NewAnsibleInventory(jsonStruct map[string]any) *AnsibleInventory {
// 	ignoredKeys := []string{"_meta", "all", "ungrouped"}
// 	hostVars := jsonStruct["_meta"]["hostvars"]
// 	return nil
// }
