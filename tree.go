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

type pathNode struct {
	// Given the next segment s, if edges[s] exists, then we'll look there first.
	edges map[string]*pathNode
	// If set, failure to match on edges will match on wildcard
	wildcard *pathNode
	// If set, and we have nothing left to match, then we match on this node
	leaves []*pathLeaf
	// If true, this pathNode has a pathparam that matches the rest of the path
	matchesFullPath bool
}

func newPathNode() *pathNode {
	return &pathNode{edges: make(map[string]*pathNode)}
}
