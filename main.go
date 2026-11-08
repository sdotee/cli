//
// Copyright (c) 2025 S.EE Development Team
//
// This source code is licensed under the MIT License,
// which is located in the LICENSE file in the source tree's root directory.
//
// File: main.go
// Author: S.EE Development Team <dev@s.ee>
// File Created: 2025-12-22 22:23:57
//
// Modified By: S.EE Development Team <dev@s.ee>
// Last Modified: 2025-12-23 11:06:22
//

package main

import (
	"runtime/debug"
	"strings"

	"github.com/sdotee/cli/cmd"
)

var (
	BuildVersion = "dev"
	BuildTime    = "unknown"
)

func main() {
	if BuildVersion == "dev" {
		if info, ok := debug.ReadBuildInfo(); ok {
			BuildVersion, BuildTime = versionFromBuildInfo(info, BuildVersion, BuildTime)
		}
	}
	cmd.Execute(BuildVersion, BuildTime)
}

// versionFromBuildInfo fills in version and build time from Go build info
// for builds without ldflags, such as go install. The "v" prefix is dropped
// to match the versions injected by make and GoReleaser.
func versionFromBuildInfo(info *debug.BuildInfo, version, buildTime string) (string, string) {
	if v := info.Main.Version; v != "" && v != "(devel)" {
		version = strings.TrimPrefix(v, "v")
	}
	for _, s := range info.Settings {
		if s.Key == "vcs.time" && s.Value != "" {
			buildTime = s.Value
		}
	}
	return version, buildTime
}
