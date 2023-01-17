package colorscheme

import (
	"reflect"
	"testing"
)

// TestColorschemeLoadNoFile if we get an error then the default colorscheme is not generated
func TestColorschemeLoadNoFile(t *testing.T) {
	// Create the colorscheme with a non-existent file path
	colors := New("non-existent")

	// Create the default colorscheme to compare with
	defaultColors := newDefault()

	if !reflect.DeepEqual(colors, defaultColors) {
		t.Errorf("Colorscheme not generated correctly")
	}
}

// TestColorschemeLoadFile if we get an error then the loading doesn't work correctly
func TestColorschemeLoadFile(t *testing.T) {
	// Create the colorscheme with a non-existent file path
	colors := New("../test/data/colorscheme.json")

	// Create the new default to compare with
	testColors := newDefault()

	if reflect.DeepEqual(colors, testColors) {
		t.Errorf("Colorscheme not loaded correctly")
	}
}
