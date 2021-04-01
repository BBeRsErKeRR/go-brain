package interfaceadapters

type HostVarsFilterI interface {
	Filter(hostvars map[string]any) map[string]any
}
