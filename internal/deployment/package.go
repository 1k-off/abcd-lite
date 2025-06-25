package deployment

type Credentials struct {
	Username string
	Password string
}

type PackageInfo struct {
	Credentials Credentials
	PackageRef  string
}

func NewCredentials(username, password string) Credentials {
	if username == "" || password == "" {
		return Credentials{}
	}
	return Credentials{
		Username: username,
		Password: password,
	}
}

func NewPackageInfo(packageRef string, credentials Credentials) PackageInfo {
	return PackageInfo{
		Credentials: credentials,
		PackageRef:  packageRef,
	}
}
