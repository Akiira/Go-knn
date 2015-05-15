/*
* TODO:
* 1. Store more of the meta data rather then recalculating it
 */

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

//Weight just uses one ingredients entropy
func (ds *Data) getEntropyWeight(recpipeID int) float64 {
	var entropy float64 = 2

	for _, ingr := range ds.recipes[recpipeID].getIngredients() {
		if ds.ingredientInfo[ingr].entropy < entropy {
			entropy = ds.ingredientInfo[ingr].entropy
		}
	}
	return 1.5 - entropy
}

//func (ds *Data) getMaxPr(recipeID int) float64 {
//	min := 99999.0
//	max := -999999.0
//	for _, ingrName := range recp.getIngredients() {
//		ingr := ds.ingredientInfo[ingrName]
//		for _, count := range ingr.cuisineCounts {

//		}

//	}

//    for (int i = 1; i <= 8; i++) {
//        double current = getProbabilityForCuisine(i + 1);
//        min = min < current ? min : current;
//        max = max > current ? max : current;
//    }

//    return max - min;
//}

//func (ds *Data) getMaxCountDifference(recipeID int) float64 {
//	recp := ds.recipes[recipeID]
//	min := 100000
//	max := 0

//	for _, ingrName := range recp.getIngredients() {
//		ingr := ds.ingredientInfo[ingrName]

//		for _, count := range ingr.cuisineCounts {
//			if count > 0 && count < min {
//				min = count
//			} else if count > max {
//				max = count
//			}
//		}
//	}

//	return 0.5 + (float64(max-min) / float64(max))
//}

func (ds *Data) getMaxDiffProb(ingrName string) float64 {
	ie := ds.ingredientInfo[ingrName]
	min := 99999.0
	max := -999999.0
	for i := 1; i <= 7; i++ {
		current := ie.getProbability(i)

		if current < min {
			min = current
		}
		if current > max {
			max = current
		}
	}

	return max - min
}

func (ds *Data) getProbability(ingrName string, cuisine int) float64 {
	ie := ds.ingredientInfo[ingrName]
	return ie.getProbability(cuisine)
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
	ie.cuisineCounts = make([]int, 8) //TODO remove hard coding

	return &ie
}

func (ie *ingredientEntry) incrementCount(cuisine int) {
	ie.count++
	ie.cuisineCounts[cuisine]++
}

func (ie *ingredientEntry) getProbability(cuisine int) float64 {
	return float64(ie.cuisineCounts[cuisine]) / float64(ie.count)
}

func (ie *ingredientEntry) calculateEntropy() float64 {
	ie.entropy = 0

	for _, cousineCount := range ie.cuisineCounts {
		var pr float64 = float64(cousineCount) / float64(ie.count)

		ie.entropy -= pr * math.Log2(pr)
	}

	return ie.entropy
}
