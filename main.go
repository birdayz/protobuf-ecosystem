package main

import (
	"context"
	"encoding/json"
	"fmt"

	"cloud.google.com/go/bigquery/storage/managedwriter"
	"cloud.google.com/go/bigquery/storage/managedwriter/adapt"
	"github.com/GoogleCloudPlatform/protoc-gen-bq-schema/v2/pkg/converter"
	"github.com/GoogleCloudPlatform/protoc-gen-bq-schema/v2/protos"
	simplev1 "github.com/birdayz/protobuf-ecosystem/pkg/bqschema/proto/gen/go/simple/v1"
	"github.com/google/uuid"
	"google.golang.org/api/bigquery/v2"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"google.golang.org/protobuf/types/pluginpb"
)

const (
	projectID = "home-net-284509"
	dataset   = "nerdentest"
	tableName = "test1"
)

func main() {
	ctx := context.Background()
	client, err := managedwriter.NewClient(ctx, projectID)
	if err != nil {
		panic(err)
	}

	defer client.Close()

	service, err := bigquery.NewService(ctx)
	if err != nil {
		panic(err)
	}

	// Can't use descriptor proto, if we load from files/BSR/

	var m simplev1.ExampleTable
	m.Id = uuid.NewString()
	m.SomeData = "test123"
	m.Timestamp = timestamppb.Now()
	// m.NewField = "def"

	bts, err := proto.Marshal(&m)
	if err != nil {
		panic(err)
	}

	// fmt.Println(descriptorProto.String())

	// fmt.Println("xx", descriptorProto.ProtoReflect().Descriptor().ParentFile().Path())
	// fmt.Println("yy", m.ProtoReflect().Descriptor().ParentFile().Path())

	// fd := descriptorProto.ProtoReflect().Descriptor().ParentFile()

	// fd := &descriptorpb.FileDescriptorProto{
	// 	Package: ptr.To(string(descriptorProto.ProtoReflect().Descriptor().ParentFile().Package().Name())),
	// 	Name:    ptr.To("test.proto"),
	// 	Extension: []*descriptorpb.FieldDescriptorProto{
	// 		{
	// 			Extendee: ptr.To("google.protobuf.MessageOptions"),
	// 		},
	// 	},
	// 	MessageType: []*descriptorpb.DescriptorProto{
	// 		descriptorProto,
	// 	},
	// }
	fd := protodesc.ToFileDescriptorProto(m.ProtoReflect().Descriptor().ParentFile())

	// messageDescriptorPb := protodesc.ToDescriptorProto(m.ProtoReflect().Descriptor())
	messageDescriptorPb := fd.MessageType[0]
	messageDescriptorPb.Options = &descriptorpb.MessageOptions{}

	proto.SetExtension(messageDescriptorPb.Options, protos.E_BigqueryOpts, &protos.BigQueryMessageOptions{
		TableName: "test1"})

	fmt.Println(messageDescriptorPb.String())

	// spew.Dump(descriptorProto.Options)

	res, err := converter.Convert(&pluginpb.CodeGeneratorRequest{
		FileToGenerate: []string{m.ProtoReflect().Descriptor().ParentFile().Path()},
		ProtoFile: []*descriptorpb.FileDescriptorProto{
			fd,
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("rez", res.File)
	table, err := service.Tables.Get(projectID, dataset, tableName).Do()
	if err != nil {
		panic(err)
	}

	schemaJSON := *res.File[0].Content

	// Translate JSON of protoc plugin to
	var bigquerySDKSchema = []*bigquery.TableFieldSchema{}
	if err := json.Unmarshal([]byte(schemaJSON), &bigquerySDKSchema); err != nil {
		panic(err)
	}

	table.Schema = &bigquery.TableSchema{
		Fields: bigquerySDKSchema,
	}

	// May be special logic to never remove a field.

	table, err = service.Tables.Patch(projectID, dataset, tableName, table).Do()
	if err != nil {
		panic(err)
	}

	// Calc. diff

	// Check if doable.

	// Use SQL to remove fields.

	// Either hard fail, or remove from message descriptor that do not work out, so we send as much as possible.

	// Delete field from BQ vs New Table vs keep in BQ schema.

	// Set Description.

	_ = bts
	////////////////////////////////////////////// BQ: use ADAPT.

	// Auto create table.

	// Self Describing Message

	normalizedForBigqueryMessageDescriptor, err := adapt.NormalizeDescriptor(m.ProtoReflect().Descriptor())
	if err != nil {
		panic(err)
		// TODO: Handle error.
	}

	stream, err := client.NewManagedStream(ctx,
		managedwriter.WithDestinationTable(fmt.Sprintf("projects/%s/datasets/%s/tables/%s", projectID, dataset, tableName)),
		managedwriter.WithSchemaDescriptor(messageDescriptorPb))
	if err != nil {
		panic(err)
	}

	// Upsert https://cloud.google.com/bigquery/docs/change-data-capture
	//
	result, err := stream.AppendRows(ctx, [][]byte{bts}, managedwriter.UpdateSchemaDescriptor(normalizedForBigqueryMessageDescriptor))
	if err != nil {
		panic(err)
	}

	rez, err := result.GetResult(ctx)
	if err != nil {
		panic(err)
	}

	_ = rez
}
