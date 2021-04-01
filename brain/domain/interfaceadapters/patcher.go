package interfaceadapters

import (
	"context"

	"github.com/bberserkerr/go-brain/brain/domain/entity"
)

type PatcherI interface {
	Load(path string) ([]entity.PatchEntity, error)
	Init(path string) error
	Migrate(ctx context.Context, patch entity.PatchEntity, path string) error
	Upgrade(ctx context.Context, patch entity.PatchEntity, path string) error
}

type PatchFilterI interface {
	Filter(*entity.PatchEntity) error
}
