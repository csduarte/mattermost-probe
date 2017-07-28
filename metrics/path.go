package metrics

import (
	"regexp"
	"strings"

	"github.com/mattermost/platform/model"
)

// rawSubstitutions holds the raw strings that will be compiled at run time
var rawSubstitutions = map[string]string{
	//Channel Routes
	"/channels/cid/":       "/channels/[a-z0-9]{26}/",       //Channel ID
	"/channels/name/cname": "/channels/name/[A-Za-z0-9_-]+", //Get Channel By Name

	//Team Routes
	"/teams/tid/": "/teams/[a-z0-9]{26}/", //Team ID
}

// Subtitutions holds the compiled regex
var Subtitutions = map[string]*regexp.Regexp{}

func init() {
	for k, v := range rawSubstitutions {
		Subtitutions[k] = regexp.MustCompile(v)
	}
}

// TokenizePath organize and clean up path name based off known formats in Subtitutions
func TokenizePath(path string) string {
	result := strings.TrimPrefix(path, model.API_URL_SUFFIX)
	for sub, reg := range Subtitutions {
		result = reg.ReplaceAllString(result, sub)
	}
	return result
}
