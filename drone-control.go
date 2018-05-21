package main

import (
	"log"
	"github.com/socketbind/drone-control/drone"
	"image"
	"github.com/socketbind/drone-control/decoder"
	"github.com/socketbind/drone-control/ui"
)

var videoChannel = make(chan *image.YCbCr)
var commandChannel = make(chan interface{})

func main() {
	err := decoder.Init()
	if err != nil {
		log.Fatal("Unable to create decoder")
	}

	defer decoder.Free()

	go drone.DroneControl(videoChannel, commandChannel)

	ui.Start(videoChannel, commandChannel)
}
