package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"

	"github.com/fatih/color"
	"github.com/go-ping/ping"
	"github.com/jessevdk/go-flags"
)

// nolint:gochecknoglobals
var pingu = []string{
	` ...        .     ...   ..    ..     .........           `,
	` ...     ....          ..  ..      ... .....  .. ..      `,
	` ...    .......      ...         ... . ..... BBBBBBB     `,
	`.....  ........ .BBBBBBBBBBBBBBB.....  ... BBBBBBBBBB.  .`,
	` .... ........BBBBBBBBBBBBBBBBBBBBB.  ... BBBBBBBBBBB    `,
	`      ....... BBWWWWBBBBBBBBBBBBBBBB.... BBBBBBBBBBBB    `,
	`.    .  .... BBWWBBWWBBBBBBBBBBWWWWBB... BBBBBBBBBBB     `,
	`   ..   ....BBBBWWWWBBRRRRRRBBWWBBWWB.. .BBBBBBBBBBB     `,
	`    .       BBBBBBBBRRRRRRRRRRBWWWWBB.   .BBBBBBBBBB     `,
	`   ....     .BBBBBBBBRRRRRRRRBBBBBBBB.      BBBBBBBB     `,
	`  .....      .  BBBBBBBBBBBBBBBBBBBB.        BBBBBBB.    `,
	`......     .. . BBBBBBBBBBBBBBBBBB . .      .BBBBBBB     `,
	`......       BBBBBBBBBBBBBBBBBBBBB  .      .BBBBBBB      `,
	`......   .BBBBBBBBBBBBBBBBBBYYWWBBBBB  ..  BBBBBBB       `,
	`...    . BBBBBBBBBBBBBBBBYWWWWWWWWWBBBBBBBBBBBBBB.       `,
	`       BBBBBBBBBBBBBBBBYWWWWWWWWWWWWWBBBBBBBBB .         `,
	`      BBBBBBBBBBBBBBBYWWWWWWWWWWWWWWWWBB    .            `,
	`     BBBBBBBBBBBBBBBYWWWWWWWWWWWWWWWWWWW  ........       `,
	`  .BBBBBBBBBBBBBBBBYWWWWWWWWWWWWWWWWWWWW    .........    `,
	` .BBBBBBBBBBBBBBBBYWWWWWWWWWWWWWWWWWWWWWW       .... . . `,
}

// nolint:gochecknoglobals
var (
	appName        = "pingu"
	appUsage       = "[OPTIONS] HOST"
	appDescription = "`ping` command but with pingu"
	appVersion     = "???"
	appRevision    = "???"
)

type exitCode int

const (
	exitCodeOK exitCode = iota
	exitCodeErrArgs
	exitCodeErrPing
)

type options struct {
	Version bool `short:"V" long:"version" description:"Show version"`
	Count int  `short:"c" long:"count" default:"20" description:"Stop after <count> replies"`
}

func main() {
	code, err := run(os.Args[1:])
	if err != nil {
		fmt.Fprintf(
			color.Error,
			"[ %v ] %s\n",
			color.New(color.FgRed, color.Bold).Sprint("ERROR"),
			err,
		)
	}

	os.Exit(int(code))
}

func run(cliArgs []string) (exitCode, error) {
	var opts options
	parser := flags.NewParser(&opts, flags.Default)
	parser.Name = appName
	parser.Usage = appUsage
	parser.ShortDescription = appDescription
	parser.LongDescription = appDescription

	args, err := parser.ParseArgs(cliArgs)
	if err != nil {
		if flags.WroteHelp(err) {
			return exitCodeOK, nil
		}

		return exitCodeErrArgs, fmt.Errorf("parse error: %w", err)
	}

	if opts.Version {
		// nolint:forbidigo
		fmt.Printf("%s: v%s-rev%s\n", appName, appVersion, appRevision)

		return exitCodeOK, nil
	}

	if len(args) == 0 {
		// nolint:goerr113
		return exitCodeErrArgs, errors.New("must requires an argument")
	}

	if 1 < len(args) {
		// nolint:goerr113
		return exitCodeErrArgs, errors.New("too many arguments")
	}

	pinger, err := initPinger(args[0], opts)
	if err != nil {
		return exitCodeOK, fmt.Errorf("an error occurred while initializing pinger: %w", err)
	}

	if err := pinger.Run(); err != nil {
		return exitCodeErrPing, fmt.Errorf("an error occurred when running ping: %w", err)
	}

	return exitCodeOK, nil
}

func initPinger(host string, opts options) (*ping.Pinger, error) {
	pinger, err := ping.NewPinger(host)
	if err != nil {
		return nil, fmt.Errorf("failed to init pinger %w", err)
	}
	pinger.Count = opts.Count

	// Listen for Ctrl-C.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	go func() {
		<-c
		pinger.Stop()
	}()

	color.New(color.FgHiWhite, color.Bold).Printf(
		"PING %s (%s) type `Ctrl-C` to abort\n",
		pinger.Addr(),
		pinger.IPAddr(),
	)

	pinger.OnRecv = pingerOnrecv
	pinger.OnFinish = pingerOnFinish

	if runtime.GOOS == "windows" {
		pinger.SetPrivileged(true)
	}

	return pinger, nil
}

// nolint:forbidigo
func pingerOnrecv(pkt *ping.Packet) {
	fmt.Printf("%s seq=%s %sbytes from %s: ttl=%s time=%s\n",
		renderASCIIArt(pkt.Seq),
		color.New(color.FgHiYellow, color.Bold).Sprintf("%d", pkt.Seq),
		color.New(color.FgHiBlue, color.Bold).Sprintf("%d", pkt.Nbytes),
		color.New(color.FgWhite, color.Bold).Sprintf("%s", pkt.IPAddr),
		color.New(color.FgHiCyan, color.Bold).Sprintf("%d", pkt.Ttl),
		color.New(color.FgHiMagenta, color.Bold).Sprintf("%v", pkt.Rtt),
	)
}

// nolint:forbidigo
func pingerOnFinish(stats *ping.Statistics) {
	color.New(color.FgWhite, color.Bold).Printf(
		"\n───────── %s ping statistics ─────────\n",
		stats.Addr,
	)
	fmt.Printf(
		"%s: %v transmitted => %v received (%v loss)\n",
		color.New(color.FgHiWhite, color.Bold).Sprintf("PACKET STATISTICS"),
		color.New(color.FgHiBlue, color.Bold).Sprintf("%d", stats.PacketsSent),
		color.New(color.FgHiGreen, color.Bold).Sprintf("%d", stats.PacketsRecv),
		color.New(color.FgHiRed, color.Bold).Sprintf("%v%%", stats.PacketLoss),
	)
	fmt.Printf(
		"%s: min=%v avg=%v max=%v stddev=%v\n",
		color.New(color.FgHiWhite, color.Bold).Sprintf("ROUND TRIP"),
		color.New(color.FgHiBlue, color.Bold).Sprintf("%v", stats.MinRtt),
		color.New(color.FgHiCyan, color.Bold).Sprintf("%v", stats.AvgRtt),
		color.New(color.FgHiGreen, color.Bold).Sprintf("%v", stats.MaxRtt),
		color.New(color.FgMagenta, color.Bold).Sprintf("%v", stats.StdDevRtt),
	)
}

func renderASCIIArt(idx int) string {
	if len(pingu) <= idx {
		return strings.Repeat(" ", len(pingu[0]))
	}

	line := pingu[idx]

	line = colorize(line, 'R', color.New(color.FgHiRed, color.Bold))
	line = colorize(line, 'Y', color.New(color.FgHiYellow, color.Bold))
	line = colorize(line, 'B', color.New(color.FgHiBlack, color.Bold))
	line = colorize(line, 'W', color.New(color.FgHiWhite, color.Bold))

	return line
}

func colorize(text string, target rune, color *color.Color) string {
	return strings.ReplaceAll(
		text,
		string(target),
		color.Sprint("#"),
	)
}
