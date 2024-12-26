package bigquerywriteapi

import "github.com/redpanda-data/benthos/v4/public/service"

func newConfigSpec() *service.ConfigSpec {
	configSpec := service.NewConfigSpec()
	// configSpec.Field(service.NewObjectField()
	configSpec.Field(service.NewStringField("project").Default(""))
	configSpec.Field(service.NewStringField("dataset").Default(""))
	configSpec.Field(service.NewStringField("table").Default(""))
	configSpec.Field(service.NewStringField("file_descriptor_set"))
	configSpec.Field(service.NewStringField("protobuf_message"))
	return configSpec
}

func setConfig() {
}
