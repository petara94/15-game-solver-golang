package game

import (
	"errors"
	"fmt"
	"log"
	"math/rand"
)

type matrix struct {
	Data [][]int
	N    int
	M    int
}

//NewMatrix создает матрицу размера N на M
func NewMatrix(n, m int) *matrix {
	data := make([][]int, n)

	if n < 0 || m < 0 {
		log.Fatal("N or M less than zero")
	}

	for i := range data {
		data[i] = make([]int, m)
	}
	res := &matrix{Data: data, N: n, M: m}
	return res
}

//PrintMatrix печатает матрицу в stdout
func (m matrix) PrintMatrix() {
	for _, row := range m.Data {
		for _, el := range row[:len(row)-1] {
			fmt.Print(el, "\t|\t")
		}
		fmt.Println(row[len(row)-1])
	}
}

//FindEl находит элемент в матрице и возвращает его координаты
//если он есть в матрице
func FindEl(m [][]int, elem int) ([2]int, error) {
	res := [2]int{-1, -1}

	for i := range m {
		for j := range m[i] {
			if m[i][j] == elem {
				res[0] = i
				res[1] = j
				break
			}
		}
	}
	if res[0] == -1 {
		return res, errors.New("not Found")
	}
	return res, nil

}

//RandFillMatrix заполняет матрицу в случайном порядке
func RandFillMatrix(m *matrix) {
	if m.N == 0 {
		return
	}

	places := make([][2]int, m.N*m.M)

	for i := range m.Data {
		for j := range m.Data[i] {
			places[i*m.M+j][0] = i
			places[i*m.M+j][1] = j
		}
	}

	pI := rand.Intn(len(places))
	places = append(places[:pI], places[pI+1:]...)

	counter := 1
	for len(places) != 0 {

		pI = rand.Intn(len(places))
		x := places[pI][0]
		y := places[pI][1]

		m.Data[x][y] = counter

		places = append(places[:pI], places[pI+1:]...)
		counter++
	}
}
