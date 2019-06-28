package class

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/nat-brown/go-mimic/example/client"
	"github.com/nat-brown/go-mimic/example/logger"
)

func makeCall(t *testing.T, filename string) (Recipes, error) {
	l := logger.NewTestLogger(t)
	c := Client{
		Logger: &l,
		Base:   client.NewTestClient("client", "class", "test_data", filename),
	}
	defer l.Hook.Print()
	return c.GetRecipes()
}

func TestClientGetRecipesJSON(t *testing.T) {
	actual, err := makeCall(t, "happy_path.json")
	require.NoError(t, err)

	expected := Recipes([]Recipe{
		{
			Directions: "Spray the pan with 2 cans of oil. Put the dough on the pan. Put the olives on the dough. Put the pan in the oven at 10 degrees. Cook for 2 hours.",
			Ingredients: Ingredients([]Ingredient{
				{
					Name: "Dough from Meijer",
				}, {
					Name:     "crushed up olives",
					Quantity: 10,
				}, {
					Name:     "oil",
					Quantity: 2,
					Unit:     "cans",
				},
			}),
			Title: "Phillip's Pizza",
		},
	})
	assert.Equal(t, expected, actual)
}

func TestBadStatusCodeErrors(t *testing.T) {
	_, err := makeCall(t, "bad_status_code.json")
	require.Error(t, err)
}

func TestBadSchemaErrors(t *testing.T) {
	_, err := makeCall(t, "invalid_schema.json")
	require.Error(t, err)
}
