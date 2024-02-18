package object

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
)

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
	num, err := strconv.Atoi(strings.Trim(o.Lines[0], "#"))
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

	// All this index shit is fucked i'm sure. Need to review the spec and
	// stop shoving things around willy nilly
	o.setContentsListEndIndex()
	o.setChildListEndIndex()

	o.parentIndex = o.contentsListEndIndex + 1
	o.childListStartIndex = o.parentIndex + 1
	parentID, err := strconv.Atoi(o.Lines[o.parentIndex])

	contentList := o.Lines[6:o.contentsListEndIndex]
	runtime.Breakpoint()
	childList := o.Lines[o.childListStartIndex:o.childListEndIndex]

	finalObj := Object{
		Number:      num,
		Name:        name,
		handles:     handles,
		Flags:       flags,
		Owner:       owner,
		Location:    location,
		ContentList: contentList,
		ParentID:    parentID,
		ChildList:   childList,
	}

	return &finalObj, nil
}
