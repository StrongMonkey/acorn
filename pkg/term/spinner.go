package term

import (
	"fmt"
	"os"
	"strings"

	"github.com/acorn-io/runtime/pkg/install/progress"
	"github.com/pterm/pterm"
)

const (
	trimSpace = 13
)

func init() {
	// Help capture text cleaner
	pterm.SetDefaultOutput(os.Stderr)
	pterm.ThemeDefault.SuccessMessageStyle = *pterm.NewStyle(pterm.FgLightGreen)
	// Customize default error.
	pterm.Success.Prefix = pterm.Prefix{
		Text:  " ✔",
		Style: pterm.NewStyle(pterm.FgLightGreen),
	}
	pterm.Error.Prefix = pterm.Prefix{
		Text:  " ✗  ERROR:",
		Style: pterm.NewStyle(pterm.FgLightRed),
	}
	pterm.Warning.Prefix = pterm.Prefix{
		Text:  " •  WARNING:",
		Style: pterm.NewStyle(pterm.FgLightYellow),
	}
	pterm.Info.Prefix = pterm.Prefix{
		Text: " •",
	}
}

type Builder struct {
}

func (b *Builder) New(msg string) progress.Progress {
	return NewSpinner(msg)
}

type Spinner struct {
	spinner *pterm.SpinnerPrinter
	text    string
	lastMsg string
}

func NewSpinner(text string) *Spinner {
	var err error
	spinner := pterm.DefaultSpinner.
		WithRemoveWhenDone(false).
		// Src: https://github.com/gernest/wow/blob/master/spin/spinners.go#L335
		WithSequence(`  ⠋ `, `  ⠙ `, `  ⠹ `, `  ⠸ `, `  ⠼ `, `  ⠴ `, `  ⠦ `, `  ⠧ `, `  ⠇ `, `  ⠏ `)
	if !pterm.RawOutput {
		spinner, err = spinner.Start(trimWidth(text))
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println()
	}

	return &Spinner{
		spinner: spinner,
		text:    text,
	}
}

func trimWidth(msg string) string {
	if width := pterm.GetTerminalWidth(); width > trimSpace && len(msg)+trimSpace > width {
		msg = msg[:width-trimSpace]
	}
	return msg
}

func (s *Spinner) Infof(format string, v ...any) {
	msg := trimWidth(strings.TrimSpace(fmt.Sprintf(s.text+": "+format, v...)))
	if s.lastMsg == msg {
		return
	}
	s.lastMsg = msg
	s.spinner.UpdateText(msg)
}

func (s *Spinner) Fail(err error) error {
	if err == nil {
		s.Success()
		return nil
	}
	s.spinner.Fail(err)
	return err
}

func (s *Spinner) SuccessWithWarning(format string, v ...any) {
	msg := fmt.Sprintf(format, v...)
	s.spinner.Warning(msg)
}

func (s *Spinner) Success() {
	s.spinner.Success(s.text)
}
