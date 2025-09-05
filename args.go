package main

import (
	"flag"
)

var disableDB *bool = flag.Bool("disable_db", false, "Disable database")
var verbose *bool = flag.Bool("verbose", false, "Enable verbose output")
