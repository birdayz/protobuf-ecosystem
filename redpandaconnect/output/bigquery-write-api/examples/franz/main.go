package main

import (
	"context"
	"fmt"
	"os"
	"sync"

	"cloud.google.com/go/bigquery/storage/managedwriter/adapt"
	bigquerywriteapi "github.com/birdayz/protobuf-ecosystem/redpandaconnect/output/bigquery-write-api"
	"github.com/twmb/franz-go/pkg/kgo"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"

	"net/http"
	_ "net/http"
	_ "net/http/pprof"
)

func init() {
	go func() {
		http.ListenAndServe("localhost:6060", nil)
	}()
}

func main() {
	fds := "./pkg/bqschema/proto/gen/filedescriptorset.binpb"
	protoMsg := "simple.v1.ExampleTable"
	data, err := os.ReadFile(fds)
	if err != nil {
		panic(err)
	}

	// --- Load protobuf descriptor set.
	// From file atm, but we can load it from anywhere; BSR or even from the message itself (interpolation?)/from self describing message
	var fileDescriptorSet descriptorpb.FileDescriptorSet
	if err := proto.Unmarshal(data, &fileDescriptorSet); err != nil {
		panic(err)
	}

	registry := &protoregistry.Files{}
	for _, file := range fileDescriptorSet.GetFile() {
		// Convert from wire format into protoreflect format
		fd, err := protodesc.NewFile(file, registry)
		if err != nil {
			panic(err)
		}
		registry.RegisterFile(fd)
	}

	// Lookup the type.
	descriptor, err := registry.FindDescriptorByName(protoreflect.FullName(protoMsg))
	if err != nil {
		panic(err)
	}

	msgDescriptor, ok := descriptor.(protoreflect.MessageDescriptor)
	if !ok {
		panic("is not msg descriptor")
	}

	descriptorProtoForWrite, err := adapt.NormalizeDescriptor(msgDescriptor)
	if err != nil {
		panic(err)
	}

	ou := bigquerywriteapi.NewOutput(descriptorProtoForWrite)
	ou.Connect(context.Background())

	cl, err := kgo.NewClient(kgo.SeedBrokers("localhost:9092"), kgo.ConsumeTopics("cdc.ExampleTable"), kgo.ConsumerGroup("lal8"))
	if err != nil {
		panic(err)
	}

	for {
		fmt.Println("VOR")
		fetches := cl.PollFetches(context.Background())
		var wg sync.WaitGroup
		fetches.EachPartition(func(ftp kgo.FetchTopicPartition) {
			fmt.Println("B", len(ftp.Records))
			wg.Add(1)
			go func() {
				defer wg.Done()
				ou.WriteKgoBatch(context.Background(), ftp.Records)
			}()
		})
		wg.Wait()
	}

}
