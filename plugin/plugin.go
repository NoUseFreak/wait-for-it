package plugin

import (
	"strings"
	"os"
)

func ParseArguments() map[string]string {
	parameters := map[string]string{}
	for _, pair := range strings.Split(os.Args[1], ",") {
		pair := strings.Split(pair, "=")
		parameters[pair[0]] = pair[1]
	}

	return parameters
}