package humanslogger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/olekukonko/tablewriter"
)

const (
	ansiBold           = "\033[1m"
	ansiReset          = "\033[0m"
	ansiFaint          = "\033[2m"
	ansiResetFaint     = "\033[22m"
	ansiBrightRed      = "\033[91m"
	ansiBrightGreen    = "\033[92m"
	ansiBrightYellow   = "\033[93m"
	ansiBrightRedFaint = "\033[91;2m"
)

func Init() {
	slog.SetDefault(NewHandler(
		WithLevel(slog.LevelDebug),
	))
}

type HandlerOption func(*Handler)

func WithWriter(in io.Writer) HandlerOption {
	return func(h *Handler) {
		h.w = in
	}
}

func WithLevel(l slog.Level) HandlerOption {
	return func(h *Handler) {
		h.l = l
	}
}

func NewHandler(opts ...HandlerOption) *Handler {
	h := Handler{
		l: slog.LevelInfo,
		w: os.Stdout,
	}

	for _, f := range opts {
		f(&h)
	}

	return &h
}

type Handler struct {
	w io.Writer
	l slog.Level
}

func (h *Handler) Enabled(context.Context, slog.Level) bool {
	return true
}

func (h *Handler) Handle(ctx context.Context, r slog.Record) error {
	fmt.Fprintf(h.w, "[%v][%v] %v%v%v\n", h.FormatColor(r.Level), r.Time.Format(time.Kitchen), ansiBold, r.Message, ansiReset)

	if r.NumAttrs() != 0 {
		table := tablewriter.NewWriter(h.w)
		r.Attrs(func(a slog.Attr) bool {
			table.Append([]string{a.Key, a.Value.String()})
			return true
		})

		table.Render()
	}

	return nil
}

func (h *Handler) WithAttrs(in []slog.Attr) slog.Handler {
	return h
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return h
}

func (h *Handler) FormatColor(level slog.Level) string {
	switch level {
	case slog.LevelInfo:
		return ansiBrightGreen + level.String() + ansiReset
	case slog.LevelError:
		return ansiBrightRed + level.String() + ansiReset
	case slog.LevelWarn:
		return ansiBrightYellow + level.String() + ansiReset
	}

	return level.String()
}
