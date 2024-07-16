package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hungerarray/go-png/png"
)

var filePath string

func init() {
	flag.StringVar(&filePath, "filepath", "", "png file path")
	flag.StringVar(&filePath, "f", "", "png file path (short)")
}

func main() {
	flag.Parse()

	if filePath == "" {
		fmt.Fprintf(os.Stderr, "usage: %s\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(-1)
	}

	png, err := png.NewPng(filePath)
	if err != nil {
		log.Fatalf("[error]: %s\n", err.Error())
	}
	png.PrintInfo()
}
