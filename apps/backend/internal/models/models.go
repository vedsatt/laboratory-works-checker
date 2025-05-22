package models

type (
	LabRequest struct {
		ID     int            `json:"id"`
		LabNum int            `json:"lab_number"`
		Code   string         `json:"code"`
		Tasks  map[string]int `json:"tasks"`
	}

	CheckerResponse struct {
		LabID       int    `json:"lab_id"`
		ResStatus   string `json:"res_status"`
		ResMsg      string `json:"res_msg"`
		SystemError string `json:"system_error"`
	}
)
