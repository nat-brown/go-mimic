package class

import (
	"encoding/json"
	"fmt"
	"strings"
)

// Recipes a list of recipes for custom formatting.
type Recipes []Recipe

// String formatting for a list of recipes.
func (rs Recipes) String() string {
	stringRs := make([]string, len(rs))
	for i, r := range rs {
		stringRs[i] = r.String()
	}
	return strings.Join(stringRs, "\n\n")
}

// Recipe instance.
type Recipe struct {
	Directions  string      `json:"directions"`
	Ingredients Ingredients `json:"ingredients"`
	Title       string      `json:"title"`
}

// String formatting for a recipe.
func (r Recipe) String() string {
	return fmt.Sprintf("%s\n=====\n%s\n%s", r.Title, r.Ingredients, r.Directions)
}

// Ingredients is a list of ingredients for custom printing.
type Ingredients []Ingredient

// String formatting for a list of ingredients.
func (is Ingredients) String() string {
	stringIs := make([]string, len(is))
	for i, ingredient := range is {
		stringIs[i] = ingredient.String()
	}
	return strings.Join(stringIs, "\n")
}

// UnmarshalJSON unmarshals custom list
func (is *Ingredients) UnmarshalJSON(data []byte) error {
	list := []Ingredient{}
	if err := json.Unmarshal(data, &list); err != nil {
		return err
	}
	*is = Ingredients(list)
	return nil
}

// Ingredient for recipe.
type Ingredient struct {
	Name     string  `json:"name,omitempty"`
	Quantity float32 `json:"quantity,omitempty"`
	Unit     string  `json:"unit,omitempty"`
}

// String formatting for an ingredient.
func (i Ingredient) String() (s string) {
	qString, _ := json.Marshal(i.Quantity) // Marshaling should not fail for floats.
	switch {
	case i.Quantity == 0 && i.Unit == "":
		s = i.Name
	case i.Quantity == 0:
		s = fmt.Sprintf("%s of %s", i.Unit, i.Name)
	case i.Unit == "":
		s = fmt.Sprintf("%s %s", qString, i.Name)
	default:
		s = fmt.Sprintf("%s %s of %s", qString, i.Unit, i.Name)
	}
	return fmt.Sprintf("* %s", s)
}
