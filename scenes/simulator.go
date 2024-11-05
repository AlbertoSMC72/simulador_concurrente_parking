package scenes

import (
    "fmt"
    "math/rand"
    "time"
    "estacionamiento-simulador/models"
)

type Simulator struct {
    ParkingLot *models.ParkingLot
}

func NewSimulator(capacity int) *Simulator {
    return &Simulator{
        ParkingLot: models.NewParkingLot(capacity),
    }
}

func (s *Simulator) StartSimulation(vehicleCount int) {
    fmt.Println("Iniciando simulación de vehículos...")
    for i := 0; i < vehicleCount; i++ {
        fmt.Printf("Generando vehículo %d...\n", i)
		time.Sleep(1 * time.Second)
        go s.SimulateVehicle(i)
    }
    fmt.Println("Finalización de generación de vehículos.")
}

func (s *Simulator) SimulateVehicle(id int) {
    parkTime := rand.Intn(6) + 10 // ! Tiempo estacionado 
    fmt.Printf("Vehículo %d intenta entrar\n", id)

    for {
        select {
        case s.ParkingLot.EntryExit <- struct{}{}: 
            if s.ParkingLot.ParkVehicle(id, parkTime) {
                fmt.Printf("Vehículo %d estacionado por %d segundos\n", id, parkTime)
                
				<-s.ParkingLot.EntryExit
                
                time.Sleep(time.Duration(parkTime) * time.Second)
                
                s.ParkingLot.LeaveVehicle(id)
                fmt.Printf("Vehículo %d salió del estacionamiento\n", id)
                return 
            } else {
                <-s.ParkingLot.EntryExit
                fmt.Printf("Vehículo %d esperando espacio en el estacionamiento\n", id)
                <-s.ParkingLot.NotifyChannel // ! Espera a ser notificado de un espacio libre
            }
        default:
            fmt.Printf("Vehículo %d esperando que la puerta esté libre\n", id)
            <-s.ParkingLot.NotifyChannel // ! Espera a ser notificado de un espacio libre
        }
    }
}
