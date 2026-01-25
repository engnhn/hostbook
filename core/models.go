package core

type Host struct {
	Name         string   `json:"name"`
	Hostname     string   `json:"hostname"`
	User         string   `json:"user"`
	Port         string   `json:"port"`
	IdentityFile string   `json:"identity_file,omitempty"`
	Tags         []string `json:"tags,omitempty"`
}

type Config struct {
	Hosts []Host `json:"hosts"`
}
