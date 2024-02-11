package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	"time"

	arg "github.com/alexflint/go-arg"
	"github.com/davecgh/go-spew/spew"
	//	"github.com/rivo/tview"
)

var args struct {
	Filename string `arg:"required"`
}

// Regexes used in predicate functions below
var regexRecycled = regexp.MustCompile(`^#\d+\srecycled$`)
var regexStartsObjDef = regexp.MustCompile(`^#\d+$`)
var regexStartsVerbBlock = regexp.MustCompile(`^#\d+:\d+$`)

var currObjLines []string

type DbHeader struct {
	VersionString    string
	TotalObjectCount int
	TotalVerbCount   int
	TotalPlayerCount int
}

// PermissionSpec holds the decoded bitfield for object permissions.
// Note that bitfield annotations describe the number of bits per member,
// not the index of the bit in the field.
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

// Fairly straightforward - contains a verb definition.
type VerbDefinition struct {
	VerbName    string
	Owner       int
	Permissions PermissionSpec
	Preposition int
}

// You know, for objects.
type Object struct {
	Number      int
	Recycled    bool
	Name        string
	handles     string
	Flags       string
	Owner       int
	Location    int
	ContentList []string
	Parent      int
	ChildList   []string
}

// Takes a string representing a recycled object and parses it into an Object
// instance. Fairly naive, but recycled object lines shouldn't require more than
// this.
func ObjectFromRecycledLine(line string) Object {
	split := strings.Split(line, " ")
	numString := strings.TrimLeft(split[0], "#")
	num, err := strconv.Atoi(numString)

	if err != nil {
		fmt.Errorf("Error parsing recycled line %s : %v", line, err)
	}

	return Object{
		Number:   num,
		Recycled: true,
	}
}

// Predicate functions using regex - fairly self-explanatory
func lineIsRecycledObject(line string) bool {
	if regexRecycled.MatchString(line) {
		return true
	}

	return false
}

func lineStartsObjectDefinition(line string) bool {
	if regexStartsObjDef.MatchString(line) {
		return true
	}

	return false
}

func lineStartsVerbBlock(line string) bool {
	if regexStartsVerbBlock.MatchString(line) {
		return true
	}

	return false
}

// End of predicate functions!

// DEPRECATED - Function using global state to parse object block
// Returns the index of the last line of the contents
// TODO Refactor to avoid global state
func currentObjectContentsListEndingIndex() int {
	for i, str := range currObjLines[6:] {
		if str == "-1" {
			return i
		}
	}
	return -999
}

// DEPRECATED - Function using global state to parse object block.
// Returns the index of the end of the child list.
// TODO Refactor to avoid global state
func currentObjectChildListEndingIndex(start int) int {
	for i, str := range currObjLines[start:] {
		if str == "-1" {
			return i
		}
	}

	// u wot m8?
	return -999
}

// DEPRECATED - Function using global state to parse object block
// TODO Refactor this one to accept object start and end index, for thread safety
func processCurrentObject() (Object, error) {
	// Simple enough to parse these
	fmt.Printf("Object lines: %v\n", len(currObjLines))

	num, err := strconv.Atoi(strings.Trim(currObjLines[0], "#"))
	name := currObjLines[1]
	handles := currObjLines[2]
	flags := currObjLines[3]
	owner, err := strconv.Atoi(currObjLines[4])
	location, err := strconv.Atoi(currObjLines[5])

	if err != nil {
		return Object{}, fmt.Errorf("Error parsing current object: %v", err)
	}
	fmt.Println("Ay! Made it!")
	contentListEndIndex := currentObjectContentsListEndingIndex()
	runtime.Breakpoint()
	parentIndex := contentListEndIndex + 1
	childListStartIndex := parentIndex + 1
	childListEndIndex := currentObjectChildListEndingIndex(childListStartIndex)

	contentList := currObjLines[6:contentListEndIndex]
	parent, err := strconv.Atoi(currObjLines[contentListEndIndex])
	childList := currObjLines[childListStartIndex:childListEndIndex]

	finalObj := Object{
		Number:      num,
		Name:        name,
		handles:     handles,
		Flags:       flags,
		Owner:       owner,
		Location:    location,
		ContentList: contentList,
		Parent:      parent,
		ChildList:   childList,
	}

	return finalObj, nil
	// TODO Grab verb definitions for object
	// TODO Grab Property names
	// TODO Gtab Property definitions

	// TODO Parse Verb Block
}

func parseObjectBlock(dbLines []string, startingIndex int) ([]Object, error) {
	var objects []Object

	for i, value := range dbLines[startingIndex:] {
		fmt.Printf("Iterating over Line %d: %s\n", i, value)

		if lineIsRecycledObject(value) {
			fmt.Printf("  It is a recycled object\n")
			objects = append(objects, ObjectFromRecycledLine(value))
		} else if lineStartsVerbBlock(value) {
			fmt.Printf("  It begins the verb block\n")
			return objects, nil
		} else if lineStartsObjectDefinition(value) {
			fmt.Printf("  It starts the object definition\n")
			fmt.Printf("Tis is it: %v", value)
			time.Sleep(1 * time.Second)
			// TODO Think this is fucked - only getting one line before it runs
			if len(currObjLines) > 0 {
				fmt.Printf("  An object is waiting to be processed...\n")
				time.Sleep(1 * time.Second)
				obj, err := processCurrentObject()
				if err != nil {
					return objects, fmt.Errorf("Error processing object: %v", err)
				}
				objects = append(objects, obj)
				currObjLines = []string{}
			}

			currObjLines = append(currObjLines, value)
		} else {
			// Add to current object?
			currObjLines = append(currObjLines, value)
		}

	}

	return objects, nil
}

func parseHeader(dbLines []string) *DbHeader {
	versionString := dbLines[0]

	// Get Object Count
	totalObjectCount, err := strconv.Atoi(dbLines[1])
	if err != nil {
		fmt.Errorf("Error parsing object count in db header: %v", err)
	}

	// Get Verb Count
	totalVerbCount, err := strconv.Atoi(dbLines[2])
	if err != nil {
		fmt.Errorf("Error parsing verb count in db header: %v", err)
	}

	// Dummy line - might as well capture it and make sure its 0 - if it ain't,
	// then we may have a malformed DB or we're overlooking something we don't
	// know about the db format.
	dummyLine, err := strconv.Atoi(dbLines[3])
	if err != nil {
		fmt.Errorf("Error parsing dummy line in db header: %v", err)
	}

	if dummyLine != 0 {
		fmt.Errorf("Error parsing DB header: Dummy line is not 0 - you may have a malformed db.")
	}

	// Get playercount
	playerCount, err := strconv.Atoi(dbLines[4])
	if err != nil {
		fmt.Errorf("Error parsing player count in db header: %v", err)
	}

	// Shove it into the struct and gtfo
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
	for i := 4; i < len(dbLines); i++ {
		if strings.HasPrefix("#", dbLines[i]) {
			return i, nil
		}
	}

	return 0, errors.New("Unable to find start of object block")

}

// Get the index of the line which starts the verb block
func GetVerbBlockSartLine(dbLines []string, objStartIdx int) (int, error) {
	for i := objStartIdx; i < len(dbLines); i++ {
		if lineStartsVerbBlock(dbLines[i]) {
			return i, nil
		}
	}
	return -1, fmt.Errorf("Error finding verb block start line: reached end of DB")
}

func main() {
	// Handle args
	arg.MustParse(&args)

	// Read in the DB and preprocess
	b, err := os.ReadFile(args.Filename)
	if err != nil {
		fmt.Errorf("Error opening database file: %v", err.Error())
	}

	lines := strings.Split(string(b), "\n")

	// Get header
	header := parseHeader(lines)

	objStartIdx, err := GetObjectBlockStartLineIndex(lines)
	if err != nil {
		fmt.Errorf("Error getting start of object block: %v", err)
	}

	verbStartIdx, err := GetVerbBlockStartLineIndex(lines, objStartIdx)
	if err != nil {
		fmt.Errorf("Error getting start of verb block: %v", err)
	}

	objEndIdx := verbStartIdx - 1

	// TODO Determine format of final four blocks (clocks, queued tasks,
	// suspended tasks, active connections w/ listeners and get relevant bounds

	fmt.Printf("Start of object block: line %d", objStartIdx)
	objects, err := parseObjectBlock(lines, objStart)
	spew.Dump(objects)
	z := lineIsRecycledObject("#4 recycled")
	spew.Dump(z)
}
