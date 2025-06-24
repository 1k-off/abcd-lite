package deployment

type Credentials struct {
	Username    string
	Password    string
	LoginServer string
}

type PackageInfo struct {
	Credentials Credentials
	Name        string
	Version     string
}

const defaultLoginServer = "docker.io"

func NewCredentials(username, password, loginServer string) Credentials {
	if username == "" || password == "" {
		return Credentials{}
	}
	if loginServer == "" {
		loginServer = defaultLoginServer
	}

	return Credentials{
		Username:    username,
		Password:    password,
		LoginServer: loginServer,
	}
}

func NewPackageInfo(name, version string, credentials Credentials) PackageInfo {
	return PackageInfo{
		Credentials: credentials,
		Name:        name,
		Version:     version,
	}
}
