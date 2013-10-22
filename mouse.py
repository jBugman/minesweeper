#! /usr/bin/env python
# -*- coding: UTF-8 -*-

from Quartz import *
import time

class Mouse:
    @staticmethod
    def press(point, button = kCGEventLeftMouseDown):
        event = CGEventCreateMouseEvent(None, button, point.asTuple(), 0)
        CGEventPost(kCGHIDEventTap, event)
    
    @staticmethod
    def release(point, button = kCGEventLeftMouseUp):
        event = CGEventCreateMouseEvent(None, button, point.asTuple(), 0)
        CGEventPost(kCGHIDEventTap, event)
    
    @staticmethod
    def move(point):
        move = CGEventCreateMouseEvent(None, kCGEventMouseMoved, point.asTuple(), 0)
        CGEventPost(kCGHIDEventTap, move)
    
    @staticmethod
    def position():
        loc = NSEvent.mouseLocation()
        return Point(loc.x, CGDisplayPixelsHigh(0) - loc.y)
    
    @staticmethod
    def click(point):
        Mouse.press(point)
        Mouse.release(point)
        time.sleep(0.025)
    
    @staticmethod
    def rightclick(point):
        Mouse.press(point, kCGEventRightMouseDown)
        Mouse.release(point, kCGEventRightMouseUp)
        time.sleep(0.025)
