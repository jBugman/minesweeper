package main

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework CoreGraphics -framework Foundation
#import <Foundation/NSObjCRuntime.h>
#import <CoreGraphics/CoreGraphics.h>
#import <CoreFoundation/CoreFoundation.h>

CFArrayRef makeCFArrayRef(CGWindowID windowID) {
	CGWindowID ids[] = {windowID};
	return CFArrayCreate(NULL, (void *)ids, 1, NULL);
}

*/
import "C"
import (
	"errors"
	"image"
	"image/png"
	"log"
	"os"
	"time"
	"unsafe"
)

func takeScreenShot(windowID int) image.Image {
	const flags = C.kCGWindowImageDefault | C.kCGWindowImageShouldBeOpaque | C.kCGWindowImageNominalResolution
	windowList := C.makeCFArrayRef(C.CGWindowID(windowID))
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

type CFString struct {
	ptr *C.struct___CFString
}

func WrapCFString(pointer unsafe.Pointer) CFString {
	return CFString{(*C.struct___CFString)(pointer)}
}

func (str CFString) String() string {
	cstring := (*C.char)(C.CFStringGetCStringPtr(str.ptr, C.kCFStringEncodingUTF8))
	return C.GoString(cstring)
}

type CFDictionary struct {
	ptr *C.struct___CFDictionary
}

func WrapCFDictionary(pointer unsafe.Pointer) CFDictionary {
	return CFDictionary{(*C.struct___CFDictionary)(pointer)}
}

func (dict CFDictionary) ObjectForKey(key C.CFStringRef) unsafe.Pointer {
	ptr := C.CFDictionaryGetValue(dict.ptr, unsafe.Pointer(key))
	return ptr
}

func (dict CFDictionary) StringForKey(key C.CFStringRef) string {
	ptr := dict.ObjectForKey(key)
	return WrapCFString(ptr).String()
}

func (dict CFDictionary) IntForKey(key C.CFStringRef) int {
	ptr := dict.ObjectForKey(key)
	var number int
	C.CFNumberGetValue((*C.struct___CFNumber)(ptr), C.kCFNumberIntType, unsafe.Pointer(&number))
	return number
}

func (dict CFDictionary) CGRectForKey(key C.CFStringRef) C.CGRect {
	dictRepresentation := WrapCFDictionary(dict.ObjectForKey(key))
	var rect C.CGRect
	C.CGRectMakeWithDictionaryRepresentation(dictRepresentation.ptr, &rect)
	return rect
}

type windowMeta struct {
	id     int
	title  string
	bounds C.CGRect
}

func findWindow(targetAppTitle string) (windowMeta, error) {
	windows := C.CGWindowListCopyWindowInfo(C.kCGWindowListOptionOnScreenOnly, C.kCGNullWindowID)
	var i C.CFIndex
	var info windowMeta
	for i = 0; i < C.CFArrayGetCount(windows); i++ {
		window := WrapCFDictionary(C.CFArrayGetValueAtIndex(windows, i))
		owner := window.StringForKey(C.kCGWindowOwnerName)
		if owner == targetAppTitle {
			info.title = window.StringForKey(C.kCGWindowName)
			info.id = window.IntForKey(C.kCGWindowNumber)
			info.bounds = window.CGRectForKey(C.kCGWindowBounds)
			return info, nil
		}
	}
	return info, errors.New("Window not found")
}

func saveImage(img image.Image) {
	outFile, _ := os.Create("test.png")
	png.Encode(outFile, img)
}

func main() {
	const targetAppTitle = "Minesweeper"
	t := time.Now()
	winMeta, err := findWindow(targetAppTitle)
	if err != nil {
		log.Fatal(err)
	}
	log.Println(winMeta)
	screenShot := takeScreenShot(winMeta.id)
	log.Println(time.Since(t))

	saveImage(screenShot)
}
