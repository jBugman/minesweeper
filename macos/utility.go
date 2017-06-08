package macos

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework CoreGraphics -framework AppKit
#import <CoreGraphics/CoreGraphics.h>
#import <AppKit/AppKit.h>

// Wrapped Objective-C code for window activation
void activateWindow(int ownerPID) {
	NSRunningApplication* app = [NSRunningApplication runningApplicationWithProcessIdentifier: ownerPID];
	[app activateWithOptions: NSApplicationActivateIgnoringOtherApps];
}

*/
import "C"
import (
	"errors"
	"image"
	"time"
	"unsafe"

	"./keycode"
)

// WindowMeta contains some info about Quartz window
type WindowMeta struct {
	ID       int
	OwnerPID int
	Title    string
	Bounds   CGRect
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
			info.OwnerPID = window.IntForKey(C.kCGWindowOwnerPID)        // TODO export constant?
			info.ID = window.IntForKey(C.kCGWindowNumber)                // TODO export constant?
			info.Bounds = CGRect{window.CGRectForKey(C.kCGWindowBounds)} // TODO export constant?
			return info, nil
		}
	}
	return info, errors.New("Window not found")
}

// TakeScreenshot takes screenshot of a window with provided ID
func TakeScreenshot(windowID int) *image.RGBA { // TODO handle errors
	const flags = C.kCGWindowImageDefault | C.kCGWindowImageShouldBeOpaque | C.kCGWindowImageNominalResolution | C.kCGWindowImageBoundsIgnoreFraming
	windowList := WrapIntAsCFArrayRef(windowID)
	screenShot := C.CGWindowListCreateImageFromArray(C.CGRectNull, windowList, flags)

	rawBytes := C.CGDataProviderCopyData(C.CGImageGetDataProvider(screenShot))
	pointer := C.CFDataGetBytePtr(rawBytes)
	length := int(C.CFDataGetLength(rawBytes))
	pixels := C.GoBytes(unsafe.Pointer(pointer), C.int(length))

	width := int(C.CGImageGetWidth(screenShot))
	height := int(C.CGImageGetHeight(screenShot))
	bytesPerRow := int(C.CGImageGetBytesPerRow(screenShot))

	// Manually fix BGRA -> RGBA
	for i := 0; i < length; i += 4 {
		pixels[i], pixels[i+2] = pixels[i+2], pixels[i]
	}

	return &image.RGBA{Pix: pixels, Stride: bytesPerRow, Rect: image.Rect(0, 0, width, height)}
}

const (
	mouseClickDuration = 75 * time.Millisecond
	keyPressDuration   = 20 * time.Millisecond
)

// LeftClick does that it sounds
func LeftClick(x, y int) {
	genericClick(C.kCGEventLeftMouseDown, C.kCGEventLeftMouseUp, C.kCGMouseButtonLeft, x, y)
}

// RightClick does that it sounds
func RightClick(x, y int) {
	genericClick(C.kCGEventRightMouseDown, C.kCGEventRightMouseUp, C.kCGMouseButtonRight, x, y)
}

func genericClick(downEventType, upEventType C.CGEventType, button C.CGMouseButton, x, y int) {
	point := C.CGPointMake(C.CGFloat(x), C.CGFloat(y))
	downEvent := CoreGraphics.CreateMouseEvent(downEventType, point, button)
	upEvent := CoreGraphics.CreateMouseEvent(upEventType, point, button)
	defer releaseEvent(downEvent)
	defer releaseEvent(upEvent)
	C.CGEventPost(C.kCGHIDEventTap, downEvent)
	time.Sleep(mouseClickDuration)
	C.CGEventPost(C.kCGHIDEventTap, upEvent)
	time.Sleep(mouseClickDuration)
}

// KeyPress emulates keyboard key press
func KeyPress(keyCode keycode.Code) {
	downEvent := CoreGraphics.CreateKeyboardEvent(keyCode, true)
	upEvent := CoreGraphics.CreateKeyboardEvent(keyCode, false)
	defer releaseEvent(downEvent)
	defer releaseEvent(upEvent)
	C.CGEventPost(C.kCGHIDEventTap, downEvent)
	time.Sleep(keyPressDuration)
	C.CGEventPost(C.kCGHIDEventTap, upEvent)
	time.Sleep(keyPressDuration)
}

// KeyPressWithModifier emulates keyboard press with modifier key
func KeyPressWithModifier(keyCode, modifierKeyCode keycode.Code) {
	modifierDownEvent := CoreGraphics.CreateKeyboardEvent(modifierKeyCode, true)
	modifierUpEvent := CoreGraphics.CreateKeyboardEvent(modifierKeyCode, false)
	defer releaseEvent(modifierDownEvent)
	defer releaseEvent(modifierUpEvent)
	C.CGEventPost(C.kCGHIDEventTap, modifierDownEvent)
	time.Sleep(keyPressDuration)
	KeyPress(keyCode)
	C.CGEventPost(C.kCGHIDEventTap, modifierUpEvent)
	time.Sleep(keyPressDuration)
}

// ActivateWindow brings window of a selected app to front
func ActivateWindow(windowOwnerPID int) {
	C.activateWindow(C.int(windowOwnerPID))
	time.Sleep(200 * time.Millisecond)
}
