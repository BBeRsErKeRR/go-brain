package interfaceadapters

import (
	"context"

	"github.com/bberserkerr/go-brain/brain/domain/entity"
)

type PatchStorageI interface {
	Connect(context.Context) error
	Close(context.Context) error
	StorePatch(context.Context, *entity.PatchEntity) error
	StorePatches(context.Context, []*entity.PatchEntity) error
	GetPatch(context.Context, string) (*entity.PatchEntity, error)
	GetPatches(context.Context, string) ([]*entity.PatchEntity, error)
}
