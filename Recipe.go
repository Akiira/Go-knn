// Recipe
package main

type Recipe struct {
	Cuisine     int
	Ingredients []string
	UniqueID    int
}

var count int = 0

func NewRecipe(cuisine int, ingredients []string) (recp *Recipe) {
	recp = new(Recipe)
	recp.Cuisine = cuisine
	recp.Ingredients = ingredients
	recp.UniqueID = count
	count++

	return recp
}

func (recp *Recipe) hasIngredient(testIngredient string) bool {
	for _, ingr := range recp.Ingredients {
		if ingr == testIngredient {
			return true
		}
	}
	return false
}

func (recp *Recipe) getIngredientCount(ingredientName string) (count int, found bool) {

	for _, ingr := range recp.Ingredients {
		if ingr == ingredientName {
			found = true
			count++
		}
	}

	return count, found
}

func (recp *Recipe) getUniqueIngredients() (ingredients map[string]bool) {
	ingredients = make(map[string]bool)

	for _, ingr := range recp.getIngredients() {
		ingredients[ingr] = true
	}

	return ingredients
}

func (recp *Recipe) getIngredients() []string {
	return recp.Ingredients
}
