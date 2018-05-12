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

	deadZoneHorizontal = 0.5
	deadZoneVertical = 0.5
)

func remapAxisInput(inputValue float64, deadZone float64, maxValue float64) int {
	if math.Abs(inputValue) > deadZone {
		if inputValue < 0 {
			return int(((inputValue + deadZone) / (1.0 - deadZone)) * maxValue)
		} else if inputValue > 0 {
			return int(((inputValue - deadZone) / (1.0 - deadZone)) * maxValue)
		}
	}

	return 0
}

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

			rotation := remapAxisInput(axis0, deadZoneHorizontal, 30)
			if rotation != 0 {
				if rotation < 0 {
					commandChannel <- drone.RotateCounterClockwiseCommand{-rotation}
				} else if rotation > 0 {
					commandChannel <- drone.RotateClockwiseCommand{rotation}
				}
			}

			altitude := remapAxisInput(axis1, deadZoneVertical, 30)
			if altitude != 0 {
				if altitude < 0 {
					commandChannel <- drone.UpCommand{-altitude}
				} else if altitude > 0 {
					commandChannel <- drone.DownCommand{altitude}
				}
			}

			axis2 := ebiten.GamepadAxis(id, 2)
			axis3 := ebiten.GamepadAxis(id, 3)

			horizontalMovement := remapAxisInput(axis2, deadZoneHorizontal, 20)
			if horizontalMovement != 0 {
				if horizontalMovement < 0 {
					commandChannel <- drone.LeftCommand{-horizontalMovement}
				} else if horizontalMovement > 0 {
					commandChannel <- drone.RightCommand{horizontalMovement}
				}
			}

			verticalMovement := remapAxisInput(axis3, deadZoneVertical, 20)
			if verticalMovement != 0 {
				if verticalMovement < 0 {
					commandChannel <- drone.ForwardCommand{-verticalMovement}
				} else if verticalMovement > 0 {
					commandChannel <- drone.BackwardCommand{verticalMovement}
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
