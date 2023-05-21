package colorscheme

import (
	"testing"
)

// TestColorschemeLoadNoFile if we get an error then the default colorscheme is not generated
func TestColorschemeLoadNoFile(t *testing.T) {
	// Create the colorscheme with a non-existent file path
	colors := New("non-existent")

	if colors != newDefault() {
		t.Error("Colorscheme not generated correctly")
	}
}

// TestColorschemeLoadFile if we get an error then the loading doesn't work correctly
func TestColorschemeLoadFile(t *testing.T) {
	// Create a colorscheme with a real file path
	colors := New("../test/data/colorscheme.json")

	if colors == newDefault() {
		t.Error("Colorscheme not loaded correctly")
	}
}

// TestPywalConvert if we get an error then the conversion doesn't work correctly
func TestPywalConvert(t *testing.T) {
	// Create the colorscheme with a non-existent file path
	colors := newDefault()

	// Try to convert the colorscheme
	if err := colors.Convert("../test/data/pywal.json"); err != nil {
		t.Error("Colorscheme couldn't convert", err)
	}

	if colors == newDefault() {
		t.Errorf("Colorscheme not converted correctly")
	}
}
