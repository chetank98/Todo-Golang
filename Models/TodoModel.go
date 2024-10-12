package Models

type Todos struct {
	NoteId      string `json:"noteid" db:"id"`
	UserId      string `json:"userid" db:"userId"`
	Name        string `json:"title" db:"name"`
	Note        string `json:"user_note" db:"description"`
	IsCompleted bool   `json:"iscompleted" db:"iscompleted"`
}
