package components_test

import (
	"strings"
	"testing"

	"github.com/vapulab/malphas/components"
)

func TestOTPRendersOneBoxPerDigit(t *testing.T) {
	html := renderToString(t, components.OTP(components.OTPProps{Name: "codigo", Label: "Código"}))
	if got := strings.Count(html, "data-otp-index"); got != 6 {
		t.Errorf("expected 6 boxes by default, got %d: %s", got, html)
	}

	four := renderToString(t, components.OTP(components.OTPProps{Name: "c", Length: 4}))
	if got := strings.Count(four, "data-otp-index"); got != 4 {
		t.Errorf("expected 4 boxes, got %d", got)
	}
}

// The whole point of the component on mobile. Without this pair the OS never
// offers the SMS code, and nothing looks broken — which is why it gets left out.
func TestOTPOffersTheSMSCode(t *testing.T) {
	html := renderToString(t, components.OTP(components.OTPProps{Name: "codigo"}))
	if !strings.Contains(html, `autocomplete="one-time-code"`) {
		t.Errorf("expected one-time-code on the first box, got: %s", html)
	}
	if strings.Count(html, `autocomplete="one-time-code"`) != 1 {
		t.Errorf("only the first box may claim the code, got: %s", html)
	}
	if !strings.Contains(html, `inputmode="numeric"`) {
		t.Errorf("expected a numeric keypad, got: %s", html)
	}
}

// Six unlabelled boxes are six anonymous fields to a screen reader.
func TestOTPBoxesAreLabelled(t *testing.T) {
	html := renderToString(t, components.OTP(components.OTPProps{Name: "codigo", Length: 3}))
	for _, want := range []string{"Digit 1 of 3", "Digit 2 of 3", "Digit 3 of 3"} {
		if !strings.Contains(html, want) {
			t.Errorf("expected %q, got: %s", want, html)
		}
	}
}

// The boxes are the interface; the hidden field is what the form submits.
func TestOTPSubmitsThroughAHiddenField(t *testing.T) {
	html := renderToString(t, components.OTP(components.OTPProps{Name: "codigo"}))
	if !strings.Contains(html, `type="hidden"`) || !strings.Contains(html, `name="codigo"`) {
		t.Errorf("expected a hidden field carrying the value, got: %s", html)
	}
	// The visible boxes must not submit anything of their own.
	if strings.Contains(html, `name="codigo"`) && strings.Count(html, `name=`) != 1 {
		t.Errorf("only the hidden field may be named, got: %s", html)
	}
}

func TestOTPGroupIsLabelled(t *testing.T) {
	html := renderToString(t, components.OTP(components.OTPProps{Name: "codigo", Label: "Código de verificação"}))
	if !strings.Contains(html, `role="group"`) || !strings.Contains(html, "aria-labelledby") {
		t.Errorf("expected the boxes grouped under the label, got: %s", html)
	}
}
