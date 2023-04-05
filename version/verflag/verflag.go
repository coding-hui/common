/*
 * MIT License
 *
 * Copyright (c) 2023 WeCoding.
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package verflag

import (
	"fmt"
	"os"
	"strconv"

	flag "github.com/spf13/pflag"

	"github.com/wecoding.top/common/version"
)

type versionValue int

// Define some const.
const (
	VersionFalse versionValue = 0
	VersionTrue  versionValue = 1
	VersionRaw   versionValue = 2
)

const strRawVersion string = "raw"

func (v *versionValue) IsBoolFlag() bool {
	return true
}

func (v *versionValue) Get() interface{} {
	return v
}

func (v *versionValue) Set(s string) error {
	if s == strRawVersion {
		*v = VersionRaw
		return nil
	}
	boolVal, err := strconv.ParseBool(s)
	if boolVal {
		*v = VersionTrue
	} else {
		*v = VersionFalse
	}
	return err
}

func (v *versionValue) String() string {
	if *v == VersionRaw {
		return strRawVersion
	}
	return fmt.Sprintf("%v", bool(*v == VersionTrue))
}

// The type of the flag as required by the pflag.Value interface.
func (v *versionValue) Type() string {
	return "version"
}

// VersionVar defines a flag with the specified name and usage string.
func VersionVar(p *versionValue, name string, value versionValue, usage string) {
	*p = value
	flag.Var(p, name, usage)
	// "--version" will be treated as "--version=true"
	flag.Lookup(name).NoOptDefVal = "true"
}

// Version wraps the VersionVar function.
func Version(name string, value versionValue, usage string) *versionValue {
	p := new(versionValue)
	VersionVar(p, name, value, usage)
	return p
}

const versionFlagName = "version"

var versionFlag = Version(versionFlagName, VersionFalse, "Print version information and quit.")

// AddFlags registers this package's flags on arbitrary FlagSets, such that they point to the
// same value as the global flags.
func AddFlags(fs *flag.FlagSet) {
	fs.AddFlag(flag.Lookup(versionFlagName))
}

// PrintAndExitIfRequested will check if the -version flag was passed
// and, if so, print the version and exit.
func PrintAndExitIfRequested() {
	if *versionFlag == VersionRaw {
		fmt.Printf("%#v\n", version.Get())
		os.Exit(0)
	} else if *versionFlag == VersionTrue {
		fmt.Printf("%s\n", version.Get())
		os.Exit(0)
	}
}
