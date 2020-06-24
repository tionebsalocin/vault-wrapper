package runner

type Config struct {
	VaultPath      []string
	VaultAddr      string
	VaultAuthPath  string
	VaultUser      string
	VaultUserRealm string
	VaultSPN       string
	Krb5           string
	Keytab         string
	EnvPrefix      string
	TokenOnly      bool
	ExportOnly     bool
	CommandLine    []string
}

type Secret struct {
	Name string
	Path string
	Key  string
}
