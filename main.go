// Knn project main.go
package main

import (
	"bufio"
	"fmt"
	"knn/knn"
	"os"
	"strconv"
	"strings"
)

//var testFile string = "testing-data-small.txt"
var testFile string = "training-data.txt"
var dataSet *Data = newDataSet()

func main() {
	fmt.Println("Hello World!")

	file, err := os.Open(testFile)
	checkError(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		data := strings.Split(scanner.Text(), " ")
		cuis, err := strconv.Atoi(data[0])
		checkError(err)

		if len(data) > 2 {
			dataSet.addRecipe(NewRecipe(cuis, data[1:]))
		}
	}
	dataSet.calculateMetaData()

	knn.RunKNN(6, dataSet, 2, JacardDistance, 7)
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func JacardDistance(v1 interface{}, v2 interface{}) float64 {

	recp1 := v1.(*Recipe)
	recp2 := v2.(*Recipe)

	var intersection int = 0
	var union int = 0

	for name, _ := range recp1.getUniqueIngredients() {
		if recp2.hasIngredient(name) {
			intersection++
		}
		union++
	}

	for name, _ := range recp2.getUniqueIngredients() {
		if !recp1.hasIngredient(name) {
			union++
		}
	}

	return 1.0 - (float64(intersection) / float64(union))
}

//TODO Haven't reviewed formula for this
func BagJacardDistance(v1 interface{}, v2 interface{}) float64 {

	recp1 := v1.(*Recipe)
	recp2 := v2.(*Recipe)

	var intersection int = 0
	var union int = 0

	for _, name := range recp1.getIngredients() {
		if recp2.hasIngredient(name) {
			intersection++
		}
		union++
	}

	for _, name := range recp2.getIngredients() {
		if !recp1.hasIngredient(name) {
			union++
		}
	}

	return 1.0 - (float64(intersection) / float64(union))
}
