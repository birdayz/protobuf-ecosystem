package bigquerywriteapi

import (
	"context"
	"errors"
	"fmt"
	"os"
	"time"

	"cloud.google.com/go/bigquery/storage/managedwriter"
	"cloud.google.com/go/bigquery/storage/managedwriter/adapt"
	"github.com/redpanda-data/benthos/v4/public/service"
	"github.com/twmb/franz-go/pkg/kgo"
	"google.golang.org/api/bigquery/v2"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/reflect/protoregistry"
	"google.golang.org/protobuf/types/descriptorpb"
)

func Ctor(conf *service.ParsedConfig, mgr *service.Resources) (out service.BatchOutput, batchPolicy service.BatchPolicy, maxInFlight int, err error) {
	bp, err := conf.FieldBatchPolicy("batching")
	if err != nil {
		panic(err)
	}

	project, err := conf.FieldString("project")
	if err != nil {
		return nil, bp, 0, err
	}
	dataset, err := conf.FieldString("dataset")
	if err != nil {
		return nil, bp, 0, err
	}
	table, err := conf.FieldString("table")
	if err != nil {
		return nil, bp, 0, err
	}

	protoMsg, err := conf.FieldString("protobuf_type")
	if err != nil {
		panic(err)
	}

	fds, err := conf.FieldString("file_descriptor_set")
	if err != nil {
		return nil, bp, 0, err
	}

	data, err := os.ReadFile(fds)
	if err != nil {
		return nil, bp, 0, err
	}

	// --- Load protobuf descriptor set.
	// From file atm, but we can load it from anywhere; BSR or even from the message itself (interpolation?)/from self describing message
	var fileDescriptorSet descriptorpb.FileDescriptorSet
	if err := proto.Unmarshal(data, &fileDescriptorSet); err != nil {
		return nil, bp, 0, err
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

	return &Output{
		log:                mgr.Logger().With("table", table),
		project:            project,
		dataset:            dataset,
		table:              table,
		descriptorForWrite: descriptorProtoForWrite,
	}, bp, 4, nil

}

func NewOutput(desc *descriptorpb.DescriptorProto) *Output {
	return &Output{
		log:                &service.Logger{},
		project:            "home-net-284509",
		dataset:            "nerdentest",
		table:              "test1",
		descriptorForWrite: desc,
	}
}

func init() {
	configSpec := service.NewConfigSpec()
	configSpec.Field(service.NewStringField("project").Default(""))
	configSpec.Field(service.NewStringField("dataset").Default(""))
	configSpec.Field(service.NewStringField("table").Default(""))
	configSpec.Field(service.NewStringField("file_descriptor_set"))
	configSpec.Field(service.NewStringField("protobuf_type"))
	configSpec.Field(service.NewBatchPolicyField("batching"))

	err := service.RegisterBatchOutput("bigquery_write_api", configSpec, Ctor)
	if err != nil {
		panic(err)
	}
}

type Output struct {
	log *service.Logger

	project string
	dataset string
	table   string

	writeClient   *managedwriter.Client
	serviceClient *bigquery.Service
	writeStream   *managedwriter.ManagedStream

	descriptorForWrite *descriptorpb.DescriptorProto
}

func (b *Output) Connect(ctx context.Context) error {
	writeClient, err := managedwriter.NewClient(ctx, b.project, managedwriter.WithMultiplexing())
	if err != nil {
		return fmt.Errorf("failed to create managedwriter: %w", err)
	}
	b.writeClient = writeClient

	writeStream, err := writeClient.NewManagedStream(ctx,
		managedwriter.WithDestinationTable(fmt.Sprintf("projects/%s/datasets/%s/tables/%s", b.project, b.dataset, b.table)),
		managedwriter.WithSchemaDescriptor(b.descriptorForWrite))
	if err != nil {
		panic(err)
	}
	b.writeStream = writeStream

	serviceClient, err := bigquery.NewService(ctx)
	if err != nil {
		return fmt.Errorf("failed to create serviceClient: %w", err)
	}
	b.serviceClient = serviceClient

	return nil
}

func (b *Output) WriteKgoBatch(ctx context.Context, batch []*kgo.Record) error {
	// b.log.With("batch_size", len(batch)).Debug("Received batch")
	var rows [][]byte
	for _, msg := range batch {
		bts := msg.Value
		rows = append(rows, bts)
	}

	fmt.Println(len(batch))
	startTime := time.Now()
	res, err := b.writeStream.AppendRows(ctx, rows)
	if err != nil {
		panic(err)
	}
	_, _ = res.GetResult(ctx)
	// b.log.With("duration", time.Since(startTime)).Debug("AppendRows finished")
	fmt.Println("Done", time.Since(startTime))

	return nil
}

func (b *Output) WriteBatch(ctx context.Context, batch service.MessageBatch) error {
	b.log.With("batch_size", len(batch)).Debug("Received batch")
	var rows [][]byte
	for _, msg := range batch {
		bts, err := msg.AsBytes()
		if err != nil {
			panic(err)
		}
		rows = append(rows, bts)
	}

	startTime := time.Now()
	res, err := b.writeStream.AppendRows(ctx, rows)
	if err != nil {
		panic(err)
	}
	_, _ = res.GetResult(ctx)
	b.log.With("duration", time.Since(startTime)).Debug("AppendRows finished")

	return nil
}

func (b *Output) Close(ctx context.Context) error {
	var errs error

	if err := b.writeStream.Close(); err != nil {
		errs = errors.Join(errs, fmt.Errorf("failed to close writeStream: %w", err))
	}

	if err := b.writeClient.Close(); err != nil {
		errs = errors.Join(errs, fmt.Errorf("failed to close writeClient: %w", err))
	}

	return errs
}
