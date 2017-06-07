package macos

/*
#cgo CFLAGS: -x objective-c
#cgo LDFLAGS: -framework CoreGraphics  -framework Foundation
#import <CoreGraphics/CoreGraphics.h>
#import <CoreFoundation/CoreFoundation.h>

// Need to do it on a C side as you can't pass ** from Cgo
CFArrayRef makeSingletonIntCFArrayRef(int i) {
	int c_array[] = {i};
	return CFArrayCreate(NULL, (void *)c_array, 1, NULL);
}

*/
import "C"
import (
	"unsafe"

	"./keycode"
)

// WrapIntAsCFArrayRef returns new CFArrayRef with a provided integer as a single element
func WrapIntAsCFArrayRef(i int) C.CFArrayRef {
	return C.makeSingletonIntCFArrayRef(C.int(i))
}

type coreFoundation struct{}

// CoreFoundation is a namespace wrapper for Core Foundation
var CoreFoundation coreFoundation

type coreGraphics struct {
	WindowName      C.CFStringRef
	WindowOwnerName C.CFStringRef
}

// CoreGraphics is a namespace wrapper for Core Graphics
var CoreGraphics = coreGraphics{
	WindowName:      C.kCGWindowName,
	WindowOwnerName: C.kCGWindowOwnerName,
}

// CFString is a basic wrapper for Core foundation CFString
type CFString struct {
	ptr *C.struct___CFString
}

// WrapCFString converts unsafe.Pointer to CFString
func (cf coreFoundation) WrapCFString(pointer unsafe.Pointer) CFString {
	return CFString{(*C.struct___CFString)(pointer)}
}

func (str CFString) String() string { // TODO can something go wrong?
	cstring := (*C.char)(C.CFStringGetCStringPtr(str.ptr, C.kCFStringEncodingUTF8))
	return C.GoString(cstring)
}

// CFDictionary is a basic wrapper for Core foundation CFDictionary
type CFDictionary struct {
	ptr *C.struct___CFDictionary
}

// WrapCFDictionary converts unsafe.Pointer to CFDictionary
func (cf coreFoundation) WrapCFDictionary(pointer unsafe.Pointer) CFDictionary { // TODO handle errors
	return CFDictionary{(*C.struct___CFDictionary)(pointer)}
}

// ObjectForKey returns generic dictionary value for key
func (dict CFDictionary) ObjectForKey(key C.CFStringRef) unsafe.Pointer { // TODO handle errors
	ptr := C.CFDictionaryGetValue(dict.ptr, unsafe.Pointer(key))
	return ptr
}

// StringForKey returns dictionary value as string
func (dict CFDictionary) StringForKey(key C.CFStringRef) string { // TODO handle errors
	ptr := dict.ObjectForKey(key)
	return CoreFoundation.WrapCFString(ptr).String()
}

// IntForKey returns dictionary value as int
func (dict CFDictionary) IntForKey(key C.CFStringRef) int { // TODO handle errors
	ptr := dict.ObjectForKey(key)
	var number int
	C.CFNumberGetValue((*C.struct___CFNumber)(ptr), C.kCFNumberIntType, unsafe.Pointer(&number))
	return number
}

// CGRectForKey returns dictionary value as C.CGRect
func (dict CFDictionary) CGRectForKey(key C.CFStringRef) C.CGRect { // TODO handle errors
	dictRepresentation := CoreFoundation.WrapCFDictionary(dict.ObjectForKey(key))
	var rect C.CGRect
	C.CGRectMakeWithDictionaryRepresentation(dictRepresentation.ptr, &rect)
	return rect
}

// CreateMouseEvent is a binding for CGEventCreateMouseEvent
func (cg coreGraphics) CreateMouseEvent(mouseType C.CGEventType, position C.CGPoint, button C.CGMouseButton) C.CGEventRef {
	return C.CGEventCreateMouseEvent(nil, mouseType, position, button)
}

// CreateKeyboardEvent is a binding for CGEventCreateKeyboardEvent
func (cg coreGraphics) CreateKeyboardEvent(keyCode keycode.Code, isKeyDown bool) C.CGEventRef {
	return C.CGEventCreateKeyboardEvent(nil, C.CGKeyCode(keyCode), C._Bool(isKeyDown))
}

func releaseEvent(event C.CGEventRef) {
	C.CFRelease(C.CFTypeRef(event))
}

// CGRect wraps C struct
type CGRect struct {
	rect C.CGRect
}

// X returns origin x of a CGRect
func (r CGRect) X() int {
	return int(r.rect.origin.x)
}

// Y returns origin y of a CGRect
func (r CGRect) Y() int {
	return int(r.rect.origin.y)
}

// Width of a CGRect
func (r CGRect) Width() uint {
	return uint(r.rect.size.width)
}

// Height of a CGRect
func (r CGRect) Height() uint {
	return uint(r.rect.size.height)
}
