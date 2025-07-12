package ui

import (
	"fmt"
	"net/url"

	"github.com/zserge/lorca"
)

// UI represents the browser UI

type UI struct {
	l lorca.UI
}

// New creates a new UI instance
func New(url string) (*UI, error) {
	ui, err := lorca.New(url, "", 800, 600)
	if err != nil {
		return nil, err
	}
	return &UI{l: ui}, nil
}

// Run starts the UI and waits for it to close
func (u *UI) Run() {
	<-u.l.Done()
}

// Load loads a new URL in the UI
func (u *UI) Load(url string) {
	u.l.Load(url)
}

// Bind exposes a Go function to the UI's Javascript context
func (u *UI) Bind(name string, f interface{}) error {
	return u.l.Bind(name, f)
}