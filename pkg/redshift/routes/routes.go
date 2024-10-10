package routes

import (
	"net/http"

	"github.com/grafana/grafana-aws-sdk/pkg/sql/routes"
	"github.com/grafana/redshift-datasource/pkg/redshift"
	"github.com/grafana/sqlds/v4"
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

func (r *RedshiftResourceHandler) clusters(rw http.ResponseWriter, req *http.Request) {
	clusters, err := r.redshift.Clusters(req.Context(), sqlds.Options{})
	routes.SendResources(rw, clusters, err)
}

func (r *RedshiftResourceHandler) workgroups(rw http.ResponseWriter, req *http.Request) {
	workgroups, err := r.redshift.Workgroups(req.Context(), sqlds.Options{})
	routes.SendResources(rw, workgroups, err)
}

func (r *RedshiftResourceHandler) Routes() map[string]func(http.ResponseWriter, *http.Request) {
	routes := r.DefaultRoutes()
	routes["/secrets"] = r.secrets
	routes["/secret"] = r.secret
	routes["/clusters"] = r.clusters
	routes["/workgroups"] = r.workgroups
	return routes
}
