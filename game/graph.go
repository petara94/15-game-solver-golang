package game

import (
	"sync"
)

type Node struct {
	Val        *State
	l, r, u, d *Node
	Parent     *Node
}

func NewNode(val *State, parent *Node) *Node {
	return &Node{Val: val, Parent: parent}
}

func (n *Node) GenChild(val *State) *Node {
	if val == nil {
		return nil
	}
	ch := NewNode(val, n)
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

func (g *Graph) SearchWinPath() []*Node {
	res := make([]*Node, 0)
	nodes := append([]*Node{}, g.nodes...)
	node := nodes[0]
	isFound := false

	for len(nodes) > 0 || !node.Val.IsWin {
		if !node.Val.IsWin {
			for _, child := range nodes[0].genChilds() {
				if child != nil {
					nodes = append(nodes, child)
				}
			}

		} else {
			isFound = true
			break
		}
		if len(nodes) == 1 {
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

func (g *Graph) DeepSearch(depth int) []*Node {
	wg := sync.WaitGroup{}
	var resultNode *Node = NewNode(nil, nil)
	isFound := false

	wg.Add(1)
	go g.StartPoint.searchAndBuildTree(&wg, resultNode, depth, &isFound)

	wg.Wait()

	if !isFound {
		return []*Node{}
	}

	res := []*Node{resultNode}
	current := resultNode.Parent
	for current != nil {
		res = append([]*Node{current}, res...)
		current = current.Parent
	}

	return res
}

func (n *Node) searchAndBuildTree(wg *sync.WaitGroup, result *Node, depth int, isFound *bool) {
	defer wg.Done()

	if *isFound {
		return
	}

	if n.Val.IsWin {
		*isFound = true
		*result = *n
		return
	}

	if depth == 0 {
		return
	}

	for _, ch := range n.genChilds() {
		wg.Add(1)
		go ch.searchAndBuildTree(wg, result, depth-1, isFound)
	}

}

func (n *Node) genChilds() []*Node {
	l, _ := n.Val.Left()
	r, _ := n.Val.Right()
	u, _ := n.Val.Up()
	d, _ := n.Val.Down()

	n.l = n.GenChild(l)
	n.r = n.GenChild(r)
	n.u = n.GenChild(u)
	n.d = n.GenChild(d)

	var res []*Node
	ch := []*Node{n.l, n.r, n.u, n.d}

	for _, c := range ch {
		if c != nil {
			res = append(res, c)
		}
	}

	return res
}
