package api

import (
	"context"

	"github.com/bberserkerr/go-brain/brain/domain/interactor"
	"github.com/bberserkerr/go-brain/brain/domain/interface_adapters"
	"github.com/bberserkerr/go-brain/v1/api/cfg"
	"github.com/bberserkerr/go-brain/v1/api/dto"
)

type BrainApi struct {
	u interactor.BrainDifferI
}

func (b *BrainApi) GetPatch() (*dto.PatchDTO, error) {
	patch, err := b.u.GetPatch()
	if err != nil {
		return nil, err
	}
	return dto.NewPatch(patch), err
}

func NewApi(ctx context.Context, logger interface_adapters.Logger, cfg cfg.Config) *BrainApi {
	return nil
}
