package app

import (
	"github.com/spf13/pflag"
	"strings"
)

// nolint: deadcode,unused,varcheck
func initFlag() {
	pflag.CommandLine.SetNormalizeFunc(WordSepNormalizeFunc)
}

// WordSepNormalizeFunc changes all flags that contain "_" separators.
func WordSepNormalizeFunc(_ *pflag.FlagSet, name string) pflag.NormalizedName {
	if strings.Contains(name, "_") {
		return pflag.NormalizedName(strings.Replace(name, "_", "-", -1))
	}

	return pflag.NormalizedName(name)
}
