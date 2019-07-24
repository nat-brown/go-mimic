package class

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gobuffalo/buffalo"
	"github.com/nat-brown/go-mimic/example/dependencies/client"
)

var tokenLock = sync.RWMutex{}
var _token string

const (
	numTries      = 4
	waitIncrement = time.Second

	urlFormat = "%s/auth/%s"
)

func token(ctx buffalo.Context) (t string) {
	tokenLock.RLock() // Lock in case this is called before auth finishes.
	t = _token
	tokenLock.Unlock()
	return t
}

func authorize(ctx buffalo.Context) error {
	l := ctx.Logger()
	reqBody := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{
		Username: "class",
		Password: "class_password",
	}
	body := new(bytes.Buffer)
	err := json.NewEncoder(body).Encode(reqBody)
	if err != nil {
		l.Errorf("error marshaling body for class authorize: %v", err)
		return err
	}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf(urlFormat, client.BaseURL, "token"), body)
	if err != nil {
		l.Errorf("error making request for class authorize: %v", err)
		return err
	}
	c := client.New()
	resp, err := c.Do(req)
	if err != nil {
		l.Errorf("error sending request for class authorize: %v", err)
		return err // Don't break; try again
	}
	defer resp.Body.Close()
	token := struct {
		Token string `json:"token"`
	}{}
	err = json.NewDecoder(resp.Body).Decode(&token)
	if err != nil {
		l.Errorf("error sending request for class authorize: %v", err)
		return err // Don't break; try again
	}
	_token = token.Token
	return nil
}

// Authorize initializes the class's auth token
func Authorize(ctx buffalo.Context) {
	l := ctx.Logger()
	l.Debug("Authorizing class")
	go func() {
		tokenLock.Lock()
		defer tokenLock.Unlock()
		for i := 0; i < numTries; i++ {
			if _token != "" {
				l.Debug("Class already authorized, skipping call")
				break
			}
			err := authorize(ctx)
			if err == nil {
				return
			}
			time.Sleep(time.Duration(i) * waitIncrement)
		}
	}()
}
