package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework CoreGraphics -framework Foundation
#import <Foundation/NSObjCRuntime.h>
#import <CoreGraphics/CoreGraphics.h>
*/
import "C"
import (
	"fmt"
	"image"
	"image/png"
	"os"
	"time"
	"unsafe"
)

func takeScreenShot() image.Image {
	const flags = C.kCGWindowImageDefault | C.kCGWindowImageShouldBeOpaque | C.kCGWindowImageNominalResolution
	screenShot := C.CGWindowListCreateImage(C.CGRectInfinite, C.kCGWindowListOptionOnScreenOnly, C.kCGNullWindowID, flags)

	rawBytes := C.CGDataProviderCopyData(C.CGImageGetDataProvider(screenShot))
	pointer := C.CFDataGetBytePtr(rawBytes)
	length := int(C.CFDataGetLength(rawBytes))
	pixels := C.GoBytes(unsafe.Pointer(pointer), C.int(length))

	width := int(C.CGImageGetWidth(screenShot))
	height := int(C.CGImageGetHeight(screenShot))
	bytesPerRow := int(C.CGImageGetBytesPerRow(screenShot))

	fmt.Println(C.CGImageGetBitmapInfo(screenShot))

	// manually fix BGRA -> RGBA
	for i := 0; i < length; i += 4 {
		pixels[i], pixels[i+2] = pixels[i+2], pixels[i]
	}

	return &image.RGBA{Pix: pixels, Stride: bytesPerRow, Rect: image.Rect(0, 0, width, height)}
}

func main() {
	t := time.Now()
	screenShot := takeScreenShot()
	fmt.Println(time.Since(t))

	outFile, _ := os.Create("test.png")
	png.Encode(outFile, screenShot)
}
