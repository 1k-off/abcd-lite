package domain

type Credentials struct {
	Username    string `json:"username"`
	Password    string `json:"password"`
	LoginServer string `json:"login_server"`
}

type PackageInfo struct {
	Credentials Credentials `json:"credentials"`
	Name        string      `json:"name"`
	Version     string      `json:"version"`
}
