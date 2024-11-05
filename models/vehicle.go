package models

// Vehicle representa la estructura de un vehículo con su ID y tiempo de estacionamiento.
type Vehicle struct {
    ID       int
    ParkTime int // Tiempo en segundos que el vehículo estará estacionado
}

// VehicleSlot representa la información de un vehículo estacionado.
type VehicleSlot struct {
    VehicleID int
    ParkTime  int // Tiempo restante de estacionamiento en segundos
}
