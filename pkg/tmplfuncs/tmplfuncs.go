package tmplfuncs

import (
	"errors"
	"net"
	"os"
	"os/user"
	"strconv"
	"strings"
)

var Funcs = map[string]interface{}{
	// built-in functions
	"netResolveIPAddr":    net.ResolveIPAddr,
	"osHostname":          os.Hostname,
	"osUserCacheDir":      os.UserCacheDir,
	"osUserConfigDir":     os.UserConfigDir,
	"osUserHomeDir":       os.UserHomeDir,
	"osGetegid":           os.Getegid,
	"osGetenv":            os.Getenv,
	"osGeteuid":           os.Geteuid,
	"osGetgid":            os.Getgid,
	"osGetgroups":         os.Getgroups,
	"osGetpagesize":       os.Getpagesize,
	"osGetpid":            os.Getpid,
	"osGetppid":           os.Getppid,
	"osGetuid":            os.Getuid,
	"osGetwd":             os.Getwd,
	"osTempDir":           os.TempDir,
	"osUserLookupGroup":   user.LookupGroup,
	"osUserLookupGroupId": user.LookupGroupId,
	"osUserCurrent":       user.Current,
	"osUserLookup":        user.Lookup,
	"osUserLookupId":      user.LookupId,
	"stringsContains":     strings.Contains,
	"stringsFields":       strings.Fields,
	"stringsIndex":        strings.Index,
	"stringsLastIndex":    strings.LastIndex,
	"stringsHasPrefix":    strings.HasPrefix,
	"stringsHasSuffix":    strings.HasSuffix,
	"stringsRepeat":       strings.Repeat,
	"stringsReplaceAll":   strings.ReplaceAll,
	"stringsSplit":        strings.Split,
	"stringsSplitN":       strings.SplitN,
	"stringsToLower":      strings.ToLower,
	"stringsToUpper":      strings.ToUpper,
	"stringsTrimPrefix":   strings.TrimPrefix,
	"stringsTrimSpace":    strings.TrimSpace,
	"stringsTrimSuffix":   strings.TrimSuffix,
	"strconvQuote":        strconv.Quote,
	"strconvUnquote":      strconv.Unquote,
	"strconvParseBool":    strconv.ParseBool,
	"strconvParseInt":     strconv.ParseInt,
	"strconvParseUint":    strconv.ParseUint,
	"strconvParseFloat":   strconv.ParseFloat,
	"strconvFormatBool":   strconv.FormatBool,
	"strconvFormatInt":    strconv.FormatInt,
	"strconvFormatUint":   strconv.FormatUint,
	"strconvFormatFloat":  strconv.FormatFloat,
	"strconvAoti":         strconv.Atoi,
	"strconvItoa":         strconv.Itoa,

	"intAdd": func(v1 int, v2 int) int {
		return v1 + v2
	},

	"intNeg": func(v1 int) int {
		return -v1
	},

	"int64Add": func(v1 int64, v2 int64) int64 {
		return v1 + v2
	},

	"int64Neg": func(v1 int64) int64 {
		return -v1
	},

	"float32Add": func(v1 float32, v2 float32) float32 {
		return v1 + v2
	},

	"float32Neg": func(v1 float32) float32 {
		return -v1
	},

	"float64Add": func(v1 float64, v2 float64) float64 {
		return v1 + v2
	},

	"float64Neg": func(v1 float64) float64 {
		return -v1
	},

	"k8sStatefulSetID": func() (id int, err error) {
		var hostname string
		if hostname = os.Getenv("HOSTNAME"); hostname == "" {
			if hostname, err = os.Hostname(); err != nil {
				return
			}
		}
		splits := strings.Split(hostname, "-")
		if len(splits) < 2 {
			err = errors.New("invalid stateful-set hostname")
			return
		}
		id, err = strconv.Atoi(splits[len(splits)-1])
		return
	},
}
