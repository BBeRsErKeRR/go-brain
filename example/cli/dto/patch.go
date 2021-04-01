package dto

import "github.com/bberserkerr/go-brain/brain/domain/entity"

type PatchDTO struct {
	name      string
	obj       interface{}
	isRemoved bool
	variants  []string
}

func NewPatch(patch entity.PatchEntity) *PatchDTO {
	return &PatchDTO{
		name:      patch.GetName(),
		obj:       patch.GetObj(),
		isRemoved: patch.IsRemoved(),
		variants:  patch.GetVariants(),
	}
}
