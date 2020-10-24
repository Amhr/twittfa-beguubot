package epimetheus

import (
	"github.com/cafebazaar/epimetheus"
	"github.com/spf13/viper"
)

func NewEpimetheus(v *viper.Viper) *epimetheus.Epimetheus {
	epimetheusServer := epimetheus.NewEpimetheusServer(v)
	go epimetheusServer.Listen()
	return epimetheusServer
}
