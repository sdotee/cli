package main

import (
	"runtime/debug"
	"testing"
)

func TestVersionFromBuildInfo(t *testing.T) {
	tests := []struct {
		name        string
		info        debug.BuildInfo
		wantVersion string
		wantTime    string
	}{
		{
			name:        "go install of a tagged release",
			info:        debug.BuildInfo{Main: debug.Module{Version: "v1.4.2"}},
			wantVersion: "1.4.2",
			wantTime:    "unknown",
		},
		{
			name: "build outside a tagged release",
			info: debug.BuildInfo{
				Main:     debug.Module{Version: "(devel)"},
				Settings: []debug.BuildSetting{{Key: "vcs.time", Value: "2026-11-09T00:00:00Z"}},
			},
			wantVersion: "dev",
			wantTime:    "2026-11-09T00:00:00Z",
		},
		{
			name:        "no build info",
			info:        debug.BuildInfo{},
			wantVersion: "dev",
			wantTime:    "unknown",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v, bt := versionFromBuildInfo(&tt.info, "dev", "unknown")
			if v != tt.wantVersion || bt != tt.wantTime {
				t.Errorf("got (%q, %q), want (%q, %q)", v, bt, tt.wantVersion, tt.wantTime)
			}
		})
	}
}
