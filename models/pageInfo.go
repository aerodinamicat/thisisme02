package models

type PageInfo struct {
	Token int `json:"token"`

	OrderBy string `json:"orderBy"`
	Size    int    `json:"size"`

	TotalPages int `json:"totalPages"`
	TotalItems int `json:"totalItems"`
}
