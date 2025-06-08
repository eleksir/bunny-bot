package moon

import (
	"os"
	"time"
)

// MyConfig describes bot config file format.
type MyConfig struct {
	Name     string `json:"name,omitempty"`
	Token    string `json:"token,omitempty"`
	LogLevel string `json:"loglevel,omitempty"`
	LogFile  string `json:"logfile,omitempty"`
	LogChats bool   `json:"logchats,omitempty"`
	CSign    string `json:"csign,omitempty"`
	DataDir  string `json:"datadir,omitempty"`
}

// CasFalse describes response from api.cas.chat when user is not blacklisted.
type CasFalse struct {
	Ok          bool   `json:"ok"`
	Description string `json:"description"`
}

// CasFalse describes response from api.cas.chat when user is blacklisted.
type CasTrue struct {
	Ok     bool `json:"ok"`
	Result struct {
		Reasons   []int     `json:"reasons"`
		Offenses  int       `json:"offenses"`
		TimeAdded time.Time `json:"time_added"`
	} `json:"result"`
}

type Chatlog struct {
	LogFilePath string
	File        *os.File
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
