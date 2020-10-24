package logrus

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"os"
)

func MakeLogrus(v *viper.Viper) {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	if viper.GetBool("dev") == true {
		fmt.Println("Logrus Dev")
		logrus.SetOutput(os.Stdout)
	} else {
		fmt.Println("Logrus Prod")
		logrus.SetLevel(logrus.WarnLevel)
	}
}
