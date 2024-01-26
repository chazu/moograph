package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	arg "github.com/alexflint/go-arg"
	"github.com/davecgh/go-spew/spew"
	//	"github.com/rivo/tview"
	bitfield "github.com/hymkor/go-bitfield"
)

var args struct {
	Filename string `arg:"required"`
}

type DbHeader struct {
	VersionString    string
	TotalObjectCount int
	TotalVerbCount   int
	TotalPlayerCount int
}

// TODO Since whoever implemented the non-standard Overridable
// Extension decided to break the bitfield by using 1024 instead
// of 128^2, we'll have to manually detect that part.
type PermissionSpec struct {
	Readable    uint `bit:"1"`
	Writable    uint `bit:"1"`
	Executable  uint `bit:"1"`
	Debug       uint `bit:"1"`
	DObjAny     uint `bit:"1"`
	DObjThis    uint `bit:"1"`
	IObjAny     uint `bit:"1"`
	IObjThis    uint `bit:"1"`
	Overridable uint
}

type VerbDefinition struct {
	VerbName    string
	Owner       int
	Permissions PermissionSpec
	Preposition int
}

func parseHeader(dbLines []string) *DbHeader {
	versionString := dbLines[0]

	totalObjectCount, err := strconv.Atoi(dbLines[1])
	if err != nil {
		fmt.Errorf("Error parsing object count in db header: %v", err)
	}

	totalVerbCount, err := strconv.Atoi(dbLines[2])
	if err != nil {
		fmt.Errorf("Error parsing verb count in db header: %v", err)
	}

	dummyLine, err := strconv.Atoi(dbLines[3])
	if err != nil {
		fmt.Errorf("Error parsing dummy line in db header: %v", err)
	}

	if dummyLine != 0 {
		fmt.Errorf("Error parsing DB header: Dummy line is not 0 - you may have a malformed db.")
	}

	playerCount, err := strconv.Atoi(dbLines[4])
	if err != nil {
		fmt.Errorf("Error parsing player count in db header: %v", err)
	}

	return &DbHeader{
		VersionString:    versionString,
		TotalObjectCount: totalObjectCount,
		TotalVerbCount:   totalVerbCount,
		TotalPlayerCount: playerCount,
	}
}

// Return the index of the first line of the object block
// TODO Check the number of players against the line found here
// to ensure integrity
func GetObjectBlockStartLine(dbLines []string) (int, error) {
	for i := 5; i < len(dbLines); i++ {
		if strings.HasPrefix("#", dbLines[i]) {
			return i, nil
		}
	}

	return 0, errors.New("Unable to find start of object block")

}

func main() {

	arg.MustParse(&args)

	b, err := os.ReadFile(args.Filename)
	if err != nil {
		fmt.Errorf("Error opening database file: %v", err.Error())
	}

	lines := strings.Split(string(b), "\n")
	header := parseHeader(lines)

	objStart, err := GetObjectBlockStartLine(lines)
	if err != nil {
		fmt.Errorf("Error getting start of object block: %v", err)
	}

	spew.Dump(header)
	spew.Dump(objStart)

	var perms PermissionSpec
	bitfield.Unpack(69, &perms)
	spew.Dump(perms)
	// box := tview.NewBox().SetBorder(true).SetTitle("Hello, world!")
	// if err := tview.NewApplication().SetRoot(box, true).Run(); err != nil {
	// 	panic(err)
	// }
}
