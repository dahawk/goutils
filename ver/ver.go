package ver

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

type Version struct {
	Major int
	Minor int
	Patch int
	Build int
}

// Parse parses the given input string into a Version. A version must have the
// format A.B.C.D and at least a major and minor number to be accepted (unlike
// the semver 2.0.0 specification which requires also a patch number to be present).
func Parse(input string) (Version, error) {
	parts := strings.Split(input, ".")
	if len(parts) < 2 {
		return Version{}, errors.New("no major.minor elements found")
	}

	v := Version{}

	switch {
	case len(parts) > 3:
		build, err := strconv.Atoi(parts[3])
		if err != nil {
			return Version{}, fmt.Errorf("invalid build number: %q", parts[3])
		}

		v.Build = build
		fallthrough

	case len(parts) > 2:
		patch, err := strconv.Atoi(parts[2])
		if err != nil {
			return Version{}, fmt.Errorf("invalid patch number: %q", parts[2])
		}

		v.Patch = patch
		fallthrough

	case len(parts) > 1:
		minor, err := strconv.Atoi(parts[1])
		if err != nil {
			return Version{}, fmt.Errorf("invalid minor number: %q", parts[1])
		}

		v.Minor = minor
		fallthrough

	case len(parts) > 0:
		txt := parts[0]
		if strings.HasPrefix(txt, "v") {
			txt = txt[1:]
		}
		major, err := strconv.Atoi(txt)

		if err != nil {
			return Version{}, fmt.Errorf("invalid major number: %q", parts[0])
		}

		v.Major = major
	}

	return v, nil
}

func (v Version) ord() int {
	return (v.Major * 1000000000) +
		(v.Minor * 1000000) +
		(v.Patch * 1000) +
		v.Build
}

// Compare compares two versions. It returns -1 if v is less than o, 0 if v is
// equal to o and 1 if v is greater than o.
func (v Version) Compare(o Version) int {
	a := v.ord()
	b := o.ord()

	if a < b {
		return -1
	}
	if a > b {
		return 1
	}
	return 0
}

func (v Version) Less(o Version) bool {
	return v.Compare(o) == -1
}

func (v Version) Equals(o Version) bool {
	return v.Compare(o) == 0
}

func (v Version) Greater(o Version) bool {
	return v.Compare(o) == 1
}

func (v Version) String() string {
	if v.Build > 0 {
		return fmt.Sprintf("%d.%d.%d.%d", v.Major, v.Minor, v.Patch, v.Build)
	}

	return fmt.Sprintf("%d.%d.%d", v.Major, v.Minor, v.Patch)
}

type Versions []Version

func (l Versions) Len() int {
	return len(l)
}

func (l Versions) Less(i, j int) bool {
	return l[i].ord() < l[j].ord()
}

func (l Versions) Swap(i, j int) {
	l[j], l[i] = l[i], l[j]
}

func (v Version) Value() (driver.Value, error) {
	return int64(v.ord()), nil
}

func (v *Version) Scan(src interface{}) error {
	if integer, ok := src.(int64); ok {
		if integer < 0 {
			return errors.New("value must be positive")
		}

		converted := int(integer)
		v.Major = converted / 1000000000
		v.Minor = (converted % 1000000000) / 1000000
		v.Patch = ((converted % 1000000000) % 1000000) / 1000
		v.Build = ((converted % 1000000000) % 1000000) % 1000
		return nil
	}

	return errors.New("invalid type")
}
