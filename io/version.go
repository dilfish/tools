package io

import (
	"runtime/debug"
	"time"
)

type VersionInfo struct {
	GoVersion  string
	VcsVersion string
	VcsTime    time.Time
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
			vi.VcsTime, _ = time.Parse("2006-01-02T15:04:05Z", setting.Value)
			vi.VcsTime = vi.VcsTime.Local()
		}
	}
	return vi
}
