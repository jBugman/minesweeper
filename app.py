#! /usr/bin/env python
# -*- coding: UTF-8 -*-

import random

from grabber import Grabber
from consts import *
from point import Point
from mouse import Mouse

class Game:
	def __init__(self):
		self.grabber = Grabber()
		self.turn = 0
		self.cellOffset = Point(16, 16)
	
	def printField(self, field):
		for y in range(FIELD_SIZE):
			print '[ ' + '  '.join(field[y]) + ' ]'
	
	def clickCell(self, cell):
		point = Point(cell.x * CELL_SIZE, cell.y * CELL_SIZE)
		Mouse.click(self.grabber.getOffset() + point + self.cellOffset)
	
	def makeTurn(self):
		self.grabber.api.activateWindow()
		self.turn += 1
		print '== Turn {} =='.format(self.turn)
		
		field = self.grabber.getField()
		self.printField(field)
		
		if len(self.getCellsByType('#', field)):
			print '[w] Unknown cells. Check it!'
			return False
		
		hiddenCells = self.getCellsByType('?', field)
		if len(hiddenCells) == FIELD_SIZE * FIELD_SIZE:
			print '[i] All cells are hidden. Clicking randomly=)'
			randomCell = hiddenCells[int(len(hiddenCells) * random.random())]
			self.clickCell(randomCell)
			return True
			
		for cell in self.getCellsByType('1', field):
			neighbours = self.getNeighbours(cell)
			
		print '[i] Don’t know what to do=('
		return False
	
	def run(self):
		while self.makeTurn():
			pass
	
	def getNeighbours(self, cell):
		cells = []
		for x in (cell.x - 1, cell.x, cell.x + 1):
			for y in (cell.y - 1, cell.y, cell.y + 1):
				if x > 0 and x < FIELD_SIZE and y > 0 and y < FIELD_SIZE and not (x == cell.x and y == cell.y):
					cells.append(Point(x, y))
		return cell
	
	def getCellsByType(self, type, field):
		cells = []
		for x in range(FIELD_SIZE):
			for y in range(FIELD_SIZE):
				if field[y][x] == type:
					cells.append(Point(x, y))
		return cells

if __name__ == '__main__':
	game = Game()
	game.run()

