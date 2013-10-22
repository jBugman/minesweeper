#! /usr/bin/env python
# -*- coding: UTF-8 -*-


class WindowException(Exception):
    def __init__(self, value):
        self.value = value

    def __str__(self):
        return repr(self.value)


class NotImplementedException(Exception):
    def __init__(self, value=None):
        pass

    def __str__(self):
        return 'Method not implemented!'
