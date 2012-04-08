#! /usr/bin/env python
# -*- coding: UTF-8 -*-

import copy

class Point:
	x = 0
	y = 0

	def __init__(self, x, y):
		self.x = x
		self.y = y

	def __str__(self):
		return '({0}, {1})'.format(self.x, self.y)
	
	def __repr__(self):
		return '({0}, {1})'.format(self.x, self.y)

	def __add__(self, other):
		return Point(self.x + other.x, self.y + other.y)

	def asTuple(self):
		return (self.x, self.y)