package main

import (
	"fmt"
	"io"
	"learn-1/game"
	"log"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixMilli())

	var depth, N, M int

	file, err := os.Open("matrix.txt")
	defer file.Close()

	if err != nil {
		log.Fatalln(err.Error())
	}
	textB, _ := io.ReadAll(file)
	text := string(textB)
	symbols := strings.Split(text, "\n")
	N, Nerr := strconv.Atoi(symbols[0])
	M, Merr := strconv.Atoi(symbols[1])
	depth, Derr := strconv.Atoi(symbols[2])

	if Nerr != nil || Merr != nil || Derr != nil {
		log.Fatalln("Неверный формат входных данных\nN M\nDepth\nmatrix...")
	}

	if len(symbols)-3 != N {
		log.Fatal("Матрица введена неверно")
	}
	for _, line := range symbols[3:] {
		if len(strings.Split(line, " ")) != M {
			log.Fatal("Матрица введена неверно")
		}
	}

	matrix := game.NewMatrix(N, M)
	symbols = symbols[3:]

	for i := 0; i < N; i++ {
		arr := strings.Split(symbols[i], " ")
		for j := 0; j < M; j++ {
			num, readErr := strconv.Atoi(arr[j])
			if readErr != nil {
				log.Fatalln(readErr.Error())
			}
			matrix.Data[i][j] = num
		}
	}

	fmt.Println("Исходная матрица:")
	matrix.PrintMatrix()
	fmt.Println("Глубина поиска: ", depth)

	startPoint := game.NewStateFromMatrix(matrix)

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
