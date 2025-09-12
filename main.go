package main

import (
	"flag"
	"log"
	"sync"

	"snipr/schemas"
	"snipr/schemas/dex"
)

func main() {
	flag.Parse()
	if *disableDB { log.Println("Database disabled, skipping connection.") } else { initDB() }
	client := auth()

	exchanges := []*schemas.Exchange{
		dex.UniswapV2(disableDB),
		dex.UniswapV3(disableDB),
		// dex.UniswapV4(disableDB),
	}

	var wg sync.WaitGroup
	for _, exchange := range exchanges {
		wg.Add(1)
		go listenForPools(exchange, &wg, client)
	}

	log.Println("Started listeners for all exchanges. Waiting for events...")
	wg.Wait()

	// Since the listeners run indefinitely, you might want to handle graceful shutdown
	// and close the channel. For this example, we'll leave it open
	// close(results)
}
