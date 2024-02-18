package database

import (
	obj "github.com/chazu/moograph/pkg/object"
)

type DbHeader struct {
	VersionString    string
	TotalObjectCount int
	TotalVerbCount   int
	TotalPlayerCount int
}

type Database struct {
	Header  DbHeader
	Lines   []string
	Objects []obj.Object
}
