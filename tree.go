package grom

import "regexp"

// pathLeaf represents a leaf path segment that corresponds to a single route.
// For the route /admin/forums/:forum_id:\d.*/suggestions/:suggestion_id:\d.*
// We'd have wildcards = ["forum_id", "suggestion_id"]
// and regexps = [/\d.*/, /\d.*/]
// For the route /admin/forums/:forum_id/suggestions/:suggestion_id:\d.*
// We'd have wildcards = ["forum_id", "suggestion_id"]
// and regexps = [nil, /\d.*/]
// For the route /admin/forums/:forum_id/suggestions/:suggestion_id
// We'd have wildcards = ["forum_id", "suggestion_id"]
// and regexps = nil
type pathLeaf struct {
	// names of wildcards that lead to this leaf, e.g. ["category_id"] for the wildcard ":category_id".
	wildcards []string
	// regexps corresponding to wildcards.
	// If a segment has regexp contraint, its entry will be nil.
	// If the route has no regexp contraints on any segments, then regexps will be nil.
	regexps []*regexp.Regexp
	// Pointer back to the route.
	route *route
	// If true, this leaf has a pathparam that matches the rest of the path.
	matchesFullPath bool
}
