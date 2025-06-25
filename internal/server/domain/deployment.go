package domain

type Credentials struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type PackageInfo struct {
	Credentials Credentials `json:"credentials"`
	PackageRef  string      `json:"package_ref"`
}
