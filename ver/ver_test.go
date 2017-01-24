package ver

import (
	"testing"
	"time"
)

func TestParse(t *testing.T) {
	table := map[string]Version{
		"0.0.0":   Version{0, 0, 0, 0},
		"0.01.0":  Version{0, 1, 0, 0},
		"2.0.0.1": Version{2, 0, 0, 1},
		"3.3.3":   Version{3, 3, 3, 0},
		"0.0.1":   Version{0, 0, 1, 0},
		"0.0.0.1": Version{0, 0, 0, 1},
		"1.3.3.7": Version{1, 3, 3, 7},
		"1.203":   Version{1, 203, 0, 0},
	}

	for input, expected := range table {
		got, err := Parse(input)
		if err != nil {
			t.Errorf("failed to parse %q: %v", input, err)
		}

		if expected.Compare(got) != 0 {
			t.Errorf("failed to parse %q: expected %v, got %v", input, expected, got)
		}
	}
}

func TestString(t *testing.T) {
	table := map[Version]string{
		Version{0, 0, 0, 0}:   "0.0.0",
		Version{0, 1, 0, 0}:   "0.1.0",
		Version{2, 0, 0, 1}:   "2.0.0.1",
		Version{3, 3, 3, 0}:   "3.3.3",
		Version{0, 0, 1, 0}:   "0.0.1",
		Version{0, 0, 0, 1}:   "0.0.0.1",
		Version{1, 3, 3, 7}:   "1.3.3.7",
		Version{1, 203, 0, 0}: "1.203.0",
	}

	for input, expected := range table {
		got := input.String()

		if expected != got {
			t.Errorf("failed to stringify %#v: expected %s, got %s", input, expected, got)
		}
	}
}

func TestErrors(t *testing.T) {
	table := map[string]string{
		"0":       "no major.minor elements found",
		"0.0.0.a": "invalid build number: \"a\"",
		"0.0.a.0": "invalid patch number: \"a\"",
		"0.a.0.0": "invalid minor number: \"a\"",
		"a.0.0.0": "invalid major number: \"a\"",
	}

	for input, expected := range table {
		got, err := Parse(input)
		if err == nil {
			t.Errorf("parse succeeded, but should not: %q -> %v", input, got)
		}

		if expected != err.Error() {
			t.Errorf("unknown error for %q: expected %v, got %v", input, expected, err.Error())
		}
	}
}

func TestCompare(t *testing.T) {
	a := Version{1, 2, 0, 0}
	b := Version{1, 2, 1, 0}
	c := Version{2, 0, 1, 0}

	if x := a.Compare(b); x != -1 {
		t.Errorf("expected %v < %v, got %d", a, b, x)
	}

	if x := b.Compare(a); x != 1 {
		t.Errorf("expected %v > %v, got %d", b, a, x)
	}

	if x := b.Compare(c); x != -1 {
		t.Errorf("expected %v < %v, got %d", b, c, x)
	}

	if x := b.Compare(b); x != 0 {
		t.Errorf("expected %v = %v, got %d", b, b, x)
	}
}

func TestValue(t *testing.T) {
	table := map[Version]int64{
		Version{8, 5, 1, 0}:    8005001000,
		Version{4, 7, 0, 0}:    4007000000,
		Version{0, 23, 0, 0}:   23000000,
		Version{14, 0, 0, 100}: 14000000100,
		Version{0, 0, 0, 0}:    0,
		Version{0, 0, 1, 1}:    1001,
		Version{0, 1, 0, 0}:    1000000,
		Version{1, 0, 0, 0}:    1000000000,
	}

	for input, expected := range table {
		got, err := input.Value()

		if err != nil {
			t.Errorf("failed to retrieve value from %#v: %v", input, err)
		}

		if expected != got {
			t.Errorf("failed to retrieve value %#v: expected %d, got %d", input, expected, got)
		}
	}
}

func TestScan(t *testing.T) {
	table := map[int64]string{
		8005001000:     "8.5.1",
		4007000000:     "4.7.0",
		23000000:       "0.23.0",
		14000000100:    "14.0.0.100",
		0:              "0.0.0",
		1001:           "0.0.1.1",
		1000000:        "0.1.0",
		1000000000:     "1.0.0",
		14000000000100: "14000.0.0.100",
	}

	var version Version
	for input, expected := range table {
		err := version.Scan(input)

		if err != nil {
			t.Errorf("failed to scan %d: %v", input, err)
		}

		got := version.String()
		if expected != got {
			t.Errorf("failed to scan %d: expected %q, got %q", input, expected, got)
		}
	}
}

func TestScanErrors(t *testing.T) {
	table := []interface{}{
		"3.15.6",
		8005001000.1514,
		-1,
		nil,
		time.Now(),
	}

	var version Version
	for _, input := range table {
		err := version.Scan(input)
		if err == nil {
			t.Errorf("scan succeeded but should not: %#v -> %q", input, version)
		}
	}
}
