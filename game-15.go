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

	pre := startPoint
	for steps > 0 {
		moves := startPoint.GenChilds()

		rand.Shuffle(len(moves), func(i, j int) {
			moves[i], moves[j] = moves[j], moves[i]
		})
		rand.Shuffle(len(moves), func(i, j int) {
			moves[i], moves[j] = moves[j], moves[i]
		})
		rand.Shuffle(len(moves), func(i, j int) {
			moves[i], moves[j] = moves[j], moves[i]
		})
		rand.Shuffle(len(moves), func(i, j int) {
			moves[i], moves[j] = moves[j], moves[i]
		})
		for _, ch := range moves {

			if ch.Hash() == pre.Hash() {
				continue
			}

			pre = startPoint
			startPoint = ch
			break
		}
		steps--
	}

	startPoint.PrintMatrix()
	fmt.Print("\n\n")

	if !startPoint.CheckWinnable() {
		fmt.Println("Решения нет")
		return
	}

	g := game.NewGraph(startPoint)

	t := time.Now()

	DeepPath := g.DeepSearch(depth)
	t2 := time.Since(t)
	fmt.Println("Путь от поиска в глубину", len(DeepPath), t2)

	t = time.Now()
	WidthPath := g.WidthSearch()
	t3 := time.Since(t)
	fmt.Println("Путь от поиска в ширину", len(WidthPath), t3)

}
