package models

type PageInfo struct {
	PageToken int `json:"pageToken"`

	OrderBy  string `json:"orderBy"`
	PageSize int    `json:"pageSize"`

	TotalPages int `json:"totalPages"`
	TotalItems int `json:"totalItems"`
}
