package main

import (
	"ctrl/cmd"
	_ "embed"
)

//go:embed .sql/schema.sql
var schema string

func main() {
	cmd.Execute()
}

func init() {
	cmd.Schema = schema
}
