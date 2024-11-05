package views

import (
    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/canvas"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/widget"
    "estacionamiento-simulador/models"
    "time"
    "fmt"
    "image/color"
)

type ParkingView struct {
    window      fyne.Window
    rects       []*canvas.Rectangle 
    labels      []*widget.Label     
    simulator   *models.ParkingLot
}

func NewParkingView(simulator *models.ParkingLot) *ParkingView {
    a := app.New()
    w := a.NewWindow("Simulador de Estacionamiento")
    rects := make([]*canvas.Rectangle, simulator.Capacity)
    labels := make([]*widget.Label, simulator.Capacity)
    grid := container.NewGridWithColumns(5)

    for i := 0; i < simulator.Capacity; i++ {
        rects[i] = canvas.NewRectangle(color.RGBA{200, 200, 200, 255})
        rects[i].SetMinSize(fyne.NewSize(60, 60))
        labels[i] = widget.NewLabel("Libre")
        slot := container.NewMax(rects[i], labels[i])
        grid.Add(slot)
    }

    w.SetContent(grid)
    w.Resize(fyne.NewSize(400, 300))

    pv := &ParkingView{window: w, rects: rects, labels: labels, simulator: simulator}
    simulator.RegisterObserver(pv)
    return pv
}

func (pv *ParkingView) Update(occupiedSlots []models.VehicleSlot) {
    for i := range pv.rects {
        if i < len(occupiedSlots) {
            slot := occupiedSlots[i]
            pv.labels[i].SetText(fmt.Sprintf("Vehículo %d (%d s)", slot.VehicleID, slot.ParkTime))

            // ! Cambiamos el color del espacio ocupado a verde
            pv.rects[i].FillColor = color.RGBA{0, 255, 0, 255} 
            pv.rects[i].Refresh()

            go func(label *widget.Label, rect *canvas.Rectangle, parkTime int) {
                for t := parkTime; t > 0; t-- {
                    time.Sleep(1 * time.Second)
                    label.SetText(fmt.Sprintf("Vehículo %d (%d s)", slot.VehicleID, t))
                }
                label.SetText("Libre")
                
                // ! Regresar a color gris cuando el vehículo se va
                rect.FillColor = color.RGBA{200, 200, 200, 255} 
                rect.Refresh()
            }(pv.labels[i], pv.rects[i], slot.ParkTime)
        } else {
            pv.labels[i].SetText("Libre")

            // ! Cambiamos el color del espacio a gris para libre
            pv.rects[i].FillColor = color.RGBA{200, 200, 200, 255} 
            pv.rects[i].Refresh()
        }
    }
    pv.window.Content().Refresh()
}

func StartInterface(simulator *models.ParkingLot) {
    view := NewParkingView(simulator)
    view.window.ShowAndRun()
}
