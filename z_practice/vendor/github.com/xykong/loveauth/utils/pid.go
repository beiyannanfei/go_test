package utils

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"os"
)

func SavePid() {

	file, err := os.Create("scripts/loveauth.pid")
	if err != nil {

		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Warn("savepid err")

		return
	}

	defer file.Close()

	file.WriteString(fmt.Sprintf("%d", os.Getpid()))
}
