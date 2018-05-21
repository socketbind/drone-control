package ui

import (
	"image"
	"github.com/veandco/go-sdl2/sdl"
	"log"
	"github.com/socketbind/drone-control/util"
	"github.com/socketbind/drone-control/drone"
)

const (
	screenWidth  = 960
	screenHeight = 720

	takeOffButton = sdl.CONTROLLER_BUTTON_A
	flipForwardButton = sdl.CONTROLLER_BUTTON_DPAD_UP
	flipBackwardButton = sdl.CONTROLLER_BUTTON_DPAD_DOWN
	flipLeftButton = sdl.CONTROLLER_BUTTON_DPAD_LEFT
	flipRightButton = sdl.CONTROLLER_BUTTON_DPAD_RIGHT

	deadZoneHorizontal = 16000
	deadZoneVertical = 16000
)

func remapAxisInput(inputValue int16, deadZone int16, maxOutputValue int16) int {
	if util.Abs(inputValue) > deadZone {
		if inputValue < 0 {
			return int((float32(inputValue + deadZone) / float32(32767 - deadZone)) * float32(maxOutputValue))
		} else if inputValue > 0 {
			return int((float32(inputValue - deadZone) / float32(32767 - deadZone)) * float32(maxOutputValue))
		}
	}

	return 0
}

func openController(index int) *sdl.GameController {
	if sdl.NumJoysticks() >= index + 1 {
		if sdl.IsGameController(index) {
			ctrl := sdl.GameControllerOpen(index)

			name := sdl.GameControllerNameForIndex(index)
			log.Println("Got controller: ", name)

			sdl.GameControllerEventState(sdl.ENABLE)

			return ctrl
		} else {
			log.Println("Joystick 0 is not a game controller somehow?")
		}
	} else {
		log.Println("No joysticks with index ", index)
	}

	return nil
}

func Start(videoChannel chan *image.YCbCr, commandChannel chan interface{}) {
	if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
		panic(err)
	}
	defer sdl.Quit()

	window, err := sdl.CreateWindow("Drone Control", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		screenWidth, screenHeight, sdl.WINDOW_SHOWN)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	util.GameControllerAddMappingsFromFile("gamecontrollerdb.txt")

	gamepad := openController(0)
	if gamepad != nil {
		defer gamepad.Close()
	}

	mainLoop(videoChannel, commandChannel, renderer)
}

func handleControllerAxisEvent(event *sdl.ControllerAxisEvent, commandChannel chan interface{}) {
	if event.Axis == sdl.CONTROLLER_AXIS_LEFTX {
		rotation := remapAxisInput(event.Value, deadZoneHorizontal, 30)
		if rotation != 0 {
			if rotation < 0 {
				commandChannel <- drone.RotateCounterClockwiseCommand{-rotation}
			} else if rotation > 0 {
				commandChannel <- drone.RotateClockwiseCommand{rotation}
			}
		}
	} else if event.Axis == sdl.CONTROLLER_AXIS_LEFTY {
		altitude := remapAxisInput(event.Value, deadZoneVertical, 30)
		if altitude != 0 {
			if altitude < 0 {
				commandChannel <- drone.UpCommand{-altitude}
			} else if altitude > 0 {
				commandChannel <- drone.DownCommand{altitude}
			}
		}
	}
}

var tookOff = false
func handleControllerButtonEvent(event *sdl.ControllerButtonEvent, commandChannel chan interface{}) {
	if event.State == sdl.PRESSED {
		switch event.Button {
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

func mainLoop(videoChannel chan *image.YCbCr, commandChannel chan interface{}, renderer *sdl.Renderer) {
	var err error
	var videoTexture *sdl.Texture = nil
	var videoBounds *sdl.Rect = nil

	screenRect := &sdl.Rect{0, 0, screenWidth, screenHeight}

	running := true
	for running {
		for polledEvent := sdl.PollEvent(); polledEvent != nil; polledEvent = sdl.PollEvent() {
			switch event := polledEvent.(type) {
			case *sdl.QuitEvent:
				running = false
			case *sdl.ControllerAxisEvent:
				handleControllerAxisEvent(event, commandChannel)
			case *sdl.ControllerButtonEvent:
				handleControllerButtonEvent(event, commandChannel)
			}
		}

		select {
		case videoImage := <-videoChannel:
			if videoTexture == nil {
				videoTexture, err = renderer.CreateTexture(sdl.PIXELFORMAT_IYUV, sdl.TEXTUREACCESS_STREAMING, 960, 720)
				if err != nil {
					panic(err)
				}

				videoBounds = &sdl.Rect{
					0, 0,
					int32(videoImage.Bounds().Max.X), int32(videoImage.Bounds().Max.Y),
				}
			}

			videoTexture.UpdateYUV(
				videoBounds,
				videoImage.Y,
				videoImage.YStride,
				videoImage.Cb,
				videoImage.CStride,
				videoImage.Cr,
				videoImage.CStride)
		default:
		}

		renderer.Clear()
		renderer.SetDrawColor(255, 0, 0, 255)
		renderer.FillRect(screenRect)

		if videoTexture != nil && videoBounds != nil {
			renderer.Copy(videoTexture, videoBounds, screenRect)
		}

		renderer.Present()

		sdl.Delay(1000 / 30)
	}

	if videoTexture != nil {
		videoTexture.Destroy()
	}
}