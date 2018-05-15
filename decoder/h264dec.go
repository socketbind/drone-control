package decoder

// Contains some parts from https://github.com/gqf2008/codec.
// Works with latest libavcodec.

import (
	/*
	#cgo LDFLAGS: -lavformat -lavutil -lavcodec
	#include "h264_decode.h"
	*/
	"C"
	"unsafe"
	"errors"
	"image"
)

var m C.h264dec_t

type frameHandlerFunc func(*image.Image)

var frameHandler frameHandlerFunc

func Init() (err error) {
	r := C.h264dec_new(&m)
	if int(r) < 0 {
		err = errors.New("open codec failed")
	}
	return
}

func Free() {
	C.h264dec_free(&m)
}

//export handleFrame
func handleFrame(f *C.AVFrame) {
	if frameHandler != nil {
		w := int(f.width)
		h := int(f.height)
		ys := int(f.linesize[0])
		cs := int(f.linesize[1])

		raw := &image.YCbCr{
			Y: fromCPtr(unsafe.Pointer(f.data[0]), ys*h),
			Cb: fromCPtr(unsafe.Pointer(f.data[1]), cs*h/2),
			Cr: fromCPtr(unsafe.Pointer(f.data[2]), cs*h/2),
			YStride: ys,
			CStride: cs,
			SubsampleRatio: image.YCbCrSubsampleRatio420,
			Rect: image.Rect(0, 0, w, h),
		}

		im := raw.SubImage(raw.Bounds())

		frameHandler(&im)
	}
}

func Decode(nal []byte, handlerFn frameHandlerFunc) (err error) {
	frameHandler = handlerFn

	r := C.h264dec_decode(
		&m,
		(*C.uint8_t)(unsafe.Pointer(&nal[0])),
		(C.int)(len(nal)),
	)

	if int(r) < 0 {
		err = errors.New("decode failed")
		return
	}

	return
}

