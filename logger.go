package main

import (
	"log"
	"os"
)

var logger = log.New(os.Stderr, "senter-cli: ", log.LstdFlags)
