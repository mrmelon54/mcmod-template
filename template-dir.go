package main

import (
	"embed"
	"log"
)

//go:embed all:template
var templateDir embed.FS

func init() {
	dir, err := templateDir.ReadDir("template")
	if err != nil {
		log.Fatal(err)
	}
	for _, i := range dir {
		println(i.Name())
	}
}
