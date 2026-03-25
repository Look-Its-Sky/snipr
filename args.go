package main

import (
	"flag"
)

var disableDB *bool = flag.Bool("disable_db", false, "Disable database")
var verbose *bool = flag.Bool("verbose", false, "Enable verbose output")
var exchangeDebug *bool = flag.Bool("debug_exchange", false, "Kill respective go routine after first event on contract")
