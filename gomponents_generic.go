//go:build go1.18
// +build go1.18

package gomponents

// Map a slice of anything to a slice of Nodes.
func Map[T any](ts []T, cb func(T) Node) []Node {
	var nodes []Node
	for _, t := range ts {
		nodes = append(nodes, cb(t))
	}
	return nodes
}
