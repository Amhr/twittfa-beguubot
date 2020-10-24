package viper

import (
	"fmt"
	"github.com/spf13/viper"
)

func NewViper() (*viper.Viper, error) {

	viper.SetConfigName("epimetheus")

	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()

	if err != nil {
		return nil, fmt.Errorf("error while reading viper $w", err)
	}
	return viper.GetViper(), nil
}
