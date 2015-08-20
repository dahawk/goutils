package ver

import "testing"

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
