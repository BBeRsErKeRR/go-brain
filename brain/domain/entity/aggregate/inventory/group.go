package inventory

import (
	"fmt"
	"slices"
)

type AnsibleGroup struct {
	name   string
	Vars   map[string]any
	Childs []*AnsibleGroup
	Hosts  []*AnsibleHost
}

func NewAnsibleGroup(name string) *AnsibleGroup {
	return &AnsibleGroup{
		name:   name,
		Childs: []*AnsibleGroup{},
		Hosts:  []*AnsibleHost{},
	}
}

func (g *AnsibleGroup) GetName() string {
	return g.name
}

func (g *AnsibleGroup) GetChilds() []*AnsibleGroup {
	return g.Childs
}

func (g *AnsibleGroup) GetHosts() []*AnsibleHost {
	return g.Hosts
}

func (g *AnsibleGroup) GetAllHosts() []*AnsibleHost {
	var hosts []*AnsibleHost = slices.Clone(g.Hosts)
	for _, child := range g.Childs {
		childHosts := child.GetAllHosts()
		hosts = append(hosts, childHosts...)
	}
	return hosts
}

func (g *AnsibleGroup) AddChild(child *AnsibleGroup) {
	g.Childs = append(g.Childs, child)
}

func (g *AnsibleGroup) AddHost(host *AnsibleHost) {
	g.Hosts = append(g.Hosts, host)
}

func (g *AnsibleGroup) GetAllChildNames() []string {
	var names []string = make([]string, 0)
	for _, child := range g.Childs {
		names = append(names, child.GetName())
		names = append(names, getAllChildNames(child)...)
	}
	return names
}

func (g *AnsibleGroup) String() string {
	return fmt.Sprintf("AnsibleGroup(name=%s, childs=%s, hosts=%s)", g.name, g.Childs, g.Hosts)
}

// Helper function to recursively get all child group names from a given group.
func getAllChildNames(grp *AnsibleGroup) []string {

	var childs []*AnsibleGroup = grp.GetChilds()
	var childsNames []string

	for _, child := range childs {
		childsNames = append(childsNames, child.GetName())
		childsNames = append(childsNames, getAllChildNames(child)...)
	}

	return childsNames
}
