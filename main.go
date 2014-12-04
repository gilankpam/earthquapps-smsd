package main

import (
	"github.com/gilankpam/smsd/earthquapps-smsd"
)

func main() {

	smsd.Serve("config.json")
}
