package class

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
	"github.com/nat-brown/go-mimic/example/dependencies/database"
)

var data = []Recipe{}

// HandleRecipes class query calls.
func HandleRecipes(ctx buffalo.Context) error {
	return ctx.Render(http.StatusOK, render.JSON(data))
}

// Recipe instance.
type Recipe struct {
	Directions  string `json:"directions"`
	Ingredients []struct {
		Name     string  `json:"name"`
		Quantity float32 `json:"quantity"`
		Unit     string  `json:"unit"`
	} `json:"ingredients"`
	Title string `json:"title"`
}

func init() {
	err := database.Unmarshal(&data, "class", "recipe_data.json")
	if err != nil {
		panic(err)
	}
}
