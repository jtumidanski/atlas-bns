package name

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"log"
	"net/http"
)

type responseList struct {
	Data []responseData `json:"data"`
}

type response struct {
	Data responseData `json:"data"`
}

type responseData struct {
	ID         string             `json:"id"`
	Type       string             `json:"type"`
	Attributes responseAttributes `json:"attributes"`
}

type responseAttributes struct {
	Name string `json:"name"`
}

func GetNames(l *log.Logger) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		vs := GetCache().GetNames()

		var result responseList
		for _, v := range vs {
			vr := responseData{
				ID:         v,
				Type:       "com.atlas.cos.rest.attribute.BlockedNameAttributes",
				Attributes: responseAttributes{Name: v},
			}
			result.Data = append(result.Data, vr)
		}
		rw.WriteHeader(http.StatusOK)
		err := toJSON(result, rw)
		if err != nil {
			l.Print(err.Error())
		}
	}
}

func GetName(l *log.Logger) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		name := mux.Vars(r)["name"]

		var result responseList
		if val, ok := GetCache().GetName(name); ok {
			vr := responseData{
				ID:         val,
				Type:       "com.atlas.cos.rest.attribute.BlockedNameAttributes",
				Attributes: responseAttributes{Name: val},
			}
			result.Data = append(result.Data, vr)
		} else {
			result.Data = make([]responseData, 0)
		}
		rw.WriteHeader(http.StatusOK)
		err := toJSON(result, rw)
		if err != nil {
			l.Print(err.Error())
		}
	}
}

// toJSON serializes the given interface into a string based JSON format
func toJSON(i interface{}, w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(i)
}
