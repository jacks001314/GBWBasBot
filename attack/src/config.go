package attack

type AttackScriptConfig struct {
	scripts []*Config `json:"scripts"`
}

type Config struct {
	Enable bool `json:"enable"`

	Name     string `json:"name"`
	Author   string `json:"author"`
	Atype    string `json:"atype"`
	Language string `json:"language"`
	App      string `json:"app"`
	Id       int    `json:"id"`
	CVECode  string `json:"CVECode"`
	Desc     string `json:"desc"`
	Suggest  string `json:"suggest"`

	DefaultPort  int    `json:"defaultPort"`
	DefaultProto string `json:"defaultProto"`

	FPath string `json:"fpath"`
}
