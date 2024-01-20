package grom

import (
	"regexp"
	"strings"
)

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

func (leaf *pathLeaf) match(wildcardValues []string) bool {
	if leaf.regexps == nil {
		return true
	}

	if len(leaf.regexps) != len(wildcardValues) {
		panic("bug: invariant violated")
	}

	for i, r := range leaf.regexps {
		if r != nil {
			if !r.MatchString(wildcardValues[i]) {
				return false
			}
		}
	}
	return true
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

// Segments is like ["admin", "users"] representing "/admin/users"
// wildcardValues are the actual values accumulated when we match on a wildcard.
func (pn *pathNode) match(segments []string, wildcardValues []string) (leaf *pathLeaf, wildcardMap map[string]string) {
	// Handle leaf nodes:
	if len(segments) == 0 {
		for _, leaf := range pn.leaves {
			if leaf.match(wildcardValues) {
				return leaf, makeWildcardMap(leaf, wildcardValues)
			}
		}
		return nil, nil
	}

	var seg string
	seg, segments = segments[0], segments[1:]
	subPn, ok := pn.edges[seg]
	if ok {
		leaf, wildcardMap = subPn.match(segments, wildcardValues)
	}

	if leaf == nil && pn.wildcard != nil {
		leaf, wildcardMap = pn.wildcard.match(segments, append(wildcardValues, seg))
	}

	if leaf == nil && pn.matchesFullPath {
		for _, leaf := range pn.leaves {
			if leaf.matchesFullPath && leaf.match(wildcardValues) {
				if len(wildcardValues) > 0 {
					wcVals := []string{wildcardValues[len(wildcardValues)-1], seg}
					for _, s := range segments {
						wcVals = append(wcVals, s)
					}
					wildcardValues[len(wildcardValues)-1] = strings.Join(wcVals, "/")
				}
				return leaf, makeWildcardMap(leaf, wildcardValues)
			}
		}
		return nil, nil
	}
	return leaf, wildcardMap
}

func makeWildcardMap(leaf *pathLeaf, wildcards []string) map[string]string {
	if leaf == nil {
		return nil
	}

	leafWildcards := leaf.wildcards
	if len(wildcards) == 0 || (len(leafWildcards) != len(wildcards)) {
		return nil
	}

	// At this point, we know that wildcards and leaf.wildcards match in length.
	assoc := make(map[string]string)
	for i, w := range wildcards {
		assoc[leafWildcards[i]] = w
	}
	return assoc
}
