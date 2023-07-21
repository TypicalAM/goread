package theme

import (
	"testing"
)

// TestThemeLoadNoFile if we get an error then the default theme is not generated
func TestThemeLoadNoFile(t *testing.T) {
	if _, err := New("non-existent"); err != nil {
		t.Error("Theme couldn't be created", err)
	}
}

// TestThemeLoadFile if we get an error then the loading doesn't work correctly
func TestThemeLoadFile(t *testing.T) {
	colors, err := New("../test/data/colorscheme.json")
	if err != nil {
		t.Error("Theme couldn't be created", err)
	}

	if err = colors.Load(); err != nil {
		t.Error("Theme couldn't load", err)
	}
}

// TestPywalConvert if we get an error then the conversion doesn't work correctly
func TestPywalConvert(t *testing.T) {
	colors := Default

	if err := colors.Convert("../test/data/pywal.json"); err != nil {
		t.Error("Theme couldn't convert", err)
	}

	if colors == Default {
		t.Errorf("Theme not converted correctly")
	}
}
