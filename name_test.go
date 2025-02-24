package fubar

import (
	"os"
	"testing"
)

func TestName(t *testing.T) {
	if _, err := os.Stat("./go.mod"); err != nil {
		t.Skip(err.Error())
	}
	var gmn string
	var need string
	t.Run("goModName", func(t *testing.T) {
		need = "github.com/yunginnanet/fubar"
		if gmn = goModName(""); gmn != need {
			t.Errorf("\ngot: '%s'\nexpected: '%s'", gmn, need)
		}
		t.Log(gmn)
	})
	t.Run("tail(goModName)", func(t *testing.T) {
		need = "yunginnanet/fubar"
		if gmn = tail(gmn); gmn != need {
			t.Errorf("\ngot: '%s'\nexpected: '%s'", gmn, need)
		}
		t.Log(gmn)
	})
	t.Run("getName", func(t *testing.T) {
		need := gmn
		gmn = getName()
		if gmn != need {
			t.Errorf("\ngot: '%s'\nexpected: '%s'", gmn, need)
		}
	})
}
