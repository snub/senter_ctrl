package senter

import (
	"log"
	"os"
)

var logger = log.New(os.Stderr, "senter: ", log.Ldate|log.Lmicroseconds|log.Ltime|log.Lshortfile)
