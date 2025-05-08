package moon

type MyConfig struct {
	Name     string `json:"name,omitempty"`
	Token    string `json:"token,omitempty"`
	LogLevel string `json:"loglevel,omitempty"`
	LogFile  string `json:"logfile,omitempty"`
	CSign    string `json:"csign,omitempty"`
	DataDir  string `json:"datadir,omitempty"`
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
