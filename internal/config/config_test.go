package config

import (
	"slices"
	"testing"

	"github.com/TypicalAM/goread/internal/ui/browser"
)

func getCfg(t *testing.T) *Config {
	myCfg, err := New("../test/data/goread.yml")
	if err != nil {
		t.Errorf("error creating config object: %v", err)
	}

	if err = myCfg.Load(); err != nil {
		t.Errorf("error loading file: %v", err)
	}

	return myCfg
}

// TestConfigLoadNoFile if we get an error then the config file is not loaded correctly
func TestConfigLoadNoFile(t *testing.T) {
	myCfg, err := New("non-existent")
	if err != nil {
		t.Errorf("error creating config object: %v", err)
	}

	if err = myCfg.Load(); err != nil {
		t.Errorf("no error returned when loading non-existent file")
	}
}

// TestConfigLoadFile if we get an error then the config file is not loaded correctly
func TestConfigLoadFile(t *testing.T) {
	cfg := getCfg(t)
	if len(cfg.Keymap) != 5 {
		t.Errorf("incorrect number of keymap categories, expected 5, got %d", len(cfg.Keymap))
	}

	if _, ok := cfg.Keymap["category"]; !ok {
		t.Error("incorrect map, expected category to exist")
	}

	keys := browser.DefaultKeymap.CloseTab.Keys()
	if !slices.Contains(keys, "EXTRA") {
		t.Errorf("incorrect keys loaded, expected 'EXTRA' in %v", keys)
	}
}

// TestConfigLoadFile if we get an error then the config loader doesn't recognize non-existent categories
func TestConfigLoadBadCategory(t *testing.T) {
	myCfg, err := New("../test/data/goread_bad_category.yml")
	if err != nil {
		t.Fatalf("error creating config object: %v", err)
	}

	if err = myCfg.Load(); err == nil {
		t.Fatalf("expected error when loading file with non-existent category, but got none")
	}
}

// TestConfigLoadFile if we get an error then the config file allows for a bad attribute of category to be defined
func TestConfigLoadBadBind(t *testing.T) {
	myCfg, err := New("../test/data/goread_bad_bind.yml")
	if err != nil {
		t.Errorf("error creating config object: %v", err)
	}

	if err = myCfg.Load(); err == nil {
		t.Error("expected error when loading file with invalid binds to the 'browser' category, but got none")
	}
}

// TestConfigLoadFile if we get an error then the config file allows for bindings with no keys attached
func TestConfigLoadNoKeys(t *testing.T) {
	myCfg, err := New("../test/data/goread_no_keys.yml")
	if err != nil {
		t.Errorf("error creating config object: %v", err)
	}

	if err = myCfg.Load(); err == nil {
		t.Error("expected error when loading file with bindings missing keys, but got none")
	}
}
