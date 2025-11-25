package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
)

type Logger struct {
	level  Level
	output io.Writer
	colors bool
}

func New(level Level, colors bool) *Logger {
	return &Logger{
		level:  level,
		output: os.Stdout,
		colors: colors,
	}
}

func (l *Logger) timestamp() string {
	return time.Now().Format("15:04:05")
}

func (l *Logger) Debug(format string, args ...interface{}) {
	if l.level <= LevelDebug {
		l.log(color.New(color.FgCyan), "DEBUG", format, args...)
	}
}

func (l *Logger) Info(format string, args ...interface{}) {
	if l.level <= LevelInfo {
		l.log(color.New(color.FgBlue), "INFO ", format, args...)
	}
}

func (l *Logger) Watch(format string, args ...interface{}) {
	if l.level <= LevelInfo {
		l.log(color.New(color.FgMagenta, color.Bold), "WATCH", format, args...)
	}
}

func (l *Logger) Runner(format string, args ...interface{}) {
	if l.level <= LevelInfo {
		l.log(color.New(color.FgYellow, color.Bold), "EXEC ", format, args...)
	}
}

func (l *Logger) Warn(format string, args ...interface{}) {
	if l.level <= LevelWarn {
		l.log(color.New(color.FgYellow), "WARN ", format, args...)
	}
}

func (l *Logger) Error(format string, args ...interface{}) {
	if l.level <= LevelError {
		l.log(color.New(color.FgRed, color.Bold), "ERROR", format, args...)
	}
}

func (l *Logger) Success(format string, args ...interface{}) {
	if l.level <= LevelInfo {
		l.log(color.New(color.FgGreen, color.Bold), "✓ OK ", format, args...)
	}
}

func (l *Logger) Banner(title, version string) {
	if l.level > LevelInfo {
		return
	}

	if l.colors {
		cyan := color.New(color.FgCyan, color.Bold)
		magenta := color.New(color.FgMagenta)

		fmt.Fprintf(l.output, "\n")
		fmt.Fprintf(l.output, "%s\n", cyan.Sprint("╔════════════════════════════════════════════════════════════╗"))
		fmt.Fprintf(l.output, "%s  %s %s\n",
			cyan.Sprint("║"),
			magenta.Sprint("🕵️  GoWatch - File Watcher & Auto-Runner"),
			cyan.Sprint("║"))
		fmt.Fprintf(l.output, "%s  %s%s%s\n",
			cyan.Sprint("║"),
			strings.Repeat(" ", 28),
			color.New(color.Faint).Sprintf("v%s", version),
			strings.Repeat(" ", 28-len(version)))
		fmt.Fprintf(l.output, "%s\n", cyan.Sprint("╚════════════════════════════════════════════════════════════╝"))
		fmt.Fprintf(l.output, "\n")
	} else {
		fmt.Fprintf(l.output, "\n%s v%s\n\n", title, version)
	}
}

func (l *Logger) Section(title string) {
	if l.level > LevelInfo {
		return
	}

	if l.colors {
		blue := color.New(color.FgBlue, color.Bold)
		fmt.Fprintf(l.output, "\n%s\n", blue.Sprintf("── %s ──", title))
	} else {
		fmt.Fprintf(l.output, "\n-- %s --\n", title)
	}
}

func (l *Logger) CommandOutput(line string, isError bool) {
	if l.level > LevelInfo {
		return
	}

	prefix := "  │ "
	if l.colors {
		if isError {
			fmt.Fprintf(l.output, "%s%s\n",
				color.New(color.Faint).Sprint(prefix),
				color.New(color.FgRed).Sprint(line))
		} else {
			fmt.Fprintf(l.output, "%s%s\n",
				color.New(color.FgCyan, color.Faint).Sprint(prefix),
				line)
		}
	} else {
		fmt.Fprintf(l.output, "%s%s\n", prefix, line)
	}
}

func (l *Logger) Separator() {
	if l.level > LevelInfo {
		return
	}

	if l.colors {
		fmt.Fprintf(l.output, "%s\n", color.New(color.Faint).Sprint(strings.Repeat("─", 60)))
	} else {
		fmt.Fprintf(l.output, "%s\n", strings.Repeat("-", 60))
	}
}

func (l *Logger) CommandStart(cmd string) {
	if l.level > LevelInfo {
		return
	}

	if l.colors {
		fmt.Fprintf(l.output, "%s %s %s\n",
			color.New(color.FgYellow, color.Bold).Sprint("▶"),
			color.New(color.FgWhite).Sprint("Running:"),
			color.New(color.FgCyan).Sprint(cmd))
	} else {
		fmt.Fprintf(l.output, "▶ Running: %s\n", cmd)
	}
}

func (l *Logger) CommandEnd(cmd string, exitCode int, duration time.Duration) {
	if l.level > LevelInfo {
		return
	}

	durationStr := l.formatDuration(duration)

	if l.colors {
		if exitCode == 0 {
			fmt.Fprintf(l.output, "%s %s %s %s\n",
				color.New(color.FgGreen, color.Bold).Sprint("✓"),
				color.New(color.FgGreen).Sprint("Completed:"),
				color.New(color.Faint).Sprint(cmd),
				color.New(color.FgGreen, color.Faint).Sprintf("(%s)", durationStr))
		} else {
			fmt.Fprintf(l.output, "%s %s %s %s %s\n",
				color.New(color.FgRed, color.Bold).Sprint("✗"),
				color.New(color.FgRed).Sprint("Failed:"),
				color.New(color.Faint).Sprint(cmd),
				color.New(color.FgRed).Sprintf("(exit: %d)", exitCode),
				color.New(color.Faint).Sprintf("(%s)", durationStr))
		}
	} else {
		if exitCode == 0 {
			fmt.Fprintf(l.output, "✓ Completed: %s (%s)\n", cmd, durationStr)
		} else {
			fmt.Fprintf(l.output, "✗ Failed: %s (exit: %d) (%s)\n", cmd, exitCode, durationStr)
		}
	}
}

func (l *Logger) formatDuration(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return fmt.Sprintf("%.2fs", d.Seconds())
}

func (l *Logger) log(c *color.Color, prefix, format string, args ...interface{}) {
	timestamp := l.timestamp()
	msg := fmt.Sprintf(format, args...)

	if l.colors {
		fmt.Fprintf(l.output, "%s %s %s\n",
			color.New(color.Faint).Sprint(timestamp),
			c.Sprintf("[%s]", prefix),
			msg)
	} else {
		fmt.Fprintf(l.output, "%s [%s] %s\n", timestamp, prefix, msg)
	}
}
