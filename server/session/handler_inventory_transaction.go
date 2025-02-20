package session

import (
	"fmt"

	"github.com/df-mc/dragonfly/server/block/cube"
	"github.com/df-mc/dragonfly/server/event"
	"github.com/df-mc/dragonfly/server/item"
	"github.com/df-mc/dragonfly/server/item/inventory"
	"github.com/df-mc/dragonfly/server/world"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
)

// InventoryTransactionHandler handles the InventoryTransaction packet.
type InventoryTransactionHandler struct{}

// Handle ...
func (h *InventoryTransactionHandler) Handle(p packet.Packet, s *Session, tx *world.Tx, c Controllable) error {
	pk := p.(*packet.InventoryTransaction)

	switch data := pk.TransactionData.(type) {
	case *protocol.NormalTransactionData:
		h.resendInventories(s)
		// Always resend inventories with normal transactions. Most of the time we do not use these
		// transactions, so we're best off making sure the client and server stay in sync.
		if err := h.handleNormalTransaction(pk, s, c); err != nil {
			s.conf.Log.Debug("process packet: InventoryTransaction: verify Normal transaction actions: " + err.Error())
			return nil
		}
		return nil
	case *protocol.MismatchTransactionData:
		// Just resend the inventory and don't do anything.
		h.resendInventories(s)
		return nil
	case *protocol.UseItemOnEntityTransactionData:
		if err := s.VerifyAndSetHeldSlot(int(data.HotBarSlot), s.stackToItem(data.HeldItem.Stack), c); err != nil {
			return err
		}
		return h.handleUseItemOnEntityTransaction(data, s, tx, c)
	case *protocol.UseItemTransactionData:
		if err := s.VerifyAndSetHeldSlot(int(data.HotBarSlot), s.stackToItem(data.HeldItem.Stack), c); err != nil {
			return err
		}
		return h.handleUseItemTransaction(data, s, c)
	case *protocol.ReleaseItemTransactionData:
		if err := s.VerifyAndSetHeldSlot(int(data.HotBarSlot), s.stackToItem(data.HeldItem.Stack), c); err != nil {
			return err
		}
		return h.handleReleaseItemTransaction(c)
	}
	return fmt.Errorf("unhandled inventory transaction type %T", pk.TransactionData)
}

// resendInventories resends all inventories of the player.
func (h *InventoryTransactionHandler) resendInventories(s *Session) {
	s.sendInv(s.inv, protocol.WindowIDInventory)
	s.sendInv(s.ui, protocol.WindowIDUI)
	s.sendInv(s.offHand, protocol.WindowIDOffHand)
	s.sendInv(s.armour.Inventory(), protocol.WindowIDArmour)
}

// handleNormalTransaction ...
func (h *InventoryTransactionHandler) handleNormalTransaction(pk *packet.InventoryTransaction, s *Session, c Controllable) error {
	if len(pk.Actions) != 2 {
		return fmt.Errorf("expected two actions for dropping an item, got %d", len(pk.Actions))
	}

	var (
		slot     int
		count    int
		expected item.Stack
	)
	for _, action := range pk.Actions {
		if action.SourceType == protocol.InventoryActionSourceWorld && action.InventorySlot == 0 {
			if old := s.stackToItem(action.OldItem.Stack); !old.Empty() {
				return fmt.Errorf("unexpected non-empty old item in transaction action: %#v", action.OldItem)
			}
			count = int(action.NewItem.Stack.Count)
		} else if action.SourceType == protocol.InventoryActionSourceContainer && action.WindowID == protocol.WindowIDInventory {
			if expected = s.stackToItem(action.OldItem.Stack); expected.Empty() {
				return fmt.Errorf("unexpected empty old item in transaction action: %#v", action.OldItem)
			}
			slot = int(action.InventorySlot)
		} else {
			return fmt.Errorf("unexpected action type in drop item transaction")
		}
	}

	actual, _ := s.inv.Item(slot)
	if count < 1 {
		return fmt.Errorf("expected at least one item to be dropped, got %d", count)
	}
	if count > actual.Count() {
		return fmt.Errorf("tried to throw %v items, but held only %v in slot", count, actual.Count())
	}
	if !expected.Equal(actual) {
		return fmt.Errorf("different item thrown than held in slot: %#v was thrown but held %#v", expected, actual)
	}

	// Explicitly don't re-use the thrown variable. This item was supplied by the user, and if some
	// logic in the Comparable() method was flawed, users would be able to cheat with item properties.
	// Only grow or shrink the held item to prevent any such issues.
	res := actual.Grow(count - actual.Count())
	if err := call(event.C(inventory.Holder(c)), int(*s.heldSlot), res, s.inv.Handler().HandleDrop); err != nil {
		return err
	}

	n := c.Drop(res)
	_ = s.inv.SetItem(slot, actual.Grow(-n))
	return nil
}

// handleUseItemOnEntityTransaction ...
func (h *InventoryTransactionHandler) handleUseItemOnEntityTransaction(data *protocol.UseItemOnEntityTransactionData, s *Session, tx *world.Tx, c Controllable) error {
	s.swingingArm.Store(true)
	defer s.swingingArm.Store(false)

	if data.TargetEntityRuntimeID == selfEntityRuntimeID {
		return fmt.Errorf("invalid entity interaction: players cannot interact with themselves")
	}

	handle, ok := s.entityFromRuntimeID(data.TargetEntityRuntimeID)
	if !ok {
		// In some cases, for example when a falling block entity solidifies, latency may allow attacking an entity that
		// no longer exists server side. This is expected, so we shouldn't kick the player.
		s.conf.Log.Debug("invalid entity interaction: no entity with runtime ID", "ID", data.TargetEntityRuntimeID)
		return nil
	}
	e, ok := handle.Entity(tx)
	if !ok {
		s.conf.Log.Debug("invalid entity interaction: entity is not in the same world (anymore)", "ID", data.TargetEntityRuntimeID)
		return nil
	}
	var valid bool
	switch data.ActionType {
	case protocol.UseItemOnEntityActionInteract:
		valid = c.UseItemOnEntity(e)
	case protocol.UseItemOnEntityActionAttack:
		valid = c.AttackEntity(e)
	default:
		return fmt.Errorf("unhandled UseItemOnEntity ActionType %v", data.ActionType)
	}
	if !valid {
		slot := int(*s.heldSlot)
		it, _ := s.inv.Item(slot)
		s.sendItem(it, slot, protocol.WindowIDInventory)
	}
	return nil
}

// handleUseItemTransaction ...
func (h *InventoryTransactionHandler) handleUseItemTransaction(data *protocol.UseItemTransactionData, s *Session, c Controllable) error {
	pos := cube.Pos{int(data.BlockPosition[0]), int(data.BlockPosition[1]), int(data.BlockPosition[2])}
	s.swingingArm.Store(true)
	defer s.swingingArm.Store(false)

	// We reset the inventory so that we can send the held item update without the client already
	// having done that client-side.
	// Because of the new inventory system, the client will expect a transaction confirmation, but instead of doing that
	// it's much easier to just resend the inventory.
	h.resendInventories(s)

	switch data.ActionType {
	case protocol.UseItemActionBreakBlock:
		c.BreakBlock(pos)
	case protocol.UseItemActionClickBlock:
		c.UseItemOnBlock(pos, cube.Face(data.BlockFace), vec32To64(data.ClickedPosition))
	case protocol.UseItemActionClickAir:
		c.UseItem()
	default:
		return fmt.Errorf("unhandled UseItem ActionType %v", data.ActionType)
	}
	return nil
}

// handleReleaseItemTransaction ...
func (h *InventoryTransactionHandler) handleReleaseItemTransaction(c Controllable) error {
	c.ReleaseItem()
	return nil
}
