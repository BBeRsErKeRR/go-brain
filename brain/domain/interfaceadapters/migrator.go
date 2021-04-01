package interfaceadapters

import (
	"context"

	"github.com/bberserkerr/go-brain/brain/domain/entity/aggregate"
)

type MigratorI interface {
	Upgrade(ctx context.Context, request aggregate.ApplyPatchRequest) error
	Rollback(ctx context.Context, request aggregate.ApplyPatchRequest) error
}
