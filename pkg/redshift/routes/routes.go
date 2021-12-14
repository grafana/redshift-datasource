package routes

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/grafana/grafana-plugin-sdk-go/backend/log"
	"github.com/grafana/redshift-datasource/pkg/redshift"
	"github.com/grafana/sqlds/v2"
)

type ResourceHandler struct {
	ds redshift.RedshiftDatasourceIface
}

func New(ds *redshift.RedshiftDatasource) *ResourceHandler {
	return &ResourceHandler{ds: ds}
}

func write(rw http.ResponseWriter, b []byte) {
	_, err := rw.Write(b)
	if err != nil {
		log.DefaultLogger.Error(err.Error())
	}
}

func parseBody(body io.ReadCloser) (sqlds.Options, error) {
	reqBody := sqlds.Options{}
	b, err := ioutil.ReadAll(body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &reqBody)
	if err != nil {
		return nil, err
	}
	return reqBody, nil
}

func sendResponse(res interface{}, err error, rw http.ResponseWriter) {
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		write(rw, []byte(err.Error()))
		return
	}
	bytes, err := json.Marshal(res)
	if err != nil {
		log.DefaultLogger.Error(err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		write(rw, []byte(err.Error()))
		return
	}
	rw.Header().Add("Content-Type", "application/json")
	write(rw, bytes)
}

func (r *ResourceHandler) secrets(rw http.ResponseWriter, req *http.Request) {
	secrets, err := r.ds.Secrets(req.Context(), sqlds.Options{})
	sendResponse(secrets, err, rw)
}

func (r *ResourceHandler) secret(rw http.ResponseWriter, req *http.Request) {
	reqBody, err := parseBody(req.Body)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		write(rw, []byte(err.Error()))
		return
	}
	secrets, err := r.ds.Secret(req.Context(), reqBody)
	sendResponse(secrets, err, rw)
}

func (r *ResourceHandler) Routes() map[string]func(http.ResponseWriter, *http.Request) {
	return map[string]func(http.ResponseWriter, *http.Request){
		"/secrets": r.secrets,
		"/secret":  r.secret,
	}
}
