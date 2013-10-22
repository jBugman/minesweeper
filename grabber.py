# -*- coding: UTF-8 -*-
from lowlevel import LowLevelApi
from point import Point
from consts import TITLE, CELL_SIZE, HEADER_HEIGHT, FOOTER_HEIGHT
import log


class Grabber:
    def __init__(self):
        self.logger = log.get_logger(__name__)
        self.logger.debug('Initialising grabber')

        self.api = LowLevelApi(TITLE)
        self.assets = self.api.load_assets()
        self.field_size = Point(int(self.api.size.x / CELL_SIZE), int((self.api.size.y - HEADER_HEIGHT - FOOTER_HEIGHT) / CELL_SIZE))

    def _get_local_offset(self):
        return Point(0, HEADER_HEIGHT)

    def get_offset(self):
        return self._get_local_offset() + self.api.offset

    def grab_field(self):
        offset = self._get_local_offset()
        snapshot = self.api.get_snapshot()
        cells = [[None for _ in range(self.field_size.x)] for _ in range(self.field_size.y)]
        for x in range(self.field_size.x):
            for y in range(self.field_size.y):
                coords = Point(CELL_SIZE * x, CELL_SIZE * y)
                cell = self.subimage(snapshot, coords)
                if self.compare(cell, 'empty'):
                    cells[y][x] = '?'
                elif self.compare(cell, '0'):
                    cells[y][x] = '0'
                elif self.compare(cell, '1'):
                    cells[y][x] = '1'
                elif self.compare(cell, '2'):
                    cells[y][x] = '2'
                elif self.compare(cell, '3'):
                    cells[y][x] = '3'
                elif self.compare(cell, '4'):
                    cells[y][x] = '4'
                elif self.compare(cell, '5'):
                    cells[y][x] = '5'
                elif self.compare(cell, '6'):
                    cells[y][x] = '6'
                elif self.compare(cell, '7'):
                    cells[y][x] = '7'
                elif self.compare(cell, '8'):
                    cells[y][x] = '8'
                elif self.compare(cell, 'mine') or self.compare(cell, 'bang'):
                    cells[y][x] = 'M'
                elif self.compare(cell, 'flag'):
                    cells[y][x] = '+'
                elif self.compare(cell, 'win'):
                    cells[y][x] = '!'
                else:
                    cells[y][x] = '#'
                    cell.save('test/test{}{}.png'.format(x, y))
        return cells

    def subimage(self, image, point, size=Point(CELL_SIZE, CELL_SIZE)):
        point = point + self._get_local_offset()
        rect = (point.x, point.y, point.x + size.x, point.y + size.y)
        result = image.crop(rect)
        result.load()
        return result

    def compare(self, image, sample):
        is0 = (sample + '@0') in self.assets and self.api.is_images_equal(image, self.assets[sample + '@0'])
        is1 = (sample + '@1') in self.assets and self.api.is_images_equal(image, self.assets[sample + '@1'])
        return is0 or is1
