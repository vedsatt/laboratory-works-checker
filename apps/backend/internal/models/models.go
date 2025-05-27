package models

type (
	LabRequest struct {
		ID     int            `json:"id"`
		LabNum int            `json:"lab_number"`
		Code   string         `json:"code"`
		Tasks  map[string]int `json:"tasks"`
	}

	CheckerResponse struct {
		ID     int    `json:"id"`
		ResMsg string `json:"res_msg"`
	}
)
