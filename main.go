package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"

	arg "github.com/alexflint/go-arg"
	"github.com/chazu/moograph/pkg/database"
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

// // Takes a string representing a recycled object and parses it into an Object
// // instance. Fairly naive, but recycled object lines shouldn't require more than
// // this.
// func ObjectFromRecycledLine(line string) Object {
// 	split := strings.Split(line, " ")
// 	numString := strings.TrimLeft(split[0], "#")
// 	num, err := strconv.Atoi(numString)

// 	if err != nil {
// 		fmt.Errorf("Error parsing recycled line %s : %v", line, err)
// 	}

// 	return Object{
// 		Number:   num,
// 		Recycled: true,
// 	}
// }

// // Predicate functions using regex - fairly self-explanatory
// func lineIsRecycledObject(line string) bool {
// 	if regexRecycled.MatchString(line) {
// 		return true
// 	}
// 	return false
// }

// func lineStartsObjectDefinition(line string) bool {
// 	if regexStartsObjDef.MatchString(line) {
// 		return true
// 	}
// 	return false
// }

// func lineStartsVerbBlock(line string) bool {
// 	if regexStartsVerbBlock.MatchString(line) {
// 		return true
// 	}
// 	return false
// }

// // Returns the index of the last line of the contents
// // TODO Should we refactor this to use a copy of the relevant lines? Not
// // sure what the best practice is in go
// // TODO Return an error instead of a random-ass sentinel value
// func objectContentsListEndingIndex(bounds [2]int, dbLines []string) int {

// 	// First six lines (0-5) are already parsed
// 	scanStartIdx := bounds[0] + 6

// 	for i, str := range dbLines[scanStartIdx:bounds[1]] {
// 		// fmt.Println(str)
// 		// fmt.Println(i)
// 		if str == "-1" {
// 			return i + scanStartIdx
// 		}
// 	}
// 	// IT NEVER ENDS AAAAAAAAAAA
// 	return -999
// }

// // Returns the index of the end of the child list.
// // TODO Use errors not sentinel values
// func objectChildListEndingIndex(startIdx int, dbLines []string) int {
// 	for i, str := range dbLines[startIdx:] {
// 		if str == "-1" {
// 			return i + startIdx
// 		}
// 	}

// 	// u wot m8?
// 	return -999
// }

// // Return the index of the first line of the object block
// // TODO Check the number of players against the line found here
// // to ensure integrity
// func GetObjectBlockStartLineIndex(dbLines []string) (int, error) {
// 	for i := 4; i < len(dbLines); i++ {
// 		// fmt.Printf("%d: %s\n", i, dbLines[i])
// 		if lineStartsObjectDefinition(dbLines[i]) {
// 			// fmt.Println("AAAAAAAAAA")
// 			// fmt.Println(dbLines[i])
// 			return i, nil
// 		}
// 	}

// 	return 0, errors.New("Unable to find start of object block")
// }

// // Get the index of the line which starts the verb block
// func GetVerbBlockStartLineIndex(dbLines []string, objStartIdx int) (int, error) {
// 	for i := objStartIdx; i < len(dbLines); i++ {
// 		if lineStartsVerbBlock(dbLines[i]) {
// 			return i, nil
// 		}
// 	}
// 	return -1, fmt.Errorf("Error finding verb block start line: reached end of DB")
// }

// // Return a slice of slices containing the start and end bounds of all
// // object definitions in the db, starting at specified index
// func gettObjDefinitionBounds(dbLines []string, objStartIdx int, objEndIdx int) ([][2]int, error) {
// 	doneParsingRecycled := false
// 	currentObjStartIdx := -1
// 	var result [][2]int

// 	for i := objStartIdx; i <= objEndIdx; i++ {

// 		// Handle recycled objects
// 		if doneParsingRecycled == false {
// 			if lineIsRecycledObject(dbLines[i]) {
// 				result = append(result, [2]int{i, i})
// 			} else {
// 				doneParsingRecycled = true
// 			}
// 		}

// 		if lineStartsObjectDefinition(dbLines[i]) {
// 			if currentObjStartIdx > 0 {
// 				// Already inside an object, note the bounds before
// 				// saving state of new object
// 				result = append(result, [2]int{currentObjStartIdx, i - 1})
// 			}
// 			currentObjStartIdx = i
// 		}
// 	}

// 	return result, nil
// }

// func printAround(lines []string, i int) {
// 	var indicator string

// 	for j := i - 4; j < i+5; j++ {
// 		if i == j-1 {
// 			indicator = "  <-- HERE\n"
// 		} else {
// 			indicator = "\n"
// 		}
// 		fmt.Printf("%d: %s %s", j, lines[j], indicator)
// 	}
// }

func main() {
	// Handle args
	arg.MustParse(&args)

	// Read in the DB and preprocess
	b, err := os.ReadFile(args.Filename)
	if err != nil {
		fmt.Errorf("Error opening database file: %v", err.Error())
	}

	lines := strings.Split(string(b), "\n")
	db := database.NewDatabase(lines)
	db.Parse()
}
