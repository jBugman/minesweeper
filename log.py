# -*- coding: UTF-8 -*-
import logging
import logging.handlers

LOG_FORMAT = '[%(asctime)-6s][%(levelname)s][%(name)s] %(message)s'
DATE_FORMAT = '%d.%m.%Y %H:%M:%S'

console_handler = logging.StreamHandler()
console_handler.setFormatter(logging.Formatter(LOG_FORMAT, datefmt=DATE_FORMAT))


def get_logger(name, debug=True):
    logger = logging.getLogger(name)
    logger.setLevel(logging.DEBUG if debug else logging.INFO)
    logger.addHandler(console_handler)
    return logger
