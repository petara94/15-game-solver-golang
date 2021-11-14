package game

import (
	"crypto/md5"
	"encoding/json"
	"errors"
)

//Movable объект примения на себе методы интерфейса
//возращает изменненую копию согласно методу
type Movable interface {
	Up() (*State, error)
	Down() (*State, error)
	Left() (*State, error)
	Right() (*State, error)
}

//State - состояние игрового поля пятнашек
type State struct {
	matrix
	hash       [16]byte
	IsWin      bool
	IsWinnable bool
}

//Hash getter
func (n *State) Hash() [16]byte {
	return n.hash
}

//Matrix getter
func (n *State) Matrix() *matrix {
	return &n.matrix
}

//checkWin узнает является ли состояние выигрышным
func (n *State) checkWin() bool {
	pre := n.matrix.Data[0][0]

	for i, row := range n.matrix.Data {
		for j, el := range row {
			if i == 0 && j == 0 {
				continue
			}
			if i == n.N-1 && j == n.M-1 {
				break
			}
			if el-pre != 1 {
				return false
			}
			pre = el
		}
	}

	return n.matrix.Data[n.N-1][n.M-1] == 0
}

func (n *State) Down() (*State, error) {
	zeroPos, e := FindEl(n.matrix.Data, 0)
	if e != nil {
		panic(e)
	}
	if zeroPos[0] == 0 {
		return nil, errors.New("can't move up")
	}

	state := n.CopySelf()

	swapInt(&state.matrix.Data[zeroPos[0]][zeroPos[1]], &state.matrix.Data[zeroPos[0]-1][zeroPos[1]])
	state.genHash()
	state.IsWin = state.checkWin()

	return state, nil
}

func (n *State) Up() (*State, error) {
	zeroPos, e := FindEl(n.matrix.Data, 0)
	if e != nil {
		panic(e)
	}
	if zeroPos[0] == n.matrix.N-1 {
		return nil, errors.New("can't move down")
	}

	state := n.CopySelf()

	swapInt(&state.matrix.Data[zeroPos[0]][zeroPos[1]], &state.matrix.Data[zeroPos[0]+1][zeroPos[1]])
	state.genHash()
	state.IsWin = state.checkWin()

	return state, nil
}

func (n *State) Right() (*State, error) {
	zeroPos, e := FindEl(n.matrix.Data, 0)
	if e != nil {
		panic(e)
	}
	if zeroPos[1] == 0 {
		return nil, errors.New("can't move left")
	}

	state := n.CopySelf()

	swapInt(&state.matrix.Data[zeroPos[0]][zeroPos[1]], &state.matrix.Data[zeroPos[0]][zeroPos[1]-1])
	state.genHash()
	state.IsWin = state.checkWin()

	return state, nil
}

func (n *State) Left() (*State, error) {
	zeroPos, e := FindEl(n.matrix.Data, 0)
	if e != nil {
		panic(e)
	}
	if zeroPos[1] == n.matrix.M-1 {
		return nil, errors.New("can't move right")
	}

	state := n.CopySelf()

	swapInt(&state.matrix.Data[zeroPos[0]][zeroPos[1]], &state.matrix.Data[zeroPos[0]][zeroPos[1]+1])
	state.genHash()
	state.IsWin = state.checkWin()

	return state, nil
}

//NewState Гнерирует состояние со случайной матрицой
func NewState() *State {
	state := &State{matrix: *NewMatrix(FieldSize, FieldSize)}

	RandFillMatrix(&state.matrix)
	state.IsWin = state.checkWin()
	state.genHash()

	return state
}

//NewStateFromMatrix возвращает состояние от матрицы
func NewStateFromMatrix(m *matrix) *State {
	state := &State{matrix: *m}

	state.IsWin = state.checkWin()
	state.genHash()

	return state
}

//CopySelf Возвращает копию состояния
func (n *State) CopySelf() *State {
	state := &State{matrix: *NewMatrix(n.N, n.N)}

	for i, row := range n.matrix.Data {
		copy(state.matrix.Data[i], row)
	}

	return state
}

func swapInt(a *int, b *int) {
	*a, *b = *b, *a
}

//genHash генерирует массив 16 байт - хеш сумму от положения элементов
//на матрице, для более быстрого стравнения состояний
func (n *State) genHash() {
	var arrBytes []byte

	for _, row := range n.matrix.Data {
		jsonBytes, _ := json.Marshal(row)
		arrBytes = append(arrBytes, jsonBytes...)
	}
	n.hash = md5.Sum(arrBytes)
}

//CheckWinnable Рассчитывает четность матрицы, узнавая
//тем самым является ли она решаема
func (n *State) CheckWinnable() bool {
	sum := 0

	for i, row := range n.matrix.Data {
		for j, el := range row {
			if el == 0 {
				sum += i + 1
			}
			sum += n.parity(i, j)
		}
	}

	n.IsWinnable = sum%2 == 0

	return n.IsWinnable
}

//parity Находит количество элементов меньших a[i][j] элемента
func (n *State) parity(i, j int) int {
	parity := 0
	pos := i*n.M + j
	for pos != n.N*n.M {
		if n.Data[pos/n.M][pos%n.M] == 0 {
			pos++
			continue
		}
		if n.Data[pos/n.M][pos%n.M] < n.Data[i][j] {
			parity++
		}
		pos++
	}
	return parity
}

//GenChilds Находит возможных потомков
func (n *State) GenChilds() []*State {
	l, _ := n.Left()
	r, _ := n.Right()
	u, _ := n.Up()
	d, _ := n.Down()

	var res []*State
	ch := []*State{l, r, u, d}

	for _, c := range ch {
		if c != nil {
			res = append(res, c)
		}
	}

	return res
}
