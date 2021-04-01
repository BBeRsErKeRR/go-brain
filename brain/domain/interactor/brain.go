package interactor

import (
	"github.com/bberserkerr/go-brain/brain/domain/entity"
	"github.com/bberserkerr/go-brain/brain/domain/entity/aggregate"
)

type BrainDifferI interface {
	GetPatch(request *aggregate.PatchRequest, strategy_request *entity.DiffEntitySelector) ([]*entity.PatchEntity, error)
	StorePatch(patch []*entity.PatchEntity) error
	ApplyPatch(request aggregate.ApplyPatchRequest) error
	RollbackPatch(request aggregate.ApplyPatchRequest) error
}
