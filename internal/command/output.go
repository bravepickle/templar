package command

import (
	"fmt"
	"io"
	"strings"
)

// Contains common helper utilities for reuse in code

// ANSI colors for Stdout terminal
const (
	colorRed    = "\033[31m"
	colorGreen  = "\033[32m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"

	textStyleBold  = "\033[1m"
	colorDarkGray  = "\033[90m"
	textStyleReset = "\033[0m" // reset and text style
)

// Text styles for Stdout formatting
var textStyles = map[string]string{
	`<info>`:    colorBlue,
	`<debug>`:   colorGreen,
	`<comment>`: colorYellow,
	`<alert>`:   colorRed,
	`<muted>`:   colorDarkGray,
	`<bold>`:    textStyleBold,
	`<reset>`:   textStyleReset,
}

// PrinterFormatter is a wrapper that formats messages for writer
type PrinterFormatter struct {
	// NoColor do not colorize output
	NoColor bool

	// Writer is a writer for input messages
	Writer io.Writer

	// Styles contains map for placeholders and replacements
	Styles map[string]string
}

// Printf prints out to Stdout formatted and colorized string
func (f *PrinterFormatter) Printf(msg string, args ...any) {
	for k, v := range textStyles {
		msg = strings.ReplaceAll(msg, k, v)
	}

	if _, err := fmt.Fprintf(f.Writer, msg, args...); err != nil {
		panic(fmt.Errorf(`failed to printf: %w`, err))
	}
}

func (f *PrinterFormatter) Sprintf(msg string, args ...any) string {
	for k, v := range textStyles {
		msg = strings.ReplaceAll(msg, k, v)
	}

	return fmt.Sprintf(msg, args...)
}

// Print prints out to Stdout and colorized string
func (f *PrinterFormatter) Print(args ...any) {
	if len(args) > 0 { // styles only first argument
		if msg, ok := args[0].(string); ok {
			for k, v := range textStyles {
				msg = strings.ReplaceAll(msg, k, v)
			}

			args[0] = msg
		}
	}

	if _, err := fmt.Fprint(f.Writer, args...); err != nil {
		panic(fmt.Errorf(`failed to print: %w`, err))
	}
}

// Println prints formatted text with styles and new line at the end
func (f *PrinterFormatter) Println(msg string, args ...any) {
	//f.Print("aaa", args...)

	for k, v := range textStyles {
		msg = strings.ReplaceAll(msg, k, v)
	}

	if _, err := fmt.Fprintln(f.Writer, append([]any{msg}, args...)...); err != nil {
		panic(fmt.Errorf(`failed to println: %w`, err))
	}
}

func (f *PrinterFormatter) Init() {
	if f.NoColor { // fill in with blanks
		for k := range textStyles {
			f.Styles[k] = ``
		}
	} else {
		f.Styles = textStyles
	}
}

func NewPrinterFormatter(NoColor bool, w io.Writer) *PrinterFormatter {
	f := &PrinterFormatter{NoColor: NoColor, Writer: w}

	f.Init()

	return f
}
