/**
 * (C) Copyright IBM Corp. 2021.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
package log

import (
	"github.com/sirupsen/logrus"
	"strings"
)

var logger *logrus.Logger

func init() {
	logger = logrus.New()
	SetLogLevel("info")
}
func GetLogger() *logrus.Logger {
	return logger
}

// SetLogger sets the logger instance
// This is useful in testing as the logger can be overridden
// with a test logger
func SetLogger(l *logrus.Logger) {
	logger = l
}
func DebugEnabled() bool {
	return logrus.GetLevel() >= logrus.DebugLevel
}

func InfoEnabled() bool {
	return logrus.GetLevel() >= logrus.InfoLevel
}

func Debug(args ...interface{}) {
	log("debug", args)
}

func Info(args ...interface{}) {
	log("info", args)
}

func Warn(args ...interface{}) {
	log("warn", args)
}

func Error(args ...interface{}) {
	log("error", args)
}

func Fatal(args ...interface{}) {
	log("fatal", args)
}

func Panic(args ...interface{}) {
	log("panic", args)
}
func SetLogLevel(level string) {
	level = strings.ToLower(level)
	switch level {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	case "fatal":
		logger.SetLevel(logrus.FatalLevel)
	case "panic":
		logger.SetLevel(logrus.PanicLevel)
	}
}

func log(level string, args []interface{}) {
	args = append([]interface{}{"AppConfiguration - "}, args...)

	switch level {
	case "debug":
		logger.Debug(args...)
	case "info":
		logger.Info(args...)
	case "warn":
		logger.Warn(args...)
	case "error":
		logger.Error(args...)
	case "fatal":
		logger.Fatal(args...)
	case "panic":
		logger.Panic(args...)
	}
}
