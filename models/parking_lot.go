package models

import (
    "sync"
)

type ParkingLot struct {
    Capacity      int
    Occupied      int
    EntryExit     chan struct{}      // Control de acceso a la puerta
    NotifyChannel chan struct{}      // Canal de notificación para reintentar entrada
    Observers     []Observer         // Lista de observadores
    OccupiedSlots []VehicleSlot      // Lista de vehículos estacionados
    Mutex         sync.Mutex         // Mutex para evitar condiciones de carrera
}

func NewParkingLot(capacity int) *ParkingLot {
    return &ParkingLot{
        Capacity:      capacity,
        EntryExit:     make(chan struct{}, 1),
        NotifyChannel: make(chan struct{}, capacity),
        Observers:     make([]Observer, 0),
    }
}

func (pl *ParkingLot) RegisterObserver(observer Observer) {
    pl.Observers = append(pl.Observers, observer)
}

func (pl *ParkingLot) NotifyObservers() {
    for _, observer := range pl.Observers {
        observer.Update(pl.OccupiedSlots)
    }
}

func (pl *ParkingLot) ParkVehicle(vehicleID, parkTime int) bool {
    pl.Mutex.Lock()
    defer pl.Mutex.Unlock()

    if pl.Occupied < pl.Capacity {
        pl.Occupied++
        pl.OccupiedSlots = append(pl.OccupiedSlots, VehicleSlot{VehicleID: vehicleID, ParkTime: parkTime})
        pl.NotifyObservers() // Notifica a los observadores sobre el nuevo estado
        return true
    }
    return false
}

func (pl *ParkingLot) LeaveVehicle(vehicleID int) {
    pl.Mutex.Lock()
    defer pl.Mutex.Unlock()

    for i, slot := range pl.OccupiedSlots {
        if slot.VehicleID == vehicleID {
            // Remueve el vehículo de los espacios ocupados
            pl.OccupiedSlots = append(pl.OccupiedSlots[:i], pl.OccupiedSlots[i+1:]...)
            pl.Occupied--
            break
        }
    }
    pl.NotifyObservers() // Notifica a los observadores sobre el nuevo estado
    // Envía señal al canal de notificación indicando que un espacio ha sido liberado
    pl.NotifyChannel <- struct{}{}
}
