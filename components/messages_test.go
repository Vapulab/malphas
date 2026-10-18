package components

import (
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/vapulab/malphas"
)

// Every field of every locale is filled: strings non-empty, funcs non-nil and
// returning text, plural maps with an "other" form. Walks the struct with
// reflect, so a field added later is covered without touching this test.
func TestEveryLocaleHasEveryMessage(t *testing.T) {
	for _, l := range malphas.Locales() {
		m := reflect.ValueOf(msgFor(l))
		for i := 0; i < m.NumField(); i++ {
			name, f := m.Type().Field(i).Name, m.Field(i)
			switch f.Kind() {
			case reflect.String:
				if strings.TrimSpace(f.String()) == "" {
					t.Errorf("%s.%s is empty", l.Segment(), name)
				}
			case reflect.Func:
				if f.IsNil() {
					t.Errorf("%s.%s is nil", l.Segment(), name)
					continue
				}
				for _, n := range []int{0, 1, 2, 5} {
					if out := callWith(f, n); strings.TrimSpace(out) == "" {
						t.Errorf("%s.%s(%d) is empty", l.Segment(), name, n)
					}
				}
			case reflect.Map:
				other := f.MapIndex(reflect.ValueOf("other"))
				if !other.IsValid() || !strings.Contains(other.String(), "%d") {
					t.Errorf("%s.%s needs an \"other\" form with %%d", l.Segment(), name)
				}
			case reflect.Array:
				for j := 0; j < f.Len(); j++ {
					if strings.TrimSpace(f.Index(j).String()) == "" {
						t.Errorf("%s.%s[%d] is empty", l.Segment(), name, j)
					}
				}
			default:
				t.Errorf("%s.%s: unexpected kind %s; teach this test about it", l.Segment(), name, f.Kind())
			}
		}
	}
}

// callWith calls a message func, feeding n to every int parameter and a
// plausible value to the others.
func callWith(f reflect.Value, n int) string {
	args := make([]reflect.Value, f.Type().NumIn())
	for i := range args {
		switch in := f.Type().In(i); in {
		case reflect.TypeOf(0):
			args[i] = reflect.ValueOf(n)
		case reflect.TypeOf(""):
			args[i] = reflect.ValueOf("report.pdf")
		case reflect.TypeOf(time.Month(0)):
			args[i] = reflect.ValueOf(time.October)
		case reflect.TypeOf(time.Time{}):
			args[i] = reflect.ValueOf(time.Date(2026, 10, 20, 0, 0, 0, 0, time.UTC))
		default:
			args[i] = reflect.Zero(in)
		}
	}
	return f.Call(args)[0].String()
}

func TestRussianResultsPlural(t *testing.T) {
	ru := msgFor(malphas.RU)
	for n, want := range map[int]string{
		1: "1 результат", 2: "2 результата", 4: "4 результата", 5: "5 результатов",
		11: "11 результатов", 12: "12 результатов", 14: "14 результатов",
		21: "21 результат", 22: "22 результата", 25: "25 результатов", 111: "111 результатов",
	} {
		if got := ru.Results(n); got != want {
			t.Errorf("Results(%d) = %q, want %q", n, got, want)
		}
	}
}

func TestEnglishResultsPlural(t *testing.T) {
	en := msgFor(malphas.EN)
	for n, want := range map[int]string{0: "0 results", 1: "1 result", 2: "2 results"} {
		if got := en.Results(n); got != want {
			t.Errorf("Results(%d) = %q, want %q", n, got, want)
		}
	}
}

func TestChineseResultsHasNoPlural(t *testing.T) {
	zh := msgFor(malphas.ZH)
	for n, want := range map[int]string{1: "1 个结果", 5: "5 个结果"} {
		if got := zh.Results(n); got != want {
			t.Errorf("Results(%d) = %q, want %q", n, got, want)
		}
	}
}

// Portuguese keeps the wording the components had before languages existed.
func TestPortugueseKeepsTodaysWording(t *testing.T) {
	pt := msgFor(malphas.PT)
	if pt.Close != "Fechar" || pt.Results(1) != "1 resultado" || pt.Results(3) != "3 resultados" || pt.Digit(2, 6) != "Dígito 2 de 6" {
		t.Errorf("pt drifted: %q %q %q %q", pt.Close, pt.Results(1), pt.Results(3), pt.Digit(2, 6))
	}
}

func TestMsgFollowsContext(t *testing.T) {
	ctx := malphas.WithLocale(t.Context(), malphas.ES)
	if got := msg(ctx).Close; got != "Cerrar" {
		t.Errorf("msg(es).Close = %q", got)
	}
	if got := msg(t.Context()).Close; got != "Close" {
		t.Errorf("msg(no locale).Close = %q, want English", got)
	}
}
