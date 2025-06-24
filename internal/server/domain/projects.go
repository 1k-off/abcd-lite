package domain

type Project struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	IISSites  []string `json:"iisSites"`
	APIKeys   []APIKey `json:"apiKeys"`
	CreatedAt string   `json:"createdAt"`
	UpdatedAt string   `json:"updatedAt"`
}

type ProjectsResponse struct {
	Projects []Project `json:"projects"`
}

type ProjectResponse struct {
	Project Project `json:"project"`
}

type APIKey struct {
	ID        string `json:"id"`
	Hash      string `json:"hash"`
	CreatedAt string `json:"createdAt"`
	Prefix    string `json:"prefix"` // first 4 chars
	Suffix    string `json:"suffix"` // last 4 chars
}

type APIKeyResponse struct {
	APIKey     string `json:"apiKey"`
	APIKeyMeta APIKey `json:"apiKeyMeta"`
}
