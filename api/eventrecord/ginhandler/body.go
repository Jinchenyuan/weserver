package ginhandler

type EventNote struct {
	ID        string  `json:"id"`
	CreatedAt string  `json:"createdAt"`
	Title     string  `json:"title"`
	Content   string  `json:"content"`
	PhotoURI  *string `json:"photoUri"`
}

type EventRecordDetail struct {
	ID            string      `json:"id"`
	Title         string      `json:"title"`
	OccurredAt    string      `json:"occurredAt"`
	Location      string      `json:"location"`
	Description   string      `json:"description"`
	CoverPhotoURI string      `json:"coverPhotoUri"`
	CreatedAt     string      `json:"createdAt"`
	UpdatedAt     string      `json:"updatedAt"`
	Notes         []EventNote `json:"notes"`
}

type EventRecordSummary struct {
	ID             string  `json:"id"`
	Title          string  `json:"title"`
	Location       string  `json:"location"`
	OccurredAt     string  `json:"occurredAt"`
	CoverPhotoURI  string  `json:"coverPhotoUri"`
	LatestNoteDate *string `json:"latestNoteDate"`
	NoteCount      int32   `json:"noteCount"`
	UpdatedAt      string  `json:"updatedAt"`
}

type CreateEventRecordRequest struct {
	Title         string `json:"title"`
	OccurredAt    string `json:"occurredAt"`
	Location      string `json:"location"`
	Description   string `json:"description"`
	CoverPhotoURI string `json:"coverPhotoUri"`
}

type UpdateEventRecordRequest struct {
	ID            string `json:"id"`
	Title         string `json:"title"`
	OccurredAt    string `json:"occurredAt"`
	Location      string `json:"location"`
	Description   string `json:"description"`
	CoverPhotoURI string `json:"coverPhotoUri"`
}

type CreateEventNoteRequest struct {
	Title    string  `json:"title"`
	Content  string  `json:"content"`
	PhotoURI *string `json:"photoUri"`
}

type UpdateEventNoteRequest struct {
	EventID  string  `json:"eventId"`
	NoteID   string  `json:"noteId"`
	Title    string  `json:"title"`
	Content  string  `json:"content"`
	PhotoURI *string `json:"photoUri"`
}

type ListEventRecordsResponse struct {
	Events []EventRecordSummary `json:"events"`
}

type EventRecordMutationResponse struct {
	Success bool              `json:"success"`
	Event   EventRecordDetail `json:"event"`
	Message *string           `json:"message"`
}
