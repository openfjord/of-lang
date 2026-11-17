// SPDX-License-Identifier: Apache-2.0

package main

import (
	"runtime/debug"
	"strings"
)

// version is set at release build time with -ldflags "-X main.version=<version>".
var version = ""

// toolVersion is the version a release build was stamped with; failing that, the
// module version Go recorded (for go install ...@v0.1.0); failing that, "dev".
func toolVersion() string {
	if version != "" {
		return version
	}
	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return strings.TrimPrefix(info.Main.Version, "v")
	}
	return "dev"
}
