package main

import (
	"context"
	"net/http"

	"github.com/redpanda-data/benthos/v4/public/service"

	// Import full suite of FOSS connect plugins
	_ "github.com/redpanda-data/connect/public/bundle/free/v4"

	// Or, in order to import both FOSS and enterprise plugins, replace the
	// above with:
	// _ "github.com/redpanda-data/connect/public/bundle/enterprise/v4"

	// Add your plugin packages here
	_ "github.com/birdayz/protobuf-ecosystem/redpandaconnect/bundle/v1"

	_ "net/http"
	_ "net/http/pprof"
)

func init() {
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()

}

func main() {
	service.RunCLI(context.Background())
}
