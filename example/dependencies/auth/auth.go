package auth

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"net/http"

	"github.com/gobuffalo/buffalo"
	"github.com/gobuffalo/buffalo/render"
)

// Permisions is a map of token to user.
var permissions = map[string]string{}

// HandleAuthentication creates a user token for authentication.
func HandleAuthentication(ctx buffalo.Context) error {
	l := ctx.Logger()
	l.Debug("HandleAuthentication called")
	req := struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}{}
	err := json.NewDecoder(ctx.Request().Body).Decode(&req)
	if err != nil {
		l.Error("error decoding handle authentication request body", err)
		return ctx.Render(http.StatusBadRequest, render.String("Invalid request format"))
	}
	if req.Username == "" || req.Password == "" {
		l.Debug("Username or password was not detected.")
		return ctx.Render(http.StatusBadRequest, render.String("Username and Password required."))
	}

	combined := []byte(req.Username + req.Password)
	hasher := md5.New()
	hash := hex.EncodeToString(hasher.Sum(combined))
	l.Debugf("Hash was %s", hash)
	permissions[hash] = req.Username
	resp := struct {
		Token string `json:"token"`
	}{
		Token: hash,
	}
	return ctx.Render(http.StatusOK, render.JSON(resp))
}

// HandleVerification returns the user associated with a given token.
func HandleVerification(ctx buffalo.Context) error {
	l := ctx.Logger()
	l.Debug("HandleVerification called")
	req := struct {
		Token string `json:"token"`
	}{}
	err := json.NewDecoder(ctx.Request().Body).Decode(&req)
	if err != nil {
		l.Error("error decoding handle verification request body", err)
		return ctx.Render(http.StatusBadRequest, render.String("Invalid request format"))
	}
	user := permissions[req.Token]
	if user == "" {
		return ctx.Render(http.StatusNotFound, nil)
	}
	resp := struct {
		Username string `json:"username"`
	}{
		Username: user,
	}
	return ctx.Render(http.StatusOK, render.JSON(resp))
}
