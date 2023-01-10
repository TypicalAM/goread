package colorscheme

import (
	"reflect"
	"testing"
)

// getTestColorscheme returns a test colorscheme
func getTestColorscheme() Colorscheme {
	return Colorscheme{
		BgDark:   "#040612",
		BgDarker: "#040612",
		Text:     "#98ccdc",
		TextDark: "#98ccdc",
		Color1:   "#625160",
		Color2:   "#BD4354",
		Color3:   "#985063",
		Color4:   "#BA9C6A",
		Color5:   "#1E5AA6",
		Color6:   "#C25C9F",
		Color7:   "#98ccdc",
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
