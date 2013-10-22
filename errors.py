# -*- coding: UTF-8 -*-


class WindowException(Exception):
    pass


class NotImplementedException(Exception):
    def __init__(self, value=None):
        pass

    def __str__(self):
        return 'Method not implemented!'
