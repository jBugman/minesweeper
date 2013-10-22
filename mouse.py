# -*- coding: UTF-8 -*-
import time

import Quartz


def _send_tap_event(event):
    Quartz.CGEventPost(Quartz.kCGHIDEventTap, event)


class Mouse:
    @staticmethod
    def press(point, button=Quartz.kCGEventLeftMouseDown):
        event = Quartz.CGEventCreateMouseEvent(None, button, point.as_tuple(), 0)
        _send_tap_event(event)

    @staticmethod
    def release(point, button=Quartz.kCGEventLeftMouseUp):
        event = Quartz.CGEventCreateMouseEvent(None, button, point.as_tuple(), 0)
        _send_tap_event(event)

    @staticmethod
    def move(point):
        move = CGEventCreateMouseEvent(None, Quartz.kCGEventMouseMoved, point.as_tuple(), 0)
        _send_tap_event(event)

    @staticmethod
    def position():
        loc = Quartz.NSEvent.mouseLocation()
        return Point(loc.x, Quartz.CGDisplayPixelsHigh(0) - loc.y)

    @staticmethod
    def click(point):
        Mouse.press(point)
        Mouse.release(point)
        time.sleep(0.025)

    @staticmethod
    def rightclick(point):
        Mouse.press(point, Quartz.kCGEventRightMouseDown)
        Mouse.release(point, Quartz.kCGEventRightMouseUp)
        time.sleep(0.025)
