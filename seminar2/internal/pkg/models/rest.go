package models

type ApiCreateTask struct {
	URL string `json:"url"`
}

type ApiCreateTaskResponse struct {
	Token  string `json:"token"`
	Status string `json:"status"`
	Info   string `json:"info"`
}

type ApiNameResponse struct {
	Name string `json:"name"`
}
