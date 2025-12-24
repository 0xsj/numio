package main

import (
	"context"
	"fmt"
	"time"

	"github.com/0xsj/numio/pkg/engine"
)

func main() {
	eng := engine.New()

	// Refresh rates from network
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	n, err := eng.RefreshRates(ctx)
	if err != nil {
		fmt.Printf("Refresh failed: %v (using cached/default rates)\n", err)
	} else {
		fmt.Printf("Refreshed %d rates\n", n)
	}

	// Now use the engine
	result := eng.Eval("$100 in EUR")
	fmt.Println(result)

	result = eng.Eval("1 BTC in USD")
	fmt.Println(result)

	result = eng.Eval("1 XAU in USD")
	fmt.Println(result)
}
