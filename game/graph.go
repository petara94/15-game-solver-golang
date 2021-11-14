package game

import (
	"sync"
)

//Node - элемент графа
type Node struct {
	Val        *State
	l, r, u, d *Node
	Parent     *Node
}

//NewNodeWithParent создает ноду из состояния и родителя
func NewNodeWithParent(val *State, parent *Node) *Node {
	return &Node{Val: val, Parent: parent}
}

//NewNode создает ноду из состояния
func NewNode(val *State) *Node {
	return &Node{Val: val, Parent: nil}
}

type Graph struct {
	StartPoint   *Node
	winPathWidth []*Node
	winPathDeep  []*Node
}

//NewGraph Создает граф от начальной точки
func NewGraph(startPoint *State) *Graph {
	g := &Graph{StartPoint: NewNode(startPoint)}

	return g
}

//WidthSearch Поиск в ширину
func (g *Graph) WidthSearch() []*Node {
	res := make([]*Node, 0)
	nodes := append([]*Node{}, g.StartPoint)
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

	g.winPathWidth = append([]*Node{}, res...)
	return res
}

//DeepSearch Поиск в глубину
func (g *Graph) DeepSearch(depth int) []*Node {
	wg := sync.WaitGroup{}
	var resultNode *Node = NewNodeWithParent(nil, nil)
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
	g.winPathDeep = res
	return g.winPathDeep
}

//searchAndBuildTree ищет выигрышную ноду, а если не находит запускает горутины для своих потомков
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

//genChilds генерирует возможных потомков
func (n *Node) genChilds() []*Node {

	var res []*Node
	if n.l != nil || n.r != nil || n.u != nil || n.d != nil {
		ch := []*Node{n.l, n.r, n.u, n.d}

		for _, c := range ch {
			if c != nil {
				res = append(res, c)
			}
		}

		return res
	}

	l, _ := n.Val.Left()
	r, _ := n.Val.Right()
	u, _ := n.Val.Up()
	d, _ := n.Val.Down()

	n.l = n.GenChild(l)
	n.r = n.GenChild(r)
	n.u = n.GenChild(u)
	n.d = n.GenChild(d)

	ch := []*Node{n.l, n.r, n.u, n.d}

	for _, c := range ch {
		if c != nil {
			res = append(res, c)
		}
	}

	return res
}

//GenChild создает ноду из возможного состояния
func (n *Node) GenChild(val *State) *Node {
	if val == nil {
		return nil
	}
	ch := NewNodeWithParent(val, n)
	return ch
}
