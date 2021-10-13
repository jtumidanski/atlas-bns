package name

import (
	json2 "atlas-bns/json"
	"atlas-bns/rest"
	"github.com/gorilla/mux"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
	"net/http"
)

type dataListContainer struct {
	Data []dataBody `json:"data"`
}

type dataBody struct {
	ID         string     `json:"id"`
	Type       string     `json:"type"`
	Attributes attributes `json:"attributes"`
}

type attributes struct {
	Name string `json:"name"`
}

const (
	getName  = "get_name"
	getNames = "get_names"
)

func InitResource(router *mux.Router, l logrus.FieldLogger) {
	r := router.PathPrefix("/names").Subrouter()
	r.HandleFunc("", registerGetName(l)).Methods(http.MethodGet).Queries("name", "{name}")
	r.HandleFunc("", registerGetNames(l)).Methods(http.MethodGet)
}

func registerGetNames(l logrus.FieldLogger) http.HandlerFunc {
	return rest.RetrieveSpan(getNames, func(span opentracing.Span) http.HandlerFunc {
		return handleGetNames(l)(span)
	})
}

func handleGetNames(l logrus.FieldLogger) func(span opentracing.Span) http.HandlerFunc {
	return func(span opentracing.Span) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			var result dataListContainer
			for _, v := range GetCache().GetNames() {
				result.Data = append(result.Data, getDataBody(v))
			}
			w.WriteHeader(http.StatusOK)
			err := json2.ToJSON(result, w)
			if err != nil {
				l.WithError(err).Errorf("Error writing response")
			}
		}
	}
}

type nameHandler func(name string) http.HandlerFunc

func parseName(l logrus.FieldLogger, next nameHandler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if val, ok := mux.Vars(r)["name"]; ok {
			next(val)(w, r)
		} else {
			l.Errorf("Name was not provided to request.")
			w.WriteHeader(http.StatusBadRequest)
		}
	}
}

func registerGetName(l logrus.FieldLogger) http.HandlerFunc {
	return rest.RetrieveSpan(getName, func(span opentracing.Span) http.HandlerFunc {
		return parseName(l, func(name string) http.HandlerFunc {
			return handleGetName(l)(span)(name)
		})
	})
}

func handleGetName(l logrus.FieldLogger) func(span opentracing.Span) func(name string) http.HandlerFunc {
	return func(span opentracing.Span) func(name string) http.HandlerFunc {
		return func(name string) http.HandlerFunc {
			return func(w http.ResponseWriter, r *http.Request) {
				var result dataListContainer
				if val, ok := GetCache().GetName(name); ok {
					result.Data = append(result.Data, getDataBody(val))
				} else {
					result.Data = make([]dataBody, 0)
				}
				w.WriteHeader(http.StatusOK)
				err := json2.ToJSON(result, w)
				if err != nil {
					l.WithError(err).Errorf("Error writing response")
				}
			}
		}
	}
}

func getDataBody(name string) dataBody {
	return dataBody{
		ID:         name,
		Type:       "com.atlas.cos.rest.attribute.BlockedNameAttributes",
		Attributes: attributes{Name: name},
	}
}
