package structures

const (
	JOURNAL_TEXT_ENTRY    string = "text"
	JOURNAL_MEETING_ENTRY string = "meeting"
	JOURNAL_VISION_ENTRY  string = "vision"
	JOURNAL_NOTE_ENTRY    string = "note"
)

type JournalEntry struct {
	Type  string `json:"type"`
	Title string `json:"title"`
	Data  string `json:"data"`
}

type JournalTextEntry struct {
	VoiceMode bool   `json:"voice_mode"`
	Response  string `json:"response"`
}
