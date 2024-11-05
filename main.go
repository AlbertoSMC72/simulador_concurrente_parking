package main

import (
    "fmt"
    "estacionamiento-simulador/scenes"
    "estacionamiento-simulador/views"
)

func main() {
    simulator := scenes.NewSimulator(20)
    fmt.Println("Iniciando la interfaz de simulación...")
    go simulator.StartSimulation(100)
    views.StartInterface(simulator.ParkingLot)
}
