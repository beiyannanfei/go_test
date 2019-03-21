package log

import (
	"github.com/heirko/go-contrib/logrusHelper"
	_ "github.com/heralight/logrus_mate/hooks/file"
	_ "github.com/heralight/logrus_mate/hooks/filewithformatter"
	_ "github.com/heralight/logrus_mate/hooks/graylog"
	_ "github.com/heralight/logrus_mate/hooks/slack"
	"github.com/sirupsen/logrus"
	"github.com/xykong/loveauth/settings"
)

func InitLogger() {

	var logSetting = settings.Get("logger")

	// Read configuration
	var c = logrusHelper.UnmarshalConfiguration(logSetting) // Unmarshal configuration from Viper
	logrusHelper.SetConfig(logrus.StandardLogger(), c)      // for e.g. apply it to logrus default instance

	// ### End Read Configuration

	//// ### Use logrus as normal
	//logrus.WithFields(logrus.Fields{
	//	"animal": "walrus",
	//}).Error("A walrus appears")
	//
	//Infof(map[string]interface{}{
	//	"aaa": 1,
	//}, "aaa", nil)
}

//// Infof logs a message at level Info on the standard logger.
//func Infof(fields map[string]interface{}, format string, args ...interface{}) {
//	//std.Infof(format, args...)
//
//	logrus.WithFields(fields).Infof(format, args)
//}

var BILogger *logrus.Logger
func InitBILogger() error {
	//if nil == BILogger {
		var logSetting = settings.Get("logger_bi")

		logger := logrus.New()
		// Read configuration
		var c = logrusHelper.UnmarshalConfiguration(logSetting) // Unmarshal configuration from Viper
		err := logrusHelper.SetConfig(logger, c)      // for e.g. apply it to logrus default instance
		if nil != err {
			logrus.WithFields(logrus.Fields{"err": err.Error()})
			return err
		}

		BILogger = logger
		return nil
	//}
	//
	//BILogger.Println("2333")
	//return BILogger
}



