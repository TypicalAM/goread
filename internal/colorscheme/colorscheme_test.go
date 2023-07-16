package colorscheme

import (
	"testing"
)

// TestColorschemeLoadNoFile if we get an error then the default colorscheme is not generated
func TestColorschemeLoadNoFile(t *testing.T) {
	// Create the colorscheme with a non-existent file path
	if _, err := New("non-existent"); err == nil {
		t.Error("Colorscheme couldn't be created", err)
	}
}

// TestColorschemeLoadFile if we get an error then the loading doesn't work correctly
func TestColorschemeLoadFile(t *testing.T) {
	// Create a colorscheme with a real file path
	colors, err := New("../test/data/colorscheme.json")
	if err != nil {
		t.Error("Colorscheme couldn't be created", err)
	}

	if colors == NewDefault() {
		t.Error("Colorscheme not loaded correctly")
	}
}

// TestPywalConvert if we get an error then the conversion doesn't work correctly
func TestPywalConvert(t *testing.T) {
	// Create the colorscheme with a non-existent file path
	colors := NewDefault()

	// Try to convert the colorscheme
	if err := colors.Convert("../test/data/pywal.json"); err != nil {
		t.Error("Colorscheme couldn't convert", err)
	}

	if colors == NewDefault() {
		t.Errorf("Colorscheme not converted correctly")
	}
}
