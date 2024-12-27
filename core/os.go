package core

import "runtime"

func GetOS() OS {
	return OS{
		Arch:   runtime.GOARCH,
		System: OsNameMap[runtime.GOOS],
	}
}

var OsNameMap = map[string]string{
	"darwin":  "macos",
	"linux":   "linux",
	"windows": "windows",
}

type OS struct {
	Version string
	Arch    string
	System  string
}

func (o *OS) String() string {
	return o.System + "-" + o.Arch
}
