package main

import (
	"flag"
	"fmt"
)

var disableDB *bool
var verbose *bool

func init() {
	disableDB = flag.Bool("disable_db", false, "Disable database")

	fmt.Println(*disableDB)
}
