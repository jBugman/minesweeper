package macos

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework CoreGraphics
#import <CoreGraphics/CoreGraphics.h>
*/
import "C"
import (
	"errors"
	"image"
	"unsafe"
)

// WindowMeta contains some info about Quartz window
type WindowMeta struct {
	ID     int
	Title  string
	Bounds C.CGRect
}

// FindWindow returns if of a first window of target app
func FindWindow(targetAppTitle string) (WindowMeta, error) {
	windows := C.CGWindowListCopyWindowInfo(C.kCGWindowListOptionOnScreenOnly, C.kCGNullWindowID)
	var i C.CFIndex
	var info WindowMeta
	for i = 0; i < C.CFArrayGetCount(windows); i++ {
		window := CoreFoundation.WrapCFDictionary(C.CFArrayGetValueAtIndex(windows, i))
		owner := window.StringForKey(CoreGraphics.WindowOwnerName)
		if owner == targetAppTitle {
			info.Title = window.StringForKey(CoreGraphics.WindowName)
			info.ID = window.IntForKey(C.kCGWindowNumber)        // TODO export constant?
			info.Bounds = window.CGRectForKey(C.kCGWindowBounds) // TODO export constant?
			return info, nil
		}
	}
	return info, errors.New("Window not found")
}

// TakeScreenshot takes screenshot of a window with provided ID
func TakeScreenshot(windowID int) image.Image { // TODO handle errors
	const flags = C.kCGWindowImageDefault | C.kCGWindowImageShouldBeOpaque | C.kCGWindowImageNominalResolution
	windowList := WrapIntAsCFArrayRef(windowID)
	screenShot := C.CGWindowListCreateImageFromArray(C.CGRectNull, windowList, flags)

	rawBytes := C.CGDataProviderCopyData(C.CGImageGetDataProvider(screenShot))
	pointer := C.CFDataGetBytePtr(rawBytes)
	length := int(C.CFDataGetLength(rawBytes))
	pixels := C.GoBytes(unsafe.Pointer(pointer), C.int(length))

	width := int(C.CGImageGetWidth(screenShot))
	height := int(C.CGImageGetHeight(screenShot))
	bytesPerRow := int(C.CGImageGetBytesPerRow(screenShot))

	// manually fix BGRA -> RGBA
	for i := 0; i < length; i += 4 {
		pixels[i], pixels[i+2] = pixels[i+2], pixels[i]
	}

	return &image.RGBA{Pix: pixels, Stride: bytesPerRow, Rect: image.Rect(0, 0, width, height)}
}
