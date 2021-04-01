package aggregate

import "github.com/bberserkerr/go-brain/brain/domain/entity"

type DiffEntitySelectorResponse struct {
	DonorDiffRequest     entity.DiffRequest
	RecipientDiffRequest entity.DiffRequest
}
