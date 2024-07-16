package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hungerarray/go-png/png"
)

func main() {
	filePath := flag.String("filepath", "", "png file path")
	flag.Parse()

	if filePath == nil || *filePath == "" {
		fmt.Fprintf(os.Stderr, "usage: %s\n", os.Args[0])
		flag.PrintDefaults()
		os.Exit(-1)
	}

	png, err := png.NewPng(*filePath)
	if err != nil {
		log.Fatalf("[error]: %s\n", err.Error())
	}
	png.PrintInfo()
}
