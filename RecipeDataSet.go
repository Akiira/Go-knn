// RecipeDataSet
package main

import (
	"math"
)

type Data struct {
	recipes         []*Recipe
	ingredientInfo  map[string]*ingredientEntry
	ingredientCount int
}

func newDataSet() *Data {
	var ds Data
	ds.ingredientInfo = make(map[string]*ingredientEntry)
	ds.recipes = make([]*Recipe, 0)

	return &ds
}

func (ds *Data) addRecipe(newRecipe *Recipe) {
	for _, ingr := range newRecipe.getIngredients() {

		if val, found := ds.ingredientInfo[ingr]; found {
			val.incrementCount(newRecipe.Cuisine)
		} else {
			ds.ingredientInfo[ingr] = newIngredientEntry(ingr)
		}
		ds.ingredientCount++
	}
	ds.recipes = append(ds.recipes, newRecipe)
}

func (ds *Data) calculateMetaData() {
	for _, val := range ds.ingredientInfo {
		val.calculateEntropy()
	}
}

//Bonus just uses one ingredients entropy
func (ds *Data) getEntropyWeight(recpipeID int) float64 {
	var entropy float64 = 2

	for _, ingr := range ds.recipes[recpipeID].getIngredients() {
		if ds.ingredientInfo[ingr].entropy < entropy {
			entropy = ds.ingredientInfo[ingr].entropy
		}
	}
	return 1.5 - entropy
}

func (ds *Data) Weight(index int) float64 {
	return ds.getEntropyWeight(index)
}

func (ds *Data) GetEntry(index int) interface{} {
	return ds.recipes[index]
}

func (ds *Data) Size() int {
	return len(ds.recipes)
}

func (ds *Data) GetClass(index int) int {
	return ds.recipes[index].Cuisine
}

type ingredientEntry struct {
	ingredient    string
	count         int
	cuisineCounts []int
	entropy       float64
}

func newIngredientEntry(name string) *ingredientEntry {
	ie := ingredientEntry{ingredient: name}
	ie.cuisineCounts = make([]int, 7)

	return &ie
}

func (ie *ingredientEntry) incrementCount(cuisine int) {
	ie.count++
	ie.cuisineCounts[cuisine-1]++
}

func (ie *ingredientEntry) calculateEntropy() float64 {
	ie.entropy = 0

	for _, cousineCount := range ie.cuisineCounts {
		var pr float64 = float64(cousineCount) / float64(ie.count)

		ie.entropy -= pr * math.Log2(pr)
	}

	return ie.entropy
}
