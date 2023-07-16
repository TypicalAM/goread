package colorscheme

import (
	"testing"
)

// TestColorschemeLoadNoFile if we get an error then the default colorscheme is not generated
func TestColorschemeLoadNoFile(t *testing.T) {
	if _, err := New("non-existent"); err != nil {
		t.Error("Colorscheme couldn't be created", err)
	}
}

// TestColorschemeLoadFile if we get an error then the loading doesn't work correctly
func TestColorschemeLoadFile(t *testing.T) {
	colors, err := New("../test/data/colorscheme.json")
	if err != nil {
		t.Error("Colorscheme couldn't be created", err)
	}

	if err = colors.Load(); err != nil {
		t.Error("Colorscheme couldn't load", err)
	}
}

// TestPywalConvert if we get an error then the conversion doesn't work correctly
func TestPywalConvert(t *testing.T) {
	colors := Default

	if err := colors.Convert("../test/data/pywal.json"); err != nil {
		t.Error("Colorscheme couldn't convert", err)
	}

	if colors == Default {
		t.Errorf("Colorscheme not converted correctly")
	}
}
