package types

type Project struct {
	ProjectID   string `json:"projectID"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Repo        string `json:"repo"`
}

type RegisterProject struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Repo        string `json:"repo"`
}
