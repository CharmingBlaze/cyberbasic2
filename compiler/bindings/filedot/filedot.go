// Package filedot exposes file.* aliases over std file foreigns.
package filedot

import (
	"cyberbasic/compiler/bindings/modfacade"
	"cyberbasic/compiler/vm"
)

var fileV2 = map[string]string{
	"read":       "ReadFile",
	"readfile":   "ReadFile",
	"write":      "WriteFile",
	"writefile":  "WriteFile",
	"loadtext":   "LoadText",
	"savetext":   "SaveText",
	"delete":     "DeleteFile",
	"deletefile": "DeleteFile",
	"copy":       "CopyFile",
	"copyfile":   "CopyFile",
	"listdir":    "ListDir",
	"dir":        "Dir",
}

// RegisterFileDot registers global "file".
func RegisterFileDot(v *vm.VM) {
	v.SetGlobal("file", modfacade.New(v, fileV2))
}
