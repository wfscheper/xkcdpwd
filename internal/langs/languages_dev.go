// +build dev

package langs

import "net/http"

var Languages http.FileSystem = http.Dir("languages")
