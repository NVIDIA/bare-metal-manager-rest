package main

import (
	"flag"

	gsv "github.com/nvidia/carbide-rest/site-workflow/pkg/grpc/server"
)

// Test the RLA grpc client
func main() {
	toutPtr := flag.Int("tout", 300, "grpc server timeout")
	flag.Parse()
	gsv.RlaTest(*toutPtr)
}
