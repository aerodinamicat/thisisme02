package models

type PageInfo struct {
	Current int `json:"current"`
	Next    int `json:"next"`

	OrderBy string `json:"orderBy"`
	Size    int    `json:"size"`

	TotalPages int `json:"totalPages"`
	TotalItems int `json:"totalItems"`
}
