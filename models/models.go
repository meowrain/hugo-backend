package models

type HugoNewRequest struct {
	ContentName string `json:"contentName"`
}

type HugoUpdateRequest struct {
	FilePath string `json:"filePath"`
	Content  string `json:"content"`
}

type PostContent struct {
	Content string `json:"content"`
}

type Post struct {
	Name string `json:"name"`
	Path string `json:"path"`
}

type Directory struct {
	Name     string      `json:"name"`
	Path     string      `json:"path"`
	Posts    []Post      `json:"posts"`
	Children []Directory `json:"children"`
}
