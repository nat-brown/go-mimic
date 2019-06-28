package class

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/nat-brown/go-mimic/example/client"
)

const urlFormat = "%s/class/%s"

// Client is a client for the class api.
type Client struct {
	Logger logrus.FieldLogger
	Base   *http.Client
}

// GetRecipes retrieves all recipes.
func (c Client) GetRecipes() (rs Recipes, err error) {
	url := fmt.Sprintf(urlFormat, client.BaseURL, "recipes")
	c.Logger.Debugf("Class.GetRecipes() sending request to %s", url)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	resp, err := c.Base.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("response from class recipes get was %d", resp.StatusCode)
	}

	err = json.NewDecoder(resp.Body).Decode(&rs)
	if err != nil {
		return nil, fmt.Errorf("reading response for get request for class recipes errored: %v", err)
	}
	return rs, nil
}
