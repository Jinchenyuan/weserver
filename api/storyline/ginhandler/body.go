package ginhandler

type StorylineNode struct {
	ID        string  `json:"id"`
	Title     string  `json:"title"`
	Date      string  `json:"date"`
	Note      string  `json:"note"`
	Location  string  `json:"location"`
	PhotoURI  *string `json:"photoUri"`
	SortOrder int32   `json:"sortOrder"`
}

type StorylineDetail struct {
	ID            string          `json:"id"`
	Title         string          `json:"title"`
	Description   string          `json:"description"`
	CoverPhotoURI *string         `json:"coverPhotoUri"`
	CreatedAt     string          `json:"createdAt"`
	UpdatedAt     string          `json:"updatedAt"`
	Nodes         []StorylineNode `json:"nodes"`
}

type StorylineSummary struct {
	ID              string  `json:"id"`
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	CoverPhotoURI   *string `json:"coverPhotoUri"`
	NodeCount       int32   `json:"nodeCount"`
	LatestNodeTitle *string `json:"latestNodeTitle"`
	LatestNodeDate  *string `json:"latestNodeDate"`
	UpdatedAt       string  `json:"updatedAt"`
}

type StorylineNodeInput struct {
	ID        *string `json:"id"`
	Title     string  `json:"title"`
	Date      string  `json:"date"`
	Note      string  `json:"note"`
	Location  string  `json:"location"`
	PhotoURI  *string `json:"photoUri"`
	SortOrder int32   `json:"sortOrder"`
}

type CreateStorylineRequest struct {
	Title         string               `json:"title"`
	Description   string               `json:"description"`
	CoverPhotoURI *string              `json:"coverPhotoUri"`
	Nodes         []StorylineNodeInput `json:"nodes"`
}

type UpdateStorylineRequest struct {
	ID            string               `json:"id"`
	Title         string               `json:"title"`
	Description   string               `json:"description"`
	CoverPhotoURI *string              `json:"coverPhotoUri"`
	Nodes         []StorylineNodeInput `json:"nodes"`
}

type ListStorylinesResponse struct {
	Storylines []StorylineSummary `json:"storylines"`
}

type StorylineMutationResponse struct {
	Success   bool            `json:"success"`
	Storyline StorylineDetail `json:"storyline"`
	Message   *string         `json:"message"`
}
