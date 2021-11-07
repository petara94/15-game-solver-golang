package game

import (
	"sync"
)

type Node struct {
	Val    *State
	childs []*Node
	Parent *Node
}

func NewNode(val *State, parent *Node) *Node {
	return &Node{Val: val, Parent: parent, childs: []*Node{}}
}

func (n *Node) Add(val *State) *Node {
	ch := NewNode(val, n)
	n.childs = append(n.childs, ch)
	return ch
}

type Graph struct {
	StartPoint *Node
	nodes      []*Node
}

func (g *Graph) Hashes() []*Node {
	return g.nodes
}

func NewGraph(startPoint *State) *Graph {
	g := &Graph{StartPoint: NewNode(startPoint, nil)}
	//g.nodes[g.StartPoint.Val.Hash()] = g.StartPoint

	g.nodes = []*Node{g.StartPoint}

	return g
}

func (g *Graph) BuildGraph(from *Node) {

	if from == nil {
		from = g.StartPoint
	}

	moves := []func() (*State, error){from.Val.Up, from.Val.Left, from.Val.Down, from.Val.Right}

	for _, moveAction := range moves {
		state, err := moveAction()
		if err == nil {
			//if _, ok := g.nodes[state.Hash()]; !ok {
			//	g.nodes[state.Hash()] = from.Add(state)
			//	g.BuildGraph(g.nodes[state.Hash()])
			//}
			next := from.Add(state)
			g.nodes = append(g.nodes, next)
			g.BuildGraph(next)
		}
	}

}

func (g *Graph) SearchWinPath() []*Node {
	res := make([]*Node, 0)
	nodes := append([]*Node{}, g.nodes...)
	node := nodes[0]
	isFound := false

	for len(nodes) > 0 || !node.Val.IsWin {
		if !node.Val.IsWin {
			moves := []func() (*State, error){
				node.Val.Up,
				node.Val.Left,
				node.Val.Down,
				node.Val.Right,
			}

			for _, action := range moves {
				nodeAfterMove, _ := action()
				if nodeAfterMove == nil {
					continue
				}
				node.Add(nodeAfterMove)
			}
			nodes = append(nodes, nodes[0].childs...)

		} else {
			isFound = true
			break
		}
		if len(nodes) == 1{
			break
		}
		node = nodes[1]
		nodes = nodes[1:]
	}

	if !isFound {
		return res
	}

	res = append(res, node)
	current := node.Parent
	for current != nil {
		res = append([]*Node{current}, res...)
		current = current.Parent
	}

	g.nodes = append([]*Node{}, res...)
	return res
}

func IsInIntAsync(arr []*Node, el *Node) bool {
	res := false
	var wg sync.WaitGroup

	step := len(arr) / 12
	for i := 0; i < len(arr); i += step {
		if len(arr)-i < step {
			step = len(arr) - i
		}
		wg.Add(1)
		go func(from, to int) {
			defer wg.Done()
			for j := from; j < to; j++ {
				if arr[j].Val.hash == el.Val.hash {
					res = true
				}
			}
		}(i, i+step)
	}
	wg.Wait()
	return res
}
