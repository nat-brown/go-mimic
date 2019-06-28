package app

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/sirupsen/logrus"

	"github.com/nat-brown/go-mimic/example/client/class"
)

// RecipeHandler handles recipe get requests.
type RecipeHandler struct {
	Logger logrus.FieldLogger
	Client *http.Client
}

// Handle recipe get requests
func (rh RecipeHandler) Handle(w http.ResponseWriter, r *http.Request) {
	rh.Logger.Debugf("get request for recipes")

	accept := r.Header.Get(HeaderAccept)
	switch accept {
	case AcceptJSON:
		rs, err := rh.getRecipes()
		if err != nil {
			err.Handle(w, rh.Logger)
			return
		}

		w.Header().Set(HeaderContentType, AcceptJSON)
		if err := json.NewEncoder(w).Encode(rs); err != nil {
			ResponseError{
				InnerError:   fmt.Sprintf("error encoding recipes response %v", err),
				ResponseCode: http.StatusInternalServerError,
			}.Handle(w, rh.Logger)
			return
		}
		return
	case AcceptText, AcceptUnspecified:
		rs, err := rh.getRecipes()
		if err != nil {
			err.Handle(w, rh.Logger)
			return
		}

		w.Header().Set(HeaderContentType, AcceptText)
		w.Write([]byte(rs.String()))
	default:
		ResponseError{
			InnerError:   fmt.Sprintf("accept value was %s", accept),
			ReturnError:  "Accept header invalid: use application/json or text/plain.",
			ResponseCode: http.StatusNotAcceptable,
		}.Handle(w, rh.Logger)
	}
}

func (rh RecipeHandler) getRecipes() (class.Recipes, ErrorHandler) {
	pc := class.Client{
		Logger: rh.Logger,
		Base:   rh.Client,
	}

	rs, err := pc.GetRecipes()
	if err != nil {
		return nil, NoMessageError{
			InnerError:   fmt.Sprintf("get request for class recipes errored: %v", err),
			ResponseCode: http.StatusInternalServerError,
		}
	}

	return rs, nil
}
