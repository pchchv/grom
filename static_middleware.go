package grom

// StaticOption configures how StaticMiddlewareDir handles url paths and index files for directories.
// If set, Prefix is removed from the start of the url path before attempting to serve a directory or file.
// If set, IndexFile is the index file to serve when the url path maps to a directory.
type StaticOption struct {
	Prefix    string
	IndexFile string
}
