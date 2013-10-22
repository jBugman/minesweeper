#! /usr/bin/env python
# -*- coding: UTF-8 -*-
import random
import time
import sys

from grabber import Grabber
from consts import CELL_SIZE
from point import Point
from mouse import Mouse
import errors
import log


class Game:
    def __init__(self):
        self.logger = log.get_logger('game')
        self.logger.info('Starting')

        try:
            self.grabber = Grabber()
        except errors.WindowException:
            self.logger.error('Can not find game window')
            sys.exit(1)

        self.field_size = self.grabber.field_size
        self.logger.info('Field size: {}'.format(self.field_size))
        self.turn = 0
        self.cell_offset = Point(16, 16)

    def print_field(self, field):
        for y in range(self.field_size.y):
            print '[ ' + ' '.join(field[y]) + ' ]'

    def click_cell(self, cell, with_right_button=False):
        point = Point(cell.x * CELL_SIZE, cell.y * CELL_SIZE)
        if with_right_button:
            Mouse.rightclick(self.grabber.get_offset() + point + self.cell_offset)
        else:
            Mouse.click(self.grabber.get_offset() + point + self.cell_offset)

    def reset_field(self):
        self.logger.info('Blocked state. Restaring game')
        restart_button_coords = Point(50, self.field_size.y * CELL_SIZE + 16)
        Mouse.click(self.grabber.get_offset() + restart_button_coords)

    def make_turn(self):
        self.turn += 1
        self.logger.info('! Iteration {}'.format(self.turn))

        field = self.grabber.grab_field()
        self.print_field(field)

        mines = self.cells_with_type('M', field)
        if len(self.cells_with_type('!', field)) and not len(mines):
            if self.turn == 1:
                self.reset_field()
                return True
            else:
                self.logger.info('ALL YOUR BASE ARE BELONG TO US!')
                return False
        elif len(mines):
            if self.turn == 1:
                self.reset_field()
                return True
            else:
                self.logger.info('BANG! Game Over..')
                return False

        if len(self.cells_with_type('#', field)):
            self.logger.error('Unknown cells. Check it!')
            return False

        hidden_cells = self.cells_with_type('?', field)
        if len(hidden_cells) == self.field_size.x * self.field_size.y:
            self.logger.info('All cells are hidden. Clicking randomly. x_x')
            randomCell = hidden_cells[int(len(hidden_cells) * random.random())]
            self.click_cell(randomCell)
            return True

        for i in (1, 2, 3, 4, 5, 6, 7, 8):
            cells = self.cells_with_type(str(i), field)
            for cell in cells:
                neighbours = self.getNeighbours(cell)
                hidden_neighbours = [item for item in neighbours if field[item.y][item.x] == '?']
                flagged_neighbours = [item for item in neighbours if field[item.y][item.x] == '+']
                if len(hidden_neighbours) and len(hidden_neighbours) == (i - len(flagged_neighbours)):
                    self.logger.debug('Flagging: {}'.format(hidden_neighbours))
                    for nb in hidden_neighbours:
                        self.click_cell(nb, with_right_button=True)
                    return True
                elif len(flagged_neighbours) == i and len(hidden_neighbours):
                    self.logger.debug(u'Clicking neighbours because of max flags: {}'.format(hidden_neighbours))
                    for nb in hidden_neighbours:
                        self.click_cell(nb)
                    return True

        hidden_cells = self.cells_with_type('?', field)
        if len(hidden_cells):
            self.logger.warn('Don’t know what to do. Clicking randomly=(')
            random_cell = hidden_cells[int(len(hidden_cells) * random.random())]
            self.click_cell(random_cell)
            return True

        self.logger.warn('Something is wrong. Don’t know how I get there')
        return False

    def run(self):
        self.grabber.api.activate_window()
        try:
            while self.make_turn():
                time.sleep(0.15)
        except KeyboardInterrupt:
            self.logger.info('Exiting now')

    def getNeighbours(self, cell):
        cells = []
        for x in (cell.x - 1, cell.x, cell.x + 1):
            for y in (cell.y - 1, cell.y, cell.y + 1):
                if x >= 0 and x < self.field_size.x and y >= 0 and y < self.field_size.y and not (x == cell.x and y == cell.y):
                    cells.append(Point(x, y))
        return cells

    def cells_with_type(self, cell_type, field):
        cells = []
        for y in range(self.field_size.y):
            for x in range(self.field_size.x):
                if field[y][x] == cell_type:
                    cells.append(Point(x, y))
        return cells


if __name__ == '__main__':
    game = Game()
    game.run()
