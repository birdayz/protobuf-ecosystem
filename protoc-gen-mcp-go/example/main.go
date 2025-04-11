package main

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create MCP server
	s := server.NewMCPServer(
		"Example auto-generated gRPC-MCP",
		"1.0.0",
	)

	srv := exampleServer{}

	RegisterExampleServiceMCP(s, &srv)

	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

type exampleServer struct {
}

func (t *exampleServer) CreateExample(ctx context.Context, in *CreateExampleRequest) (*CreateExampleResponse, error) {
	return &CreateExampleResponse{
		SomeString: "HAHA " + in.GetNested().GetNested2().GetNested3().GetOptionalString(),
	}, nil
}
