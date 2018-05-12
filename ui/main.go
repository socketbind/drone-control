package ui

import (
	"github.com/hajimehoshi/ebiten"
	"log"
	"image"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/socketbind/drone-control/drone"
	"math"
)

const (
	screenWidth  = 960
	screenHeight = 720

	takeOffButton = 1
	flipForwardButton = 12
	flipBackwardButton = 13
	flipLeftButton = 14
	flipRightButton = 15
)

func Start(videoChannel chan *image.Image, commandChannel chan interface{}) {
	var lastImage *ebiten.Image = nil
	var tookOff = false

	update := func (screen *ebiten.Image) error {
		for _, id := range inpututil.JustConnectedGamepadIDs() {
			log.Printf("gamepad connected: id: %d", id)
		}

		ids := ebiten.GamepadIDs()
		if len(ids) > 0 {
			id := ids[0]

			axis0 := ebiten.GamepadAxis(id, 0)
			axis1 := ebiten.GamepadAxis(id, 1)

			if math.Abs(axis0) > 0.05 {
				value := int(axis0 * 30)

				if value < 0 {
					commandChannel <- drone.RotateCounterClockwiseCommand{-value}
				} else if value > 0 {
					commandChannel <- drone.RotateClockwiseCommand{value}
				}
			}

			if math.Abs(axis1) > 0.05 {
				value := int(axis1 * 30)

				if value < 0 {
					commandChannel <- drone.UpCommand{-value}
				} else if value > 0 {
					commandChannel <- drone.DownCommand{value}
				}
			}

			axis2 := ebiten.GamepadAxis(id, 2)
			axis3 := ebiten.GamepadAxis(id, 3)

			if math.Abs(axis2) > 0.2 {
				value := int(axis2 * 30)

				if value < 0 {
					commandChannel <- drone.LeftCommand{-value}
				} else if value > 0 {
					commandChannel <- drone.RightCommand{value}
				}
			}

			if math.Abs(axis3) > 0.2 {
				value := int(axis3 * 30)

				if value < 0 {
					commandChannel <- drone.ForwardCommand{-value}
				} else if value > 0 {
					commandChannel <- drone.BackwardCommand{value}
				}
			}

			maxButton := ebiten.GamepadButton(ebiten.GamepadButtonNum(id))
			for b := ebiten.GamepadButton(id); b < maxButton; b++ {
				if inpututil.IsGamepadButtonJustPressed(id, b) {
					log.Printf("button pressed: id: %d, button: %d", id, b)
				}
				if inpututil.IsGamepadButtonJustReleased(id, b) {
					log.Printf("button released: id: %d, button: %d", id, b)

					switch b {

					case takeOffButton:
						if tookOff {
							commandChannel <- drone.LandCommand{}
						} else {
							commandChannel <- drone.TakeOffCommand{}
						}
						tookOff = !tookOff

					case flipForwardButton:
						commandChannel <- drone.FlipForwardCommand{}

					case flipBackwardButton:
						commandChannel <- drone.FlipBackwardCommand{}

					case flipLeftButton:
						commandChannel <- drone.FlipLeftCommand{}

					case flipRightButton:
						commandChannel <- drone.FlipRightCommand{}
					}
				}
			}
		}

		if ebiten.IsRunningSlowly() {
			return nil
		}

		select {
		case videoImage := <-videoChannel:
			var err error
			lastImage, err = ebiten.NewImageFromImage(*videoImage, ebiten.FilterDefault)
			if err != nil {
				panic("Unable to create image")
			}
		default:
		}

		if lastImage != nil {
			screen.DrawImage(lastImage, nil)
		}

		return nil
	}

	if err := ebiten.Run(update, screenWidth, screenHeight, 1, "Drone Control"); err != nil {
		log.Fatal(err)
	}
}
