package main

import (
	"log"
	"os"

	"github.com/ciiim/cloudborad/cmd/backpack"
	_ "github.com/ciiim/cloudborad/cmd/backpack"
)

func main() {
	if err := backpack.App.Run(os.Args); err != nil {
		log.Fatalln(err)
	}
}
