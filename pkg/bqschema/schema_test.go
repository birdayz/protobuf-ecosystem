package bqschema

import (
	"testing"

	simplev1 "github.com/birdayz/protobuf-ecosystem/pkg/bqschema/proto/gen/go/simple/v1"
	. "github.com/onsi/gomega"
	"google.golang.org/api/bigquery/v2"
)

func TestSimpleSchema(t *testing.T) {
	RegisterTestingT(t)
	p := simplev1.ExampleTable{}

	schema, err := SchemaFromProto(&p)
	Expect(err).ToNot(HaveOccurred())
	Expect(schema).To(ConsistOf(
		[]*bigquery.TableFieldSchema{
			&bigquery.TableFieldSchema{
				Name: "id",
				Type: "STRING",
				Mode: "NULLABLE",
			},
			&bigquery.TableFieldSchema{
				Name: "timestamp",
				Type: "TIMESTAMP",
				Mode: "NULLABLE",
			},
			&bigquery.TableFieldSchema{
				Name: "some_data",
				Type: "STRING",
				Mode: "NULLABLE",
			},
			&bigquery.TableFieldSchema{
				Name: "bla",
				Type: "INTEGER",
				Mode: "NULLABLE",
			},
		}))
}
