package inventory

import (
	"fmt"
	"slices"
)

type AnsibleGroup struct {
	name   string          `json:"-"`
	vars   map[string]any  `json:"vars,omitempty"`
	childs []*AnsibleGroup `json:"children,omitempty"`
	hosts  []*AnsibleHost  `json:"hosts,omitempty"`
}

func NewAnsibleGroup(name string) *AnsibleGroup {
	return &AnsibleGroup{
		name:   name,
		childs: []*AnsibleGroup{},
		hosts:  []*AnsibleHost{},
	}
}

func (g *AnsibleGroup) GetName() string {
	return g.name
}

func (g *AnsibleGroup) GetChilds() []*AnsibleGroup {
	return g.childs
}

func (g *AnsibleGroup) GetHosts() []*AnsibleHost {
	return g.hosts
}

func (g *AnsibleGroup) GetAllHosts() []*AnsibleHost {
	var hosts []*AnsibleHost = slices.Clone(g.hosts)
	for _, child := range g.childs {
		childHosts := child.GetAllHosts()
		hosts = append(hosts, childHosts...)
	}
	return hosts
}

func (g *AnsibleGroup) AddChild(child *AnsibleGroup) {
	g.childs = append(g.childs, child)
}

func (g *AnsibleGroup) AddHost(host *AnsibleHost) {
	g.hosts = append(g.hosts, host)
}

func (g *AnsibleGroup) GetAllChildNames() []string {
	var names []string = make([]string, 0)
	for _, child := range g.childs {
		names = append(names, child.GetName())
		names = append(names, getAllChildNames(child)...)
	}
	return names
}

func (g *AnsibleGroup) String() string {
	return fmt.Sprintf("AnsibleGroup(name=%s, childs=%s, hosts=%s)", g.name, g.childs, g.hosts)
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
