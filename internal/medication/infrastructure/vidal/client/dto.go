package client

import "github.com/FSO-VK/final-project-vk-backend/internal/medication/infrastructure/vidal"

// Response represents the top-level structure of the API response from vidal.ru.
type Response struct {
	Success    bool            `json:"success"`
	Products   []vidal.Product `json:"products"`
	Pagination Pagination      `json:"pagination"`
}

// Pagination represents the pagination details of the response.
type Pagination struct {
	Page            int `json:"page"`
	Limit           int `json:"limit"`
	PageCount       int `json:"pageCount"`
	ItemsCount      int `json:"itemsCount"`
	TotalItemsCount int `json:"totalItemsCount"`
}
