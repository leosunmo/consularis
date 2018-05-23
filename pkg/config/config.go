package config

type Config struct {
	Kubeconfig string `json:"kubeconfig"`
	MasterURL  string `json:"master"`
	Namespace  string `json:"ns"`
	Consul     string `json:"consul"`
	Port       string `json:"port"`
}

// New creates new config object
func New() *Config {
	c := &Config{}

	return c
}
