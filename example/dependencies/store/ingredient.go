package store

import (
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"

	"github.com/nat-brown/go-mimic/example/dependencies/database"
)

var data = []Ingredient{}

// Handle store query calls.
func Handle(ctx buffalo.Context) error {
	return ctx.Render(http.StatusOK, render.JSON(data))
}

// Ingredient instance.
type Ingredient struct {
	Name        string   `json:"name"`
	Permissions []string `json:"permissions"`
}

func init() {
	err := database.Unmarshal(&data, "store", "store_data.json")
	if err != nil {
		panic(err)
	}
}
