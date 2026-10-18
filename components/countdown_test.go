package components_test

import (
	"strings"
	"testing"

	"github.com/vapulab/vapulab/malphas/components"
)

func TestCountdownRendersInitialTimeServerSide(t *testing.T) {
	// 1h 2m 3s
	html := renderToString(t, components.Countdown(components.CountdownProps{Seconds: 3723}))
	if !strings.Contains(html, "01:02:03") {
		t.Errorf("expected the starting time rendered before JS runs, got: %s", html)
	}
}

func TestCountdownShortFormatDropsHours(t *testing.T) {
	html := renderToString(t, components.Countdown(components.CountdownProps{Seconds: 125, Format: "mm:ss"}))
	if !strings.Contains(html, "02:05") {
		t.Errorf("expected mm:ss, got: %s", html)
	}
	if strings.Contains(html, "00:02:05") {
		t.Errorf("mm:ss must not include hours, got: %s", html)
	}
}

// role=timer marks it as a live-updating clock. aria-live is deliberately off:
// a value announced every second makes the page unusable with a screen reader.
func TestCountdownIsATimerThatDoesNotAnnounceEverySecond(t *testing.T) {
	html := renderToString(t, components.Countdown(components.CountdownProps{Seconds: 60}))
	if !strings.Contains(html, `role="timer"`) {
		t.Errorf("expected role=timer, got: %s", html)
	}
	if !strings.Contains(html, `aria-live="off"`) {
		t.Errorf("expected aria-live=off, got: %s", html)
	}
}

func TestCountdownZeroRendersZero(t *testing.T) {
	html := renderToString(t, components.Countdown(components.CountdownProps{Seconds: 0}))
	if !strings.Contains(html, "00:00:00") {
		t.Errorf("expected a zeroed clock, got: %s", html)
	}
}

// Negative input would otherwise render something like 0-1:-1:-30, and hand a
// negative starting value to the ticker.
func TestCountdownClampsNegativeInput(t *testing.T) {
	html := renderToString(t, components.Countdown(components.CountdownProps{Seconds: -30}))
	if !strings.Contains(html, "00:00:00") {
		t.Errorf("expected the clock clamped to zero, got: %s", html)
	}
	// The clamp has to reach the behavior too, not just the rendered text.
	if !strings.Contains(html, "countdown(0,") {
		t.Errorf("expected the behavior to start from zero, got: %s", html)
	}
}
