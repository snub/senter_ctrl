package main

import (
	"log"
	"os"
)

var logger = log.New(os.Stderr, "senter-ctrl: ", log.Ldate|log.Lmicroseconds|log.Ltime|log.Lshortfile)
