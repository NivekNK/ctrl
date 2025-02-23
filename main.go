package main

import (
	"ctrl/cmd"
	"ctrl/util"
	_ "embed"
)

//go:embed .sql/schema.sql
var schema string

func main() {
	cmd.Execute()
}

func init() {
	util.Schema = &schema
}
