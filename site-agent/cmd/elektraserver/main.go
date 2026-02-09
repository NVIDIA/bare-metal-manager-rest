package main

import (
	"flag"

	gsv "github.com/nvidia/carbide-rest/site-workflow/pkg/grpc/server"
)

// Test the carbide grpc client
func main() {
	toutPtr := flag.Int("tout", 300, "grpc server timeout")
	flag.Parse()
	gsv.ForgeTest(*toutPtr)
}
