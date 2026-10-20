package elog

// The chat-logs root: where the tree lives, and whether logging is on at all.
//
// The state itself is in globals.go with the rest of the package-level variables;
// this file holds only the accessors, because "is it on and where" is one question
// with one answer and must not be asked in two pieces.

// ChatLoggingEnabled reports whether chat logging is currently switched on.
//
// Exported so another package can assert the post-shutdown state instead of
// inferring it from the absence of files.
func ChatLoggingEnabled() (string, bool) {
	return chatLogRoot()
}

// chatLogRoot returns the configured root and reports whether chat logging is
// on, in a *single* read.
//
// The two must not be asked separately: a Close landing between an enabled()
// check and a rootDir() read would leave the caller with an empty root, and
// filepath.Join("", ...) silently yields a path relative to the working
// directory — chat logs would be written wherever the supervisor happened to
// start the process.
func chatLogRoot() (string, bool) {
	stateMu.RLock()
	defer stateMu.RUnlock()

	return root, root != ""
}

// setRoot switches chat logging on or off.
func setRoot(dir string) {
	stateMu.Lock()
	defer stateMu.Unlock()

	root = dir

	if dir == "" {
		failing = false
	}
}

/* vim: set ft=go noet ai ts=4 sw=4 sts=4: */
