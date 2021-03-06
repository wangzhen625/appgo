package server

import (
	log "github.com/Sirupsen/logrus"
	"github.com/oxfeeefeee/appgo/redis2"
	"net/http"
)

const (
	zsetNamespace = "metrics"
	totalKey      = "total"

	contextKeyUser = "user"
)

type Metrics struct {
	schema MetricsSchema
	zsets  *redis2.Zsets
}

func newMetrics(schema MetricsSchema) *Metrics {
	zsets := redis2.NewZsets(zsetNamespace)
	return &Metrics{
		schema, zsets,
	}
}

func (m *Metrics) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	next(rw, r)

	var params []redis2.ZsetIncrbyParams
	keys := m.schema.KeysGen(r)
	for k, v := range keys {
		params = append(params, redis2.ZsetIncrbyParams{k, v, 1})
	}
	go func() {
		err := m.zsets.BatchIncrby(params)
		if err != nil {
			log.WithField("params", params).Errorln("BatchIncrby metrics error: ", err)
		}
	}()
}

type DefaultSchema struct {
}

func NewDefaultSchema() *DefaultSchema {
	return &DefaultSchema{}
}

func (s *DefaultSchema) KeysGen(r *http.Request) map[string]string {
	id := GetUserFromToken(r)
	item := r.Method + r.URL.Path
	return map[string]string{
		"raw": item,
		"raw:u" + id.String(): item,
		"raw:users":           id.String(),
	}
}
