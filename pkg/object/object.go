package object

import (
	"fmt"
	"strconv"
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

// Sets the index of the last line of the contents
// TODO Return an error instead of a random-ass sentinel value
func (o *Object) setContentsListEndingIndex() error {
	startIdx := 6
	// First six lines (0-5) are already parsed

	for i, str := range o.Lines[startIdx:] {
		if str == "-1" {
			o.contentsListEndIdx = i + startIdx
		}
	}

	return fmt.Errorf("Error setting contents list ending index: end of object defintion reached")
}

// Set the index of the end of the child list.
// TODO Use errors not sentinel values
func (o *Object) setObjectChildListEndingIndex() error {
	startIdx := o.contentsListEndIdx + 1
	for i, str := range o.Lines[startIdx:] {
		if str == "-1" {
			o.childListEndIdx = i + startIdx
		}
	}

	return fmt.Errorf("Error setting child list ending index: end of object defintion reached")

}

func (o *Object) Parse() (*Object, error) {
	num, err := strconv.Atoi(o.Lines[0], "#")
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

	o.setObjectContentsListEndingIndex()
	parentIndex := contentListEndIndex + 1
	childListStartIndex := parentIndex + 1
	childListEndIndex := o.setObjectChildListEndingIndex(childListStartIndex)

	contentList := o.Lines[6:contentListEndIndex]
	parent, err := strconv.Atoi(o.Lines[contentListEndIndex])
	childList := o.Lines[childListStartIndex:childListEndIndex]

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
}
