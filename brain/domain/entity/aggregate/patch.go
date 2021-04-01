package aggregate

import "github.com/bberserkerr/go-brain/brain/domain/entity"

type PatchRequest struct {
	DonorInventoryRequest     entity.InventoryRequest
	RecipientInventoryRequest entity.InventoryRequest
}

type ApplyPatchRequest struct {
	Target    entity.InventoryRequest
	Patch     []*entity.PatchEntity
	CanEmpty  bool
	StorePath string
}
