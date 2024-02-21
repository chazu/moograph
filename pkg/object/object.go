package object

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var regexRecycled = regexp.MustCompile(`^#\d+\srecycled$`)

func NewFromLines(lines []string) *Object {
	// Deep copy lines
	c := make([]string, len(lines))
	copy(c, lines)

	o := Object{
		Lines: c,
	}
	o.Parse()

	return &o
}

func (o *Object) isRecycled() bool {
	if regexRecycled.MatchString(o.Lines[0]) {
		return true
	}
	return false
}

// TODO THIS SHIT DONT WERK - DO TDD
// Sets the index of the last line of the contents
// TODO Return an error instead of a random-ass sentinel value
func (o *Object) setContentsListEndIndex() error {
	startIdx := 6
	// First six lines (0-5) are already parsed

	for i, str := range o.Lines[startIdx:] {
		if str == "-1" {
			o.contentsListEndIndex = i + startIdx
		}
	}

	return fmt.Errorf("Error setting contents list ending index: end of object defintion reached")
}

// Set the index of the end of the child list.
func (o *Object) setChildListEndIndex() error {
	startIdx := o.contentsListEndIndex + 1
	for i, str := range o.Lines[startIdx:] {
		if str == "-1" {
			o.childListEndIndex = i + startIdx
		}
	}

	return fmt.Errorf("Error setting child list ending index: end of object defintion reached")

}

func (o *Object) Parse() (*Object, error) {
	fmt.Println("Parsing Object")
	num, err := strconv.Atoi(strings.Trim(o.Lines[0], "#"))

	if o.isRecycled() {
		o.Recycled = true
		return o, nil
	}

	name := o.Lines[1]
	handles := o.Lines[2]
	flags := o.Lines[3]
	owner, err := strconv.Atoi(o.Lines[4])
	if err != nil {
		return &Object{}, fmt.Errorf("Error parsing current object owner: %v", err)
	}

	location, err := strconv.Atoi(o.Lines[5])
	if err != nil {
		return &Object{}, fmt.Errorf("Error parsing current location: %v", err)
	}
	firstContainedItem, err := strconv.Atoi(oLines[6])
	if err != nil {
		return &Object{}, fmt.Errorf("Error parsing first contained item: %v", err)
	}

	nextColocatedItem, err := strconv.Atoi(o.Lines[7])
	if err != nil {
		return &Object{}, fmt.Errorf("Error parsing next colocated item: %v", err)
	}

	parentID, err := strconv.Atoi(o.Lines[8])
	if err != nil {
		return &Object{}, fmt.Errorf("Error parsing parent ID: %v", err)
	}

	firstChild, err := strconv.Atoi(o.Lines[9])
	if err != nil {
		return &Object{}, fmt.Errorf("Error parsing child ID: %v", err)
	}

	nextSibling, err := strconv.Atoi(o.Lines[10])
	if err != nil {
		return &Object{}, fmt.Errorf("Error parsing next sibling ID: %v", err)
	}

	verbCount, err := strconv.Atoi(o.Lines[11])
	if err != nil {
		return &Object{}, fmt.Errorf("Error parsing verb count for object: %v", err)
	}

	finalObj := Object{
		Number:             num,
		Name:               name,
		handles:            handles,
		Flags:              flags,
		Owner:              owner,
		Location:           location,
		FirstContainedItem: firstContainedItem,
		NextColocatedItem:  nextColocatedItem,
		FirstChild:         firstChild,
		NextSibling:        nextSibling,
		VerbCount:          verbCount,
		ParentID:           parentID,
	}

	return &finalObj, nil
}
