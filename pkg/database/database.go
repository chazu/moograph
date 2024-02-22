package database

import (
	"fmt"
	"regexp"
	"strconv"

	obj "github.com/chazu/moograph/pkg/object"
	"github.com/davecgh/go-spew/spew"
)

// Regexes used in predicate functions below
var regexRecycled = regexp.MustCompile(`^#\d+\srecycled$`)
var regexStartsObjDef = regexp.MustCompile(`^#\d+$`)
var regexStartsVerbBlock = regexp.MustCompile(`^#\d+:\d+$`)

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

func NewDatabase(lines []string) *Database {
	return &Database{
		Lines: lines,
	}
}

// Return the index of the first line of the object block
// TODO Check the number of players against the line found here
// to ensure integrity
func (d Database) GetObjectBlockStartLineIndex() int {
	pBounds := d.playerListBounds()

	return pBounds[1] + 1
}

// Get the index of the line which starts the verb block
func (d Database) GetVerbBlockStartLineIndex() (int, error) {
	for i := d.GetObjectBlockStartLineIndex(); i < len(d.Lines); i++ {
		if lineStartsVerbBlock(d.Lines[i]) {
			return i, nil
		}
	}
	return -1, fmt.Errorf("Error finding verb block start line: reached end of DB")
}

// Return a slice of slices containing the start and end bounds of all
// object definitions in the db, starting at specified index
func (d Database) getObjDefinitionBounds(objStartIdx int, objEndIdx int) ([][2]int, error) {
	doneParsingRecycled := false
	currentObjStartIdx := -1
	var result [][2]int

	for i := objStartIdx; i <= objEndIdx; i++ {

		// Handle recycled objects
		if doneParsingRecycled == false {
			if lineIsRecycledObject(d.Lines[i]) {
				result = append(result, [2]int{i, i})
				continue
			} else {
				doneParsingRecycled = true
			}
		}

		if lineStartsObjectDefinition(d.Lines[i]) {
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

// Beginning and ending indices of player list.
// Per https://www.mars.org/home/rob/docs/lmdb.html, the start is always 5
func (d Database) playerListBounds() [2]int {
	return [2]int{5, d.Header.TotalPlayerCount + 5}
}

func (d Database) Parse() error {
	d.ParseHeader()

	pBounds := d.playerListBounds()
	// TODO Copy here
	d.Players = d.Lines[pBounds[0] : pBounds[1]+1]

	objStartIdx := d.GetObjectBlockStartLineIndex()
	verbStartIdx, err := d.GetVerbBlockStartLineIndex()
	if err != nil {
		fmt.Errorf("Error getting start of verb block: %v", err)
	}

	objEndIdx := verbStartIdx - 1

	// fmt.Printf("Start of object block: index %d\n", objStartIdx)
	// printAround(lines, objStartIdx-1)
	oDefBounds, err := d.getObjDefinitionBounds(objStartIdx, objEndIdx)

	// fmt.Println(oDefBounds[:3])
	var objInstances []obj.Object
	for _, b := range oDefBounds {
		theseLines := d.Lines[b[0]:b[1]]
		o := obj.NewFromLines(theseLines)
		if err := o.Parse(); err != nil {
			return fmt.Errorf("Error parsing object: %v", err)
		}
		spew.Dump(o)
		objInstances = append(objInstances, *o)
	}

	return nil
}

func (d Database) ParseHeader() error {
	versionString := d.Lines[0]

	// Get Object Count
	totalObjectCount, err := strconv.Atoi(d.Lines[1])
	if err != nil {
		return fmt.Errorf("Error parsing object count in db header: %v", err)
	}

	// Get Verb Count
	totalVerbCount, err := strconv.Atoi(d.Lines[2])
	if err != nil {
		return fmt.Errorf("Error parsing verb count in db header: %v", err)
	}
	// Dummy line - might as well capture it and make sure its 0 - if it ain't,
	// then we may have a malformed DB or we're overlooking something we don't
	// know about the db format.
	dummyLine, err := strconv.Atoi(d.Lines[3])
	if err != nil {
		return fmt.Errorf("Error parsing dummy line in db header: %v", err)
	}

	if dummyLine != 0 {
		return fmt.Errorf("Error parsing DB header: Dummy line is not 0 - you may have a malformed db.")
	}

	// Get playercount
	playerCount, err := strconv.Atoi(d.Lines[4])
	if err != nil {
		return fmt.Errorf("Error parsing player count in db header: %v", err)
	}

	// Shove it into the struct and gtfo
	d.Header = DbHeader{
		VersionString:    versionString,
		TotalObjectCount: totalObjectCount,
		TotalVerbCount:   totalVerbCount,
		TotalPlayerCount: playerCount,
	}

	return nil
}
