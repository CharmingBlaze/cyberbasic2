package std

// stdV2 is a small facade over common file/env helpers; Print and math stay flat.
var stdV2 = map[string]string{
	"readfile":         "ReadFile",
	"writefile":        "WriteFile",
	"loadtext":         "LoadText",
	"savetext":         "SaveText",
	"deletefile":       "DeleteFile",
	"copyfile":         "CopyFile",
	"listdir":          "ListDir",
	"dir":              "Dir",
	"directorylist":    "DirectoryList",
	"getenv":           "GetEnv",
	"iswindowprocess":  "IsWindowProcess",
	"getwindowtitle":   "GetWindowTitle",
	"getwindowwidth":   "GetWindowWidth",
	"getwindowheight":  "GetWindowHeight",
	"spawnwindow":      "SpawnWindow",
	"help":             "HELP",
}
