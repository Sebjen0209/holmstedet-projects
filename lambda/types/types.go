package types

type Project struct {
	ProjectID   string `json:"projectID"`
	Title       string `json:"title"`
	Description string `json:"desription"`
	Repo        string `json:"repo"`
}
