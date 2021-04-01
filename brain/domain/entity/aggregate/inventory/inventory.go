package inventory

import (
	"maps"
	"slices"
)

type AnsibleInventory struct {
	jsonStruct map[string]any
	headGroups []string
	groups     map[string]*AnsibleGroup `json:"-"`
}

func (i *AnsibleInventory) GetGroupsNames() []string {
	return slices.Collect(maps.Keys(i.groups))
}

func (i *AnsibleInventory) GetInventoryStruct() map[string]any {
	return i.jsonStruct
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
