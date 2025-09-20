package schema

import (
	gql "github.com/mununki/gqlmerge/lib"
)

func NewSchema() *string {
	schema := gql.Merge("  ", "./schema")

	return schema
}
