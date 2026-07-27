package malphas_test

import (
	"io/fs"
	"testing"

	"github.com/darkanum/vapulab/malphas"
)

func TestAssetsContainsFavicon(t *testing.T) {
	data, err := fs.ReadFile(malphas.Assets(), "favicon.svg")
	if err != nil {
		t.Fatalf("read favicon.svg: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("favicon.svg is empty")
	}
}

func TestAssetsContainsUIJS(t *testing.T) {
	data, err := fs.ReadFile(malphas.Assets(), "malphas.js")
	if err != nil {
		t.Fatalf("read ui.js: %v", err)
	}
	if len(data) == 0 {
		t.Fatal("ui.js is empty")
	}
}
