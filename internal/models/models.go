package models

type Data struct {
	Repository string `json:"repository"`
	Name       string `json:"name"`
	Tag        string `json:"tag"`
}

type Information struct {
	LayersCount int   `json:"layers_count"`
	TotalSize   int64 `json:"total_size"`
}
