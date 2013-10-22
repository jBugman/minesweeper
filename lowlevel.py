#! /usr/bin/env python
# -*- coding: UTF-8 -*-

from Quartz import *
from AppKit import NSImage, NSZeroSize
from PIL import Image
import os, os.path
import time

from errors import WindowException, NotImplementedException
from point import Point

class LowLevelApi:
    def __init__(self, title):
        windowList = CGWindowListCopyWindowInfo(kCGWindowListOptionAll, kCGNullWindowID)
        self.window = None
        for window in windowList:
            if kCGWindowName in window and window[kCGWindowName] == title:
                self.window = window
                break
        if self.window:
            self.windowId = window[kCGWindowNumber]
            bounds = self.window[kCGWindowBounds]
            self.size = Point(int(bounds['Width']), int(bounds['Height']))
            self.offset = Point(int(bounds['X']), int(bounds['Y']))
        else:
            raise WindowException('Can not find window')

    def loadAssets(self):
        assets = {}
        for filename in os.listdir('assets'):
            if filename.endswith('.png'):
                image = Image.open(os.path.join('assets', filename))
                assets[os.path.splitext(filename)[0]] = image
        return assets

    def getSnapshot(self):
        windowImage = CGWindowListCreateImageFromArray(CGRectNull, [self.windowId], kCGWindowImageBoundsIgnoreFraming)
        nsImage = NSImage.alloc().initWithCGImage_size_(windowImage, NSZeroSize)
        tiffRepresentation = nsImage.TIFFRepresentation()
        image = Image.fromstring('RGBA', self.size.asTuple(), tiffRepresentation)
        image.load()
        b1 = (0, 1, 2, self.size.y)
        r1 = image.crop(b1)
        r1.load()
        b2 = (2, 0, self.size.x, self.size.y)
        r2 = image.crop(b2)
        r2.load()
        image.paste(r1, (self.size.x - 2, 0, self.size.x, self.size.y - 1))
        image.paste(r2, (0, 0, self.size.x - 2, self.size.y))
        return image

    def isImagesEqual(self, imageA, imageB):
        bytesA = imageA.load()
        bytesB = imageB.load()
        for x in range(imageA.size[0]):
            for y in range(imageA.size[1]):
                if not bytesA[x, y] == bytesB[x, y]:
                    return False
        return True

    def activateWindow(self):
        app = NSRunningApplication.runningApplicationWithProcessIdentifier_(self.window[kCGWindowOwnerPID])
        app.activateWithOptions_(NSApplicationActivateIgnoringOtherApps)
        time.sleep(0.05)
