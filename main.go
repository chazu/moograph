package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

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

var dbLines []string

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

// Returns the index of the last line of the contents
// TODO Should we refactor this to use a copy of the relevant lines? Not
// sure what the best practice is in go
// TODO Return an error instead of a random-ass sentinel value
func objectContentsListEndingIndex(bounds [2]int, dbLines []string) int {

	// First six lines (0-5) are already parsed
	scanStartIdx := bounds[0] + 6

	for i, str := range dbLines[scanStartIdx:bounds[1]] {
		// fmt.Println(str)
		// fmt.Println(i)
		if str == "-1" {
			return i + scanStartIdx
		}
	}
	// IT NEVER ENDS AAAAAAAAAAA
	return -999
}

// Returns the index of the end of the child list.
// TODO Use errors not sentinel values
func objectChildListEndingIndex(startIdx int, dbLines []string) int {
	for i, str := range dbLines[startIdx:] {
		if str == "-1" {
			return i + startIdx
		}
	}

	// u wot m8?
	return -999
}

// Process the object at the passed-in bounds
func processObject(bounds [2]int, dbLines []string) (Object, error) {
	startIdx := bounds[0]

	num, err := strconv.Atoi(strings.Trim(dbLines[startIdx], "#"))
	name := dbLines[startIdx+1]
	handles := dbLines[startIdx+2]
	flags := dbLines[startIdx+3]
	owner, err := strconv.Atoi(dbLines[startIdx+4])
	if err != nil {
		return Object{}, fmt.Errorf("Error parsing current object owner: %v", err)
	}

	location, err := strconv.Atoi(dbLines[startIdx+5])

	if err != nil {
		return Object{}, fmt.Errorf("Error parsing current location: %v", err)
	}

	contentListEndIndex := objectContentsListEndingIndex(bounds, dbLines)
	parentIndex := contentListEndIndex + 1
	childListStartIndex := parentIndex + 1
	childListEndIndex := objectChildListEndingIndex(childListStartIndex, dbLines)

	contentList := dbLines[6:contentListEndIndex]
	parent, err := strconv.Atoi(dbLines[contentListEndIndex])
	childList := dbLines[childListStartIndex:childListEndIndex]

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
func GetObjectBlockStartLineIndex(dbLines []string) (int, error) {
	for i := 4; i < len(dbLines); i++ {
		fmt.Printf("%d: %s\n", i, dbLines[i])
		if lineStartsObjectDefinition(dbLines[i]) {
			fmt.Println("AAAAAAAAAA")
			fmt.Println(dbLines[i])
			return i, nil
		}
	}

	return 0, errors.New("Unable to find start of object block")
}

// Get the index of the line which starts the verb block
func GetVerbBlockStartLineIndex(dbLines []string, objStartIdx int) (int, error) {
	for i := objStartIdx; i < len(dbLines); i++ {
		if lineStartsVerbBlock(dbLines[i]) {
			return i, nil
		}
	}
	return -1, fmt.Errorf("Error finding verb block start line: reached end of DB")
}

// Return a slice of slices containing the start and end bounds of all
// object definitions in the db, starting at specified index
func getObjDefinitionBounds(dbLines []string, objStartIdx int, objEndIdx int) ([][2]int, error) {
	doneParsingRecycled := false
	currentObjStartIdx := -1
	var result [][2]int

	for i := objStartIdx; i <= objEndIdx; i++ {

		// Handle recycled objects
		if doneParsingRecycled == false {
			if lineIsRecycledObject(dbLines[i]) {
				result = append(result, [2]int{i, i})
			} else {
				doneParsingRecycled = true
			}
		}

		if lineStartsObjectDefinition(dbLines[i]) {
			if currentObjStartIdx > 0 {
				// Already inside an object, note the bounds before
				// saving state of new object
				result = append(result, [2]int{currentObjStartIdx, i - 1})
			}
			currentObjStartIdx = i
		}
	}

	return result, nil
}

func printAround(lines []string, i int) {
	var indicator string

	for j := i - 4; j < i+5; j++ {
		if i == j-1 {
			indicator = "  <-- HERE\n"
		} else {
			indicator = "\n"
		}
		fmt.Printf("%d: %s %s", j, lines[j], indicator)
	}
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
	// header := parseHeader(lines)

	objStartIdx, err := GetObjectBlockStartLineIndex(lines)
	if err != nil {
		fmt.Errorf("Error getting start of object block: %v", err)
	}

	verbStartIdx, err := GetVerbBlockStartLineIndex(lines, objStartIdx)
	if err != nil {
		fmt.Errorf("Error getting start of verb block: %v", err)
	}

	objEndIdx := verbStartIdx - 1

	// TODO Get
	// TODO Determine format of final four blocks (clocks, queued tasks,
	// suspended tasks, active connections w/ listeners and get relevant bounds

	fmt.Printf("Start of object block: index %d\n", objStartIdx)
	printAround(lines, objStartIdx-1)
	oDefBounds, err := getObjDefinitionBounds(lines, objStartIdx, objEndIdx)

	fmt.Println(oDefBounds[:3])
	var objInstances []Object
	for _, b := range oDefBounds {
		obj, err := processObject(b, lines)
		if err != nil {
			fmt.Errorf("Error processing object: %v", err)
		}
		objInstances = append(objInstances, obj)
	}

	spew.Dump(objInstances[0])
	// z := lineIsRecycledObject("#4 recycled")
	// spew.Dump(z)
}
