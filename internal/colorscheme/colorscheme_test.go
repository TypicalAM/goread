package colorscheme

import (
	"reflect"
	"testing"
)

// getTestColorscheme returns a test colorscheme
func getTestColorscheme() Colorscheme {
	return Colorscheme{
		BgDark:   "#141018",
		BgDarker: "#141018",
		Text:     "#e0d3cf",
		TextDark: "#e0d3cf",
		Color1:   "#6DA96E",
		Color2:   "#F89763",
		Color3:   "#C9C755",
		Color4:   "#3A409E",
		Color5:   "#6056A7",
		Color6:   "#9D6CB7",
		Color7:   "#e0d3cf",
	}
}

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

	// Create the test colorscheme to compare with
	testColors := getTestColorscheme()

	if reflect.DeepEqual(colors, testColors) {
		t.Errorf("Colorscheme not loaded correctly")
	}
}

// TestColorschemeConvert if we get an error then the conversion from a pywal file doesn't work correctly
func TestColorschemeConvert(t *testing.T) {
	// Create the colorscheme with a non-existent file path
	colors := newDefault()
	err := colors.Convert("../test/data/non-existent.json")
	if err == nil {
		t.Errorf("Colorscheme conversion didn't fail")
	}

	// Create the colorscheme with a real pywal file
	colors = New("")
	err = colors.Convert("../test/data/pywal.json")
	if err != nil {
		t.Errorf("Colorscheme conversion failed")
	}

	// Create the test colorscheme to compare with
	testColors := getTestColorscheme()

	// Check if the colors are the same
	if !reflect.DeepEqual(colors, testColors) {
		t.Errorf("Colorscheme not converted correctly")
	}
}
