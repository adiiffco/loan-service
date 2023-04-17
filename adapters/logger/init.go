package logger

import (
	"context"
	"fmt"
	"strings"

	"loanapp/common"

	"github.com/sirupsen/logrus"
	"golang.org/x/exp/maps"
)

type Log struct {
	Tag string
}

func Initialize() {
	fmt.Println("Initializing logger")
	logrus.SetReportCaller(true)
	formatter := &logrus.TextFormatter{
		TimestampFormat:        "02-01-2006 15:04:05", // the "time" field configuratiom
		FullTimestamp:          true,
		DisableLevelTruncation: true, // log level field configuration
	}
	logrus.SetFormatter(formatter)
}

func (l *Log) LogData(
	ctx context.Context,
	level logrus.Level,
	message string,
	fields logrus.Fields,
	identifier ...string,
) {
	if fields == nil {
		fields = make(logrus.Fields)
	}
	requestID := ctx.Value("request_id")
	maps.Copy(fields, logrus.Fields{
		"request_id": requestID,
	})

	if level == logrus.ErrorLevel {
		callerDetails := common.GetCallerMethodDetails(2)
		if callerDetails != nil {
			callerMap := logrus.Fields{
				"method": callerDetails.Method,
				"line":   callerDetails.Line,
			}
			maps.Copy(fields, callerMap)
		}
	}

	identifiers := strings.Join(identifier, "-")
	identifiers = fmt.Sprintf("%s|%s", l.Tag, identifiers)
	message = fmt.Sprintf("%s: %s", identifiers, message)
	logrus.WithContext(ctx).WithFields(fields).Log(level, message)
}
