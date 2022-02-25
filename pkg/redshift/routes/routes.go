package routes

import (
	"net/http"

	"github.com/grafana/grafana-aws-sdk/pkg/sql/routes"
	"github.com/grafana/redshift-datasource/pkg/redshift"
	"github.com/grafana/sqlds/v2"
)

type RedshiftResourceHandler struct {
	routes.ResourceHandler
	redshift redshift.RedshiftDatasourceIface
}

func New(api redshift.RedshiftDatasourceIface) *RedshiftResourceHandler {
	return &RedshiftResourceHandler{routes.ResourceHandler{API: api}, api}
}

func (r *RedshiftResourceHandler) secrets(rw http.ResponseWriter, req *http.Request) {
	secrets, err := r.redshift.Secrets(req.Context(), sqlds.Options{})
	routes.SendResources(rw, secrets, err)
}

func (r *RedshiftResourceHandler) secret(rw http.ResponseWriter, req *http.Request) {
	reqBody, err := routes.ParseBody(req.Body)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		routes.Write(rw, []byte(err.Error()))
		return
	}
	secret, err := r.redshift.Secret(req.Context(), reqBody)
	routes.SendResources(rw, secret, err)
}

func (r *RedshiftResourceHandler) cluster(rw http.ResponseWriter, req *http.Request) {
	reqBody, err := routes.ParseBody(req.Body)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		routes.Write(rw, []byte(err.Error()))
		return
	}
	cluster, err := r.redshift.Cluster(req.Context(), reqBody)
	routes.SendResources(rw, cluster, err)
}

func (r *RedshiftResourceHandler) Routes() map[string]func(http.ResponseWriter, *http.Request) {
	routes := r.DefaultRoutes()
	routes["/secrets"] = r.secrets
	routes["/secret"] = r.secret
	routes["/cluster"] = r.cluster
	return routes
}
