package showcase

import "strconv"

// The Message Scroller demo's conversation: fixed lines, picked by index. The
// endpoint is stateless, so "new" and "older" messages are just more lines.

var chatLines = []string{
	"The curve for cycle 32 came out above the threshold.",
	"Should I rerun it first thing tomorrow?",
	"Rerun it, and set aside a control aliquot.",
	"The primers arrived, they are in fridge 2.",
	"Is the thermocycler free this afternoon?",
	"I booked it from 2 to 4 pm.",
	"I logged the deviations in the protocol.",
	"The negative control came out clean.",
}

var chatAuthors = []struct{ Name, Initials string }{
	{"Amara Okafor", "AO"}, {"Ben Carter", "BC"}, {"Chloé Martin", "CM"},
}

// ChatMessage is message number n of the demo conversation.
type ChatMessage struct {
	Author, Initials, Text, Time string
	Mine                         bool
}

func DemoChatMessage(n int) ChatMessage {
	if n < 0 {
		n = -n
	}
	a := chatAuthors[n%len(chatAuthors)]
	m := ChatMessage{Author: a.Name, Initials: a.Initials, Text: chatLines[n%len(chatLines)]}
	m.Time = "14:" + pad2((n*3)%60)
	if n%4 == 2 {
		m.Author, m.Initials, m.Mine = "You", "ME", true
	}
	return m
}

func pad2(n int) string {
	if n < 10 {
		return "0" + strconv.Itoa(n)
	}
	return strconv.Itoa(n)
}
