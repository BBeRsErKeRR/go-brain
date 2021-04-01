package service

import (
	"context"
	"fmt"

	"github.com/bberserkerr/go-brain/brain/domain/entity"
	"github.com/bberserkerr/go-brain/brain/domain/entity/aggregate"
	"github.com/bberserkerr/go-brain/brain/domain/entity/aggregate/inventory"
	"github.com/bberserkerr/go-brain/brain/domain/interactor"
	"github.com/bberserkerr/go-brain/brain/domain/interfaceadapters"
	"github.com/bberserkerr/go-brain/pkg/objects"
)

type BrainDiffer struct {
	logger                           interfaceadapters.Logger
	ctx                              context.Context
	inventoryBuilderService          interfaceadapters.AnsibleInventoryBuilderI
	hostvarsFilterService            interfaceadapters.HostVarsFilterI
	diffService                      interfaceadapters.DifferI
	diffEntitySelectorServiceBuilder interfaceadapters.DiffEntitySelectorStrategyBuilderI
	patchStorage                     interfaceadapters.PatchStorageI
	patchFilterService               interfaceadapters.PatchFilterI
	migrator                         interfaceadapters.MigratorI
}

// Brain implements the BrainI
var _ interactor.BrainDifferI = &BrainDiffer{}

func (b *BrainDiffer) getInventory(inventoryRequest *entity.InventoryRequest) (*inventory.AnsibleInventory, error) {
	b.logger.Debug(fmt.Sprintf("get inventory: %s", inventoryRequest))
	return b.inventoryBuilderService.Build(inventoryRequest)
}

func (b *BrainDiffer) getDiffApplyVariants(dGrp string, dI, rI *inventory.AnsibleInventory, isRInD bool) []string {
	var res []string
	if !isRInD {
		return res
	}

	dGroups := dI.GetGroup(dGrp).GetAllChildNames()
	patchGroup := rI.GetGroup(dGrp)
	rGroups := patchGroup.GetAllChildNames()
	allRHosts := patchGroup.GetAllHosts()
	intersect := objects.Intersection(dGroups, rGroups)

	res = make([]string, 0, len(intersect))
	for _, grpName := range intersect {
		variantGrp := rI.GetGroup(grpName)
		if objects.EqSlices(allRHosts, variantGrp.GetAllHosts()) {
			res = append(res, grpName)
		}
	}
	return res
}

func (b *BrainDiffer) getDiffEntity(inv *inventory.AnsibleInventory, grp *inventory.AnsibleGroup) *entity.DIffEntity {
	hosts := inv.GetHostsByGroup(grp.GetName())
	var target map[string]any
	if len(hosts) > 0 {
		target = b.hostvarsFilterService.Filter(hosts[0].GetVars())
	} else {
		target = make(map[string]any)
	}
	return entity.NewDIffEntity(grp.GetName(), target)
}

func (b *BrainDiffer) generateDiffEntities(req *entity.DiffRequest, inv *inventory.AnsibleInventory) map[string]*entity.DIffEntity {
	res := make(map[string]*entity.DIffEntity)
	groups := inv.GetGroups(req.Groups)
	for _, g := range groups {
		res[g.GetName()] = b.getDiffEntity(inv, g)
	}
	return res
}

func (b *BrainDiffer) GetPatch(request *aggregate.PatchRequest, strategy_request *entity.DiffEntitySelector) ([]*entity.PatchEntity, error) {
	dI, err := b.getInventory(&request.DonorInventoryRequest)
	if err != nil {
		return nil, err
	}
	rI, err := b.getInventory(&request.RecipientInventoryRequest)
	if err != nil {
		return nil, err
	}

	// 1. Build compare strategy
	// TODO: Add strategy to get compared groups (head groups or groups from playbook)
	selectStrategy := b.diffEntitySelectorServiceBuilder.Build(strategy_request)

	// 2. Get diff request entities for donor and recipient inventories
	selectRequest, err := selectStrategy.GetDiffRequest(dI, rI)
	if err != nil {
		return nil, err
	}

	// 3. Get diff entities for donor and recipient inventories
	donorEntities := b.generateDiffEntities(&selectRequest.DonorDiffRequest, dI)
	recipientEntities := b.generateDiffEntities(&selectRequest.RecipientDiffRequest, rI)

	// TODO: recipient can contains fake head groups, for override settings
	// We need collect all groups for recipient inventory

	res := make([]*entity.PatchEntity, 0, len(donorEntities))

	// 4. Compare the donor and recipient inventories
	for dGrp, dDiff := range donorEntities {
		dCompareValue := dDiff
		rCompareValue, isRInD := recipientEntities[dGrp]

		if !isRInD {
			rCompareValue = entity.NewDIffEntity(dGrp, make(map[string]any))
		}

		patch, err := b.diffService.GetPatch(dCompareValue, rCompareValue)
		if err != nil {
			return nil, err
		}

		// If patch not exist (eq objects) continue
		if patch == nil {
			continue
		}

		err = b.patchFilterService.Filter(patch)
		if err != nil {
			return nil, err
		}

		// Add possible variants to apply
		patch.SetVariants(b.getDiffApplyVariants(dGrp, dI, rI, isRInD))
		res = append(res, patch)
	}

	for rGrp, rDiff := range recipientEntities {
		_, ok := donorEntities[rGrp]
		if ok {
			continue
		}
		patch, err := b.diffService.GetPatch(entity.NewDIffEntity(rGrp, make(map[string]any)), rDiff)
		if err != nil {
			return nil, err
		}
		patch.SetRemoved(true)
		// TODO: calculate res length before
		res = append(res, patch)
	}

	return res, nil
}

func (b *BrainDiffer) StorePatch(patch []*entity.PatchEntity) error {
	return b.patchStorage.StorePatches(b.ctx, patch)
}

func (b *BrainDiffer) LoadPatch(path string) ([]*entity.PatchEntity, error) {
	return b.patchStorage.GetPatches(b.ctx, path)
}

func (b *BrainDiffer) ApplyPatch(request aggregate.ApplyPatchRequest) error {
	return b.migrator.Upgrade(b.ctx, request)
}

func (b *BrainDiffer) RollbackPatch(request aggregate.ApplyPatchRequest) error {
	return b.migrator.Rollback(b.ctx, request)
}
