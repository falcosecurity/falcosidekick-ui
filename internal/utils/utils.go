// SPDX-License-Identifier: Apache-2.0
/*
Copyright (C) 2023 The Falco Authors.
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at
    http://www.apache.org/licenses/LICENSE-2.0
Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package utils

import (
	"flag"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/falcosecurity/falcosidekick-ui/configuration"
)

const (
	extractNumber = "^[0-9]+"
	extractUnity  = "[a-z-A-Z]+$"
	trimPrefix    = "(?i)^\\d{2}:\\d{2}:\\d{2}\\.\\d{9}\\:\\ (Debug|Info|Informational|Notice|Warning|Error|Critical|Alert|Emergency)"
)

const (
	debugLog = iota
	infoLog
	warningLog
	errorLog
	fatalLog
)

var regExtractNumber, regExtractUnity, regTrimPrefix *regexp.Regexp

func init() {
	regExtractNumber, _ = regexp.Compile(extractNumber)
	regExtractUnity, _ = regexp.Compile(extractUnity)
	regTrimPrefix, _ = regexp.Compile(trimPrefix)
}

func CheckErr(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}

func WriteLog(level, message string) {
	config := configuration.GetConfiguration()
	if GetPriortiyInt(level) < GetPriortiyInt(config.LogLevel) {
		return
	}
	var prefix string
	switch level {
	case "fatal":
		prefix = "[ERROR]:"
		log.Fatalf("%v %v\n", prefix, message)
	case "info":
		prefix = "[INFO] :"
	case "warning":
		prefix = "[WARN] :"
	case "error":
		prefix = "[ERROR]:"
	}
	log.Printf("%v %v\n", prefix, message)
}

func ConvertToSeconds(s string) int {
	n, err := strconv.Atoi(regExtractNumber.FindString(s))
	if err != nil {
		return 0
	}
	u := regExtractUnity.FindString(s)
	switch u {
	case "s", "second", "seconds":
		return n
	case "m", "min", "minute", "minutes":
		return n * 60
	case "h", "hour", "hours":
		return n * 60 * 60
	case "d", "day", "days":
		return n * 24 * 60 * 60
	case "w", "week", "weeks":
		return n * 7 * 24 * 60 * 60
	case "M", "month", "months":
		return n * 30 * 24 * 60 * 60
	case "y", "year", "years":
		return n * 365 * 24 * 60 * 60
	default:
		o, err := strconv.Atoi(s)
		if err != nil {
			WriteLog("fatal", "invalid TTL")
		}
		return o
	}
}

func RemoveDuplicate(input []string) []string {
	allKeys := make(map[string]bool)
	singleKeys := []string{}
	for _, i := range input {
		if _, value := allKeys[i]; !value {
			allKeys[i] = true
			singleKeys = append(singleKeys, i)
		}
	}
	sort.Strings(singleKeys)
	return singleKeys
}

func GetStringFlagOrEnvParam(flagString string, envVar string, defaultValue string, usage string) *string {
	envvar, present := os.LookupEnv(envVar)
	if present {
		defaultValue = envvar
	}
	return flag.String(flagString, defaultValue, usage)
}

func GetBoolFlagOrEnvParam(flagString string, envVar string, defaultValue bool, usage string) *bool {
	envvar, present := os.LookupEnv(envVar)
	if present {
		if envvar == "true" {
			defaultValue = true
		}
	}
	return flag.Bool(flagString, defaultValue, usage)
}

func GetIntFlagOrEnvParam(flagString string, envVar string, defaultValue int, usage string) *int {
	envvar, present := os.LookupEnv(envVar)
	if present {
		val, err := strconv.Atoi(envvar)
		if err == nil {
			defaultValue = val
		}
	}
	return flag.Int(flagString, defaultValue, usage)
}

func GetPriortiyInt(prio string) int {
	switch prio {
	case "debug":
		return debugLog
	case "info":
		return infoLog
	case "warning":
		return warningLog
	case "error":
		return errorLog
	case "fatal":
		return fatalLog
	case "default":
		return debugLog
	default:
		return -1
	}
}

func Escape(s string) string {
	s = strings.ReplaceAll(s, "-", `\-`)
	s = strings.ReplaceAll(s, "/", "\\/")
	s = strings.ReplaceAll(s, ".", "\\.")
	return s
}

func UnEscape(s string) string {
	s = strings.ReplaceAll(s, `\-`, "-")
	s = strings.ReplaceAll(s, `\\/`, "/")
	return s
}

func TrimPrefix(s string) string {
	return regTrimPrefix.ReplaceAllString(s, "")
}
