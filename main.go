// Knn project main.go
package main

import (
	"bufio"
	"fmt"
	"knn/knn"
	"os"
	"runtime"
	"strconv"
	"strings"
)

//var testFile string = "testing-data-small.txt"
var testFile string = "training-data.txt"
var dataSet *Data = newDataSet()
var cleanDataSet *Data = newDataSet()

func main() {
	runtime.GOMAXPROCS(4)

	file, err := os.Open(testFile)
	checkError(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		data := strings.Split(scanner.Text(), " ")
		cuis, err := strconv.Atoi(data[0])
		checkError(err)

		dataSet.addRecipe(NewRecipe(cuis, data[1:]))
	}
	dataSet.calculateMetaData()

	cleanData()
}

func RunKNN() {
	knn.RunKNN(8, cleanDataSet, 2, WeightedJacardDistance, 7)
}

func cleanData() {
	tooFewIngr := 0
	badIngr := 0
	badIngredients := make(map[string]bool)

	ingrInfo := dataSet.ingredientInfo

	for ingr, info := range ingrInfo {
		if info.count < 2 {
			badIngredients[ingr] = true
		}
	}

	file, err := os.Open(testFile)
	checkError(err)
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {

		data := strings.Split(scanner.Text(), " ")
		ingredients := data[1:]
		cuis, err := strconv.Atoi(data[0])
		checkError(err)

		recp := NewRecipe2(cuis)

		for _, ingr := range ingredients {
			if found, _ := badIngredients[ingr]; !found {
				recp.AddIngredient(ingr)
			} else {
				badIngr++
			}
		}
		if len(recp.Ingredients) > 2 {
			cleanDataSet.addRecipe(recp)
		} else {
			tooFewIngr++
		}

	}

	fmt.Println("Size of original dataset: ", len(dataSet.recipes))
	fmt.Println("Size of cleaned dataset : ", len(cleanDataSet.recipes))
	fmt.Println("Number of recipes below cutoff: ", tooFewIngr)
	fmt.Println("Number of unique ingredients removed: ", len(badIngredients))
	fmt.Println("Number of ingredients removed: ", badIngr)
	i := 0
	for key, _ := range badIngredients {
		fmt.Print("\t", key)
		i++
		if i%7 == 6 {
			fmt.Print("\n")
		}
	}

	cleanDataSet.calculateMetaData()
	RunKNN()
}

func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

func WeightedJacardDistance(v1 interface{}, v2 interface{}) float64 {

	recp1 := v1.(*Recipe)
	recp2 := v2.(*Recipe)

	intersection := 0.0
	union := 0.0

	for name, _ := range recp1.getUniqueIngredients() {
		if recp2.hasIngredient(name) {
			intersection += (1.0 * dataSet.getMaxDiffProb(name))
		}
		union += (1.0 * dataSet.getMaxDiffProb(name))
		//fmt.Println("\tDiff: ", dataSet.getMaxDiffProb(name))
	}

	for name, _ := range recp2.getUniqueIngredients() {
		if !recp1.hasIngredient(name) {
			union += (1.0 * dataSet.getMaxDiffProb(name))
		}
	}

	return 1.0 - (float64(intersection) / float64(union))
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
