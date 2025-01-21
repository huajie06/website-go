package notes

import "time"

// above lines are scratch
type Note struct {
	Fname, Content string
	NoteTimestamp  time.Time
}

type Note_Metadata struct {
	Fname         string `json:"Fname"`
	NoteTimestamp string `json:"NoteTimestamp"`
}

type NewNote_Meta struct {
	Note_Metadata
	DatetimeStr string
}
