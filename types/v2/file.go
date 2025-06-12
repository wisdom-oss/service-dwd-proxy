package v2

type File struct {
	Name     string `json:"name"`
	MimeType string `json:"mime"`
	Content  string `json:"content"`
}
