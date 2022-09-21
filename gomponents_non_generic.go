//go:build !go1.18
// +build !go1.18

package gomponents

// Map something enumerable to a list of Nodes.
func Map(length int, cb func(i int) Node) []Node {
	var nodes []Node
	for i := 0; i < length; i++ {
		nodes = append(nodes, cb(i))
	}
	return nodes
}
