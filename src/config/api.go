package config

type ApiSetting struct {
	ApiToken  string `toml:"api_token"`
	Email     string `toml:"email"`
	Domain    string `toml:"domain"`
	SubDomain string `toml:"sub_domain"`
}
