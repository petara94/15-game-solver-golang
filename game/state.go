package game

import (
	"crypto/md5"
	"encoding/json"
	"errors"
)

type Movable interface {
	Up() (*State, error)
	Down() (*State, error)
	Left() (*State, error)
	Right() (*State, error)
}

type State struct {
	matrix
	hash       [16]byte
	IsWin      bool
	IsWinnable bool
}

func (s *State) Hash() [16]byte {
	return s.hash
}

func (s *State) Matrix() *matrix {
	return &s.matrix
}

func (s *State) checkWin() bool {
	pre := s.matrix.Data[0][0]

	for i, row := range s.matrix.Data {
		for j, el := range row {
			if i == 0 && j == 0 {
				continue
			}
			if i == s.N-1 && j == s.M-1 {
				break
			}
			if el-pre != 1 {
				return false
			}
			pre = el
		}
	}

	return s.matrix.Data[s.N-1][s.M-1] == 0
}

func (s *State) Down() (*State, error) {
	zeroPos, e := FindEl(s.matrix.Data, 0)
	if e != nil {
		panic(e)
	}
	if zeroPos[0] == 0 {
		return nil, errors.New("can't move up")
	}

	state := s.CopySelf()

	swapInt(&state.matrix.Data[zeroPos[0]][zeroPos[1]], &state.matrix.Data[zeroPos[0]-1][zeroPos[1]])
	state.genHash()
	state.IsWin = state.checkWin()

	return state, nil
}

func (s *State) Up() (*State, error) {
	zeroPos, e := FindEl(s.matrix.Data, 0)
	if e != nil {
		panic(e)
	}
	if zeroPos[0] == s.matrix.N-1 {
		return nil, errors.New("can't move down")
	}

	state := s.CopySelf()

	swapInt(&state.matrix.Data[zeroPos[0]][zeroPos[1]], &state.matrix.Data[zeroPos[0]+1][zeroPos[1]])
	state.genHash()
	state.IsWin = state.checkWin()

	return state, nil
}

func (s *State) Right() (*State, error) {
	zeroPos, e := FindEl(s.matrix.Data, 0)
	if e != nil {
		panic(e)
	}
	if zeroPos[1] == 0 {
		return nil, errors.New("can't move left")
	}

	state := s.CopySelf()

	swapInt(&state.matrix.Data[zeroPos[0]][zeroPos[1]], &state.matrix.Data[zeroPos[0]][zeroPos[1]-1])
	state.genHash()
	state.IsWin = state.checkWin()

	return state, nil
}

func (s *State) Left() (*State, error) {
	zeroPos, e := FindEl(s.matrix.Data, 0)
	if e != nil {
		panic(e)
	}
	if zeroPos[1] == s.matrix.M-1 {
		return nil, errors.New("can't move right")
	}

	state := s.CopySelf()

	swapInt(&state.matrix.Data[zeroPos[0]][zeroPos[1]], &state.matrix.Data[zeroPos[0]][zeroPos[1]+1])
	state.genHash()
	state.IsWin = state.checkWin()

	return state, nil
}

func NewState() *State {
	state := &State{matrix: *NewMatrix(FieldSize, FieldSize)}

	RandFillMatrix(&state.matrix)
	state.IsWin = state.checkWin()
	state.genHash()

	return state
}

func NewStateFromMatrix(m *matrix) *State {
	state := &State{matrix: *m}

	state.IsWin = state.checkWin()
	state.genHash()

	return state
}

func (s *State) CopySelf() *State {
	state := &State{matrix: *NewMatrix(s.N, s.N)}

	for i, row := range s.matrix.Data {
		copy(state.matrix.Data[i], row)
	}

	return state
}

func swapInt(a *int, b *int) {
	*a, *b = *b, *a
}

func (s *State) genHash() {
	var arrBytes []byte

	for _, row := range s.matrix.Data {
		jsonBytes, _ := json.Marshal(row)
		arrBytes = append(arrBytes, jsonBytes...)
	}
	s.hash = md5.Sum(arrBytes)
}

func (s *State) CheckWinnable() bool {
	sum := 0

	for i, row := range s.matrix.Data {
		for j, el := range row {
			if el == 0 {
				sum += i + 1
			}
			sum += s.parity(i, j)
		}
	}

	s.IsWinnable = sum%2 == 0

	return s.IsWinnable
}

func (s *State) parity(i, j int) int {
	parity := 0
	pos := i*s.M + j
	for pos != s.N*s.M {
		if s.Data[pos/s.M][pos%s.M] ==0 {
			pos++
			continue
		}
		if s.Data[pos/s.M][pos%s.M] < s.Data[i][j] {
			parity++
		}
		pos++
	}
	return parity
}
