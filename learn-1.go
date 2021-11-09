package main

import (
	"fmt"
	"learn-1/game"
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixMilli())

	var steps, depth int

	matrix := game.NewMatrix(4, 4)
	matrix.Data[0] = []int{1, 2, 3, 4}
	matrix.Data[1] = []int{5, 6, 7, 8}
	matrix.Data[2] = []int{9, 10, 11, 12}
	matrix.Data[3] = []int{13, 14, 15, 0}

	fmt.Print("Количество случайных ходов на матрице: ")
	_, _ = fmt.Scanf("%d", &steps)
	fmt.Print("Глубина поиска: ")
	_, _ = fmt.Scanf("%d", &depth)

	startPoint := game.NewStateFromMatrix(matrix)
	//startPoint := game.NewState()

	for steps > 0 {
		moves := []func() (*game.State, error){startPoint.Up, startPoint.Down, startPoint.Left, startPoint.Right}
		next, _ := moves[rand.Int31()%4]()
		if next == nil {
			continue
		}
		startPoint = next
		steps--
	}

	startPoint.PrintMatrix()
	fmt.Print("\n\n")

	if !startPoint.CheckWinnable() {
		fmt.Println("Решения нет")
		return
	}

	g := game.NewGraph(startPoint)

	path := g.DeepSearch(depth)

	fmt.Println(len(path))

}
