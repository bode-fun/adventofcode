package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
)

func main() {
	filePointer, err := os.Open("./input")
	if err != nil {
		log.Fatalln(err)
	}

	defer filePointer.Close()

	caloriesTxtChannel := make(chan string)
	go readCaloriesFromFile(filePointer, caloriesTxtChannel)

	caloriesListChannel := make(chan []int)
	go accumulateCalories(caloriesTxtChannel, caloriesListChannel)

	maxCalories := calculateMaxCalories(caloriesListChannel)

	fmt.Printf("The elf with the max calories has %d calories\n", maxCalories)
}

func readCaloriesFromFile(filePointer *os.File, caloriesChannel chan string) {
	defer close(caloriesChannel)

	fileScanner := bufio.NewScanner(filePointer)
	fileScanner.Split(bufio.ScanLines)

	for fileScanner.Scan() {
		caloriesTxt := fileScanner.Text()
		caloriesChannel <- caloriesTxt
	}
}

func accumulateCalories(caloriesChannel chan string, caloriesListChannel chan []int) {
	defer close(caloriesListChannel)

	caloriesList := make([]int, 0)

	for caloriesTxt := range caloriesChannel {
		calories, err := strconv.Atoi(caloriesTxt)

		if caloriesTxt == "" {
			caloriesListChannel <- caloriesList
			caloriesList = make([]int, 0)
		} else {
			if err != nil {
				panic(err)
			}

			caloriesList = append(caloriesList, calories)
		}
	}
}

func calculateMaxCalories(caloriesListChannel chan []int) int {
	var maxCaloriesSum int

	for caloriesList := range caloriesListChannel {
		caloriesSum := sumCaloriesPerElf(caloriesList)
		if caloriesSum > maxCaloriesSum {
			maxCaloriesSum = caloriesSum
		}
	}

	return maxCaloriesSum
}

func sumCaloriesPerElf(caloriesList []int) int {
	var caloriesSum int

	for _, calories := range caloriesList {
		caloriesSum += calories
	}

	return caloriesSum
}
