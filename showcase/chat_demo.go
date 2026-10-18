package showcase

import "strconv"

// The Message Scroller demo's conversation: fixed lines, picked by index. The
// endpoint is stateless, so "new" and "older" messages are just more lines.

var chatLines = []string{
	"A curva do ciclo 32 saiu acima do limiar.",
	"Repito a corrida amanhã cedo?",
	"Repete, e separa uma alíquota de controle.",
	"Os primers chegaram, estão na geladeira 2.",
	"Alguém viu o termociclador livre à tarde?",
	"Reservei das 14h às 16h.",
	"Anotei os desvios no protocolo.",
	"O controle negativo ficou limpo.",
}

var chatAuthors = []struct{ Name, Initials string }{
	{"Ana Lima", "AL"}, {"Bruno Reis", "BR"}, {"Carla Souza", "CS"},
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
		m.Author, m.Initials, m.Mine = "Você", "VC", true
	}
	return m
}

func pad2(n int) string {
	if n < 10 {
		return "0" + strconv.Itoa(n)
	}
	return strconv.Itoa(n)
}
