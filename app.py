#! /usr/bin/env python
# -*- coding: UTF-8 -*-

import random
import time

from grabber import Grabber
from consts import *
from point import Point
from mouse import Mouse

class Game:
	def __init__(self):
		self.grabber = Grabber()
		self.fieldSize = self.grabber.fieldSize
		print '[i] Field size:', self.fieldSize
		self.turn = 0
		self.cellOffset = Point(16, 16)
	
	def printField(self, field):
		for y in range(self.fieldSize.y):
			print '[ ' + ' '.join(field[y]) + ' ]'
	
	def clickCell(self, cell, rightButton = False):
		point = Point(cell.x * CELL_SIZE, cell.y * CELL_SIZE)
		if rightButton:
			Mouse.rightclick(self.grabber.getOffset() + point + self.cellOffset)
		else:
			Mouse.click(self.grabber.getOffset() + point + self.cellOffset)
	
	def resetField(self):
		print '[i] Blocked state. Restaring game'
		restartButton = Point(50, self.fieldSize.y * CELL_SIZE + 16)
		Mouse.click(self.grabber.getOffset() + restartButton)
	
	def makeTurn(self):
		self.turn += 1
		print '== Iteration {} =='.format(self.turn)
		
		field = self.grabber.getField()
		self.printField(field)
		
		mines = self.getCellsByType('M', field)
		if len(self.getCellsByType('!', field)) and not len(mines):
			if self.turn == 1:
				self.resetField()
				return True
			else:
				print '[i] Win! All your base are belong to us!'
				return False
		elif len(mines):
			if self.turn == 1:
				self.resetField()
				return True
			else:
				print '[i] BANG! Game Over..'
				return False
		
		if len(self.getCellsByType('#', field)):
			print '[w] Unknown cells. Check it!'
			return False
		
		hiddenCells = self.getCellsByType('?', field)
		if len(hiddenCells) == self.fieldSize.x * self.fieldSize.y:
			print '[i] All cells are hidden. Clicking randomly=)'
			randomCell = hiddenCells[int(len(hiddenCells) * random.random())]
			self.clickCell(randomCell)
			return True
		
		for i in (1, 2, 3, 4, 5, 6, 7, 8):
			cells = self.getCellsByType(str(i), field)
			for cell in cells:
				neighbours = self.getNeighbours(cell)
				hiddenNeighbours = [item for item in neighbours if field[item.y][item.x] == '?']
				flaggedNeighbours = [item for item in neighbours if field[item.y][item.x] == '+']
				if len(hiddenNeighbours) and len(hiddenNeighbours) == (i - len(flaggedNeighbours)):
					print '[i] Flagging:', hiddenNeighbours
					for nb in hiddenNeighbours:
						self.clickCell(nb, True)
					return True
				elif len(flaggedNeighbours) == i and len(hiddenNeighbours):
					print '[i] Clicking neighbours because of max flags:', hiddenNeighbours
					for nb in hiddenNeighbours:
						self.clickCell(nb)
					return True
		
		hiddenCells = self.getCellsByType('?', field)
		if len(hiddenCells):
			print '[w] Don’t know what to do=( Clicking randomly'
			randomCell = hiddenCells[int(len(hiddenCells) * random.random())]
			self.clickCell(randomCell)
			return True
		
		print '[w] Something is wrong. Don’t know how I get there'
		return False
	
	def run(self):
		self.grabber.api.activateWindow()
		try:
			while self.makeTurn():
				time.sleep(0.15)
		except KeyboardInterrupt:
			print '[i] Exit'
	
	def getNeighbours(self, cell):
		cells = []
		for x in (cell.x - 1, cell.x, cell.x + 1):
			for y in (cell.y - 1, cell.y, cell.y + 1):
				if x >= 0 and x < self.fieldSize.x and y >= 0 and y < self.fieldSize.y and not (x == cell.x and y == cell.y):
					cells.append(Point(x, y))
		return cells
	
	def getCellsByType(self, type, field):
		cells = []
		for y in range(self.fieldSize.y):
			for x in range(self.fieldSize.x):
				if field[y][x] == type:
					cells.append(Point(x, y))
		return cells

if __name__ == '__main__':
	game = Game()
	game.run()