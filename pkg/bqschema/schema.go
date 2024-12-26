package bqschema

import (
	"encoding/json"
	"fmt"

	"github.com/GoogleCloudPlatform/protoc-gen-bq-schema/v2/pkg/converter"
	"github.com/GoogleCloudPlatform/protoc-gen-bq-schema/v2/protos"
	"google.golang.org/api/bigquery/v2"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/reflect/protodesc"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

func SchemaFromProto(msg proto.Message) ([]*bigquery.TableFieldSchema, error) {
	fd := protodesc.ToFileDescriptorProto(msg.ProtoReflect().Descriptor().ParentFile())
	messageDescriptorPb := fd.MessageType[0]

	if messageDescriptorPb.Options == nil {
		messageDescriptorPb.Options = &descriptorpb.MessageOptions{}
	}

	proto.SetExtension(messageDescriptorPb.Options, protos.E_BigqueryOpts,
		&protos.BigQueryMessageOptions{
			TableName: "fake", // We don't need the real table name, we're only interested in the fields.
		},
	)

	res, err := converter.Convert(&pluginpb.CodeGeneratorRequest{
		FileToGenerate: []string{msg.ProtoReflect().Descriptor().ParentFile().Path()},
		ProtoFile: []*descriptorpb.FileDescriptorProto{
			fd,
		},
	})
	if err != nil {
		return nil, fmt.Errorf("failed to convert proto descriptor to bigquery schema: %w", err)
	}

	if len(res.File) == 0 {
		return nil, fmt.Errorf("failed to convert proto descriptor to BigQuery schema: protoc-gen-bq-schema returned no file")
	}

	if res.File[0].Content == nil {
		return nil, fmt.Errorf("failed to convert proto descriptor to BigQuery schema: protoc-gen-bq-schema returned file with no content")
	}

	schemaJSON := *res.File[0].Content

	// Translate JSON of protoc plugin to
	var fieldSchema = []*bigquery.TableFieldSchema{}
	if err := json.Unmarshal([]byte(schemaJSON), &fieldSchema); err != nil {
		return nil, fmt.Errorf("failed to unmarshal schema JSON into *bigquery.TableFieldSchema: %w", err)
	}

	return fieldSchema, nil
}
