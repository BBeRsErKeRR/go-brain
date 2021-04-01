package hostvarsfilter

import "github.com/bberserkerr/go-brain/brain/domain/interfaceadapters"

type Multiple struct {
	filters []interfaceadapters.HostVarsFilterI
}

func (f *Multiple) Filter(vars map[string]any) map[string]any {
	res := vars
	for _, filter := range f.filters {
		res = filter.Filter(res)
	}
	return res
}
