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

// Map a map[key:string]any to a slice of Nodes.
func MapKey[T any](ts map[string]T, cb func(string, T) Node) []Node {
	var nodes []Node
	for k, t := range ts {
		nodes = append(nodes, cb(k, t))
	}
	return nodes
}
