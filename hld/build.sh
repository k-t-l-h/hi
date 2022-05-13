#!/bin/bash

gcc -std=gnu99 -g -pg -o ws  ./main.c -levent -levent_pthreads -lpthread
