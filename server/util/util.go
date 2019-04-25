package util

import (
	"crypto/sha256"
	"encoding/base64"
	"github.com/mattermost/mattermost-server/model"
	"github.com/pkg/errors"
	"github.com/standup-raven/standup-raven/server/otime"
	"log"
	"net/http"
	"net/http/httputil"
	"regexp"
	"strings"
)

func SplitArgs(s string) ([]string, error) {
	indexes := regexp.MustCompile("\"").FindAllStringIndex(s, -1)
	if len(indexes)%2 != 0 {
		return []string{}, errors.New("quotes not closed")
	}

	indexes = append([][]int{{0, 0}}, indexes...)

	if indexes[len(indexes)-1][1] < len(s) {
		indexes = append(indexes, [][]int{{len(s), 0}}...)
	}

	var args []string
	for i := 0; i < len(indexes)-1; i++ {
		start := indexes[i][1]
		end := Min(len(s), indexes[i+1][0])

		if i%2 == 0 {
			args = append(args, strings.Split(strings.Trim(s[start:end], " "), " ")...)
		} else {
			args = append(args, s[start:end])
		}

	}

	cleanedArgs := make([]string, len(args))
	count := 0

	for _, arg := range args {
		if arg != "" {
			cleanedArgs[count] = arg
			count++
		}
	}

	return cleanedArgs[0:count], nil
}

// Because math.Min is for floats and
// casting to and from floats is dangerous.
func Min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func Max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func SendEphemeralText(msg string) (*model.CommandResponse, *model.AppError) {
	return &model.CommandResponse{
		Type: model.COMMAND_RESPONSE_TYPE_EPHEMERAL,
		Text: msg,
	}, nil
}

// Set Difference: A - B
func Difference(a, b []string) (diff []string) {
	m := make(map[string]bool)

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}
	return
}

func GetCurrentDateString() string {
	return otime.Now().Format("20060102")
}

func GetKeyHash(key string) string {
	hash := sha256.New()
	hash.Write([]byte(key))
	return base64.StdEncoding.EncodeToString(hash.Sum(nil))
}

type LogWriter struct {
	http.ResponseWriter
}

func (w LogWriter) Write(p []byte) (n int, err error) {
	n, err = w.ResponseWriter.Write(p)
	if err != nil {
		log.Printf("Write failed: %v", err)
	}
	return
}

func DumpRequest(r *http.Request) string {
	d, err := httputil.DumpRequest(r, r.Method == http.MethodPost)
	if err != nil {
		d = []byte{}
	}

	return string(d)
}

func ContainsDuplicates(data *[]string) (string, bool) {
	seen := make(map[string]bool, len(*data))

	for _, item := range *data {
		if _, ok := seen[item]; ok {
			return item, true
		} else {
			seen[item] = true
		}
	}

	return "", false
}
