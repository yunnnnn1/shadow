#!/usr/bin/env  python
# -*- coding:utf-8 -*-
#@Time    : 2022/1/25 9:47 下午
#@Author  : wangjun
#@File    : test.py

# please import your package
import time
import logging
import logging.handlers
import sys
import time

def initlogger():
    logfile ='/tmp/log.txt'
    handler = logging.handlers.RotatingFileHandler(logfile,maxBytes=1024 * 1024 * 100 ,backupCount=3)
    fmt= '%(asctime)s - %(name)s - %(levelname)-8s - %(message)s'
    formatter = logging.Formatter(fmt)
    handler.setFormatter(formatter)

    console = logging.StreamHandler()
    console.setFormatter(formatter)

    logger = logging.getLogger("test")
    logger.addHandler(handler)
    logger.addHandler(console)
    logger.setLevel(logging.INFO)

    return logger



if __name__ == '__main__':
    logger = initlogger()

    cnt = 0
    while 1 :
        print ("print_test")
        logger.info('logger_test')
        time.sleep(1)
        cnt = cnt +1
