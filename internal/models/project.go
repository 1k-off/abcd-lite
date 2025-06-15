package models

type Project struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	IISSites  []string `json:"iisSites"`
	APIKeys   []string `json:"apiKeys"`
	CreatedAt string   `json:"createdAt"`
	UpdatedAt string   `json:"updatedAt"`
}
