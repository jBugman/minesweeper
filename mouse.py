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
		time.sleep(0.01)
		Mouse.press(point)
		time.sleep(0.01)
		Mouse.release(point)
		time.sleep(0.01)
	
	@staticmethod
	def rightclick(point):
		time.sleep(0.01)
		Mouse.press(point, kCGEventRightMouseDown)
		time.sleep(0.01)
		Mouse.release(point, kCGEventRightMousUp)
		time.sleep(0.01)