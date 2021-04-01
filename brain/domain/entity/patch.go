package entity

type PatchEntity struct {
	name      string
	obj       interface{}
	isRemoved bool
	variants  []string
}

func NewPatchEntity(name string, obj interface{}) *PatchEntity {
	return &PatchEntity{
		name:      name,
		obj:       obj,
		isRemoved: false,
		variants:  []string{},
	}
}

func (p *PatchEntity) SetRemoved(isRemoved bool) {
	p.isRemoved = isRemoved
}

func (p *PatchEntity) IsRemoved() bool {
	return p.isRemoved
}

func (p *PatchEntity) SetObj(obj interface{}) {
	p.obj = obj
}

func (p *PatchEntity) GetObj() interface{} {
	return p.obj
}

func (p *PatchEntity) GetName() string {
	return p.name
}

func (p *PatchEntity) SetVariants(groupNames []string) {
	p.variants = append(p.variants, groupNames...)
}

func (p *PatchEntity) GetVariants() []string {
	var result []string = make([]string, len(p.variants)+1)
	result[0] = p.name
	copy(result[1:], p.variants)
	return result
}
