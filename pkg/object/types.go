package object

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
	Number             int
	Recycled           bool
	Name               string
	handles            string
	Flags              string
	Owner              int
	Location           int
	ContentList        []string
	FirstContainedItem int
	NextColocatedItem  int
	FirstChild         int
	NextSibling        int
	ParentID           int
	VerbCount          int
	Lines              []string
	// Contents List Starting index is always 6 (The 7th Line)
	contentsListEndIndex int
	childListStartIndex  int
	childListEndIndex    int
	parentIndex          int
}
