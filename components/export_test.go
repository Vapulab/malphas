package components

import "github.com/vapulab/malphas"

// msgFor exposes one locale's messages to the tests in this package.
func msgFor(l malphas.Locale) messages { return allMessages[l] }
