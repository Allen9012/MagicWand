package log

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
)

var KVLen = 2

type verboseModule map[string]int32
type logFilter []string
type logExtraResource map[string]interface{}

func (f *logFilter) String() string {
	return fmt.Sprint(*f)
}

// Set sets the value of the named command-line flag.
// format: -log.filter key1,key2
func (f *logFilter) Set(value string) error {
	for _, i := range strings.Split(value, ",") {
		*f = append(*f, strings.TrimSpace(i))
	}
	return nil
}

func (m verboseModule) String() string {
	// FIXME strings.Builder
	var buf bytes.Buffer
	for k, v := range m {
		buf.WriteString(k)
		buf.WriteString(strconv.FormatInt(int64(v), 10))
		buf.WriteString(",")
	}
	return buf.String()
}

// Set sets the value of the named command-line flag.
// format: -log.module file=1,file2=2
func (m verboseModule) Set(value string) error {
	for _, i := range strings.Split(value, ",") {
		kv := strings.Split(i, "=")
		if len(kv) == 2 {
			if v, err := strconv.ParseInt(kv[1], 10, 64); err == nil {
				m[strings.TrimSpace(kv[0])] = int32(v)
			}
		}
	}
	return nil
}

// Set sets the value of the named command-line flag.
// format: -log.module file=1,file2=2
func (m logExtraResource) Set(value string) error {
	for _, i := range strings.Split(value, ",") {
		kv := strings.Split(i, "=")
		if len(kv) == KVLen {
			if strings.HasPrefix(kv[1], "$") {
				envKey := strings.TrimSpace(strings.TrimPrefix(kv[1], "$"))
				envValue := os.Getenv(envKey)
				if len(envValue) != 0 {
					m[strings.TrimSpace(kv[0])] = envValue
				}
			} else if len(kv[1]) != 0 {
				m[strings.TrimSpace(kv[0])] = kv[1]
			}
		}
	}
	return nil
}

func (m logExtraResource) String() string {
	// FIXME strings.Builder
	var buf bytes.Buffer
	for k, v := range m {
		buf.WriteString(k)
		buf.WriteString(fmt.Sprintf("%v", v))
		buf.WriteString(",")
	}
	return buf.String()
}
