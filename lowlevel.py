# -*- coding: UTF-8 -*-
import os
import time

import Quartz
from AppKit import NSImage, NSZeroSize
from PIL import Image

import errors
from point import Point
import log


class LowLevelApi:
    def __init__(self, title):
        self.logger = log.get_logger(__name__)

        self.logger.debug('Finding game window')
        self.window = None
        windows = Quartz.CGWindowListCopyWindowInfo(Quartz.kCGWindowListOptionAll, Quartz.kCGNullWindowID)
        for window in windows:
            if Quartz.kCGWindowName in window and window[Quartz.kCGWindowName] == title:
                self.window = window
                break
        if self.window:
            self.window_id = window[Quartz.kCGWindowNumber]
            bounds = self.window[Quartz.kCGWindowBounds]
            self.size = Point(int(bounds['Width']), int(bounds['Height']))
            self.offset = Point(int(bounds['X']), int(bounds['Y']))
        else:
            raise errors.WindowException

    def load_assets(self):
        assets = {}
        for filename in os.listdir('assets'):
            if filename.endswith('.png'):
                image = Image.open(os.path.join('assets', filename))
                assets[os.path.splitext(filename)[0]] = image
        return assets

    def get_snapshot(self):
        window_image = Quartz.CGWindowListCreateImageFromArray(Quartz.CGRectNull, [self.window_id], Quartz.kCGWindowImageBoundsIgnoreFraming)
        ns_image = Quartz.NSImage.alloc().initWithCGImage_size_(window_image, Quartz.NSZeroSize)
        tiff_representation = ns_image.TIFFRepresentation()
        image = Image.fromstring('RGBA', self.size.as_tuple(), tiff_representation)
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

    def is_images_equal(self, image_a, image_b):
        bytes_a = image_a.load()
        bytes_b = image_b.load()
        for x in range(image_a.size[0]):
            for y in range(image_a.size[1]):
                if not bytes_a[x, y] == bytes_b[x, y]:
                    return False
        return True

    def activate_window(self):
        self.logger.debug('Activating window')
        app = Quartz.NSRunningApplication.runningApplicationWithProcessIdentifier_(self.window[Quartz.kCGWindowOwnerPID])
        app.activateWithOptions_(Quartz.NSApplicationActivateIgnoringOtherApps)
        time.sleep(0.2)
