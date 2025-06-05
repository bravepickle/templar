package core

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestPrintAll(t *testing.T) {
	must := require.New(t)

	buf := bytes.NewBuffer([]byte{})
	fm := NewPrinterFormatter(true, buf)

	fm.Printf("Hello %s\n", "printf")
	fm.Println("Hello", "println")
	fm.Print("Hello ", "print!")

	t.Logf("output: %s", buf.String())

	must.Contains(buf.String(), "Hello printf")
	must.Contains(buf.String(), "Hello println")
	must.Contains(buf.String(), "Hello print!")
	must.Contains("Hello sprintf", fm.Sprintf("Hello %s", "sprintf"))

	fx := NewPrinterFormatter(false, nil)
	must.Panics(func() {
		fx.Printf("Hello %s\n", "printf")
	})
	must.Panics(func() {
		fx.Println("Hello %s\n", "printf")
	})
	must.Panics(func() {
		fx.Print("Hello print")
	})
}
