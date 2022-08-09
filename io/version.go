package io

import (
	"runtime/debug"
)

type VersionInfo struct {
	GoVersion  string
	VcsVersion string
	BuildTime  string
}

func Version() VersionInfo {
	var vi VersionInfo
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return vi
	}
	vi.GoVersion = bi.GoVersion
	for _, setting := range bi.Settings {
		if setting.Key == "vcs.revision" {
			vi.VcsVersion = setting.Value
		}
		if setting.Key == "vcs.time" {
			vi.BuildTime = setting.Value
		}
	}
	return vi
}
