// Code generated by ent, DO NOT EDIT.

package migrate

import (
	"entgo.io/ent/dialect/sql/schema"
	"entgo.io/ent/schema/field"
)

var (
	// HTTPRulesColumns holds the columns for the "http_rules" table.
	HTTPRulesColumns = []*schema.Column{
		{Name: "id", Type: field.TypeUint32, Increment: true},
		{Name: "gw_url", Type: field.TypeString, Size: 128, Default: ""},
		{Name: "http_type", Type: field.TypeString, Size: 128, Default: ""},
		{Name: "status", Type: field.TypeUint8, Default: 0},
		{Name: "application", Type: field.TypeString, Unique: true, Size: 128},
		{Name: "interface_type", Type: field.TypeUint8, Default: 0},
		{Name: "interface_url", Type: field.TypeString, Size: 128, Default: ""},
		{Name: "config", Type: field.TypeString, Size: 2000, Default: ""},
		{Name: "create_time", Type: field.TypeUint32, Default: 0},
		{Name: "update_time", Type: field.TypeUint32, Default: 0},
	}
	// HTTPRulesTable holds the schema information for the "http_rules" table.
	HTTPRulesTable = &schema.Table{
		Name:       "http_rules",
		Columns:    HTTPRulesColumns,
		PrimaryKey: []*schema.Column{HTTPRulesColumns[0]},
	}
	// Tables holds all the tables in the schema.
	Tables = []*schema.Table{
		HTTPRulesTable,
	}
)

func init() {
}
