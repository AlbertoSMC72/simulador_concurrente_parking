package models

type Observer interface {
    Update(occupiedSlots []VehicleSlot)
}