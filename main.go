package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"strings"

	"github.com/gookit/color"
	"github.com/jessevdk/go-flags"
	probing "github.com/prometheus-community/pro-bing"
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
	Count     int  `short:"c" long:"count" default:"20" description:"Stop after <count> replies"`
	Privilege bool `short:"P" long:"privilege" description:"Enable privileged mode"`
	Version   bool `short:"V" long:"version" description:"Show version"`
}

func main() {
	code, err := run(os.Args[1:])
	if err != nil {
		fmt.Fprintf(
			os.Stderr,
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

func initPinger(host string, opts options) (*probing.Pinger, error) {
	pinger, err := probing.NewPinger(host)
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

	color.New(color.FgWhite, color.Bold).Printf(
		"PING %s (%s) type `Ctrl-C` to abort\n",
		pinger.Addr(),
		pinger.IPAddr(),
	)

	pinger.OnRecv = pingerOnrecv
	pinger.OnFinish = pingerOnFinish

	if opts.Privilege || runtime.GOOS == "windows" {
		pinger.SetPrivileged(true)
	}

	return pinger, nil
}

func pingerOnrecv(pkt *probing.Packet) {
	fmt.Printf(
		"%s seq=%s %sbytes from %s: ttl=%s time=%s\n",
		renderASCIIArt(pkt.Seq),
		color.New(color.FgYellow, color.Bold).Sprintf("%d", pkt.Seq),
		color.New(color.FgBlue, color.Bold).Sprintf("%d", pkt.Nbytes),
		color.New(color.FgWhite, color.Bold).Sprintf("%s", pkt.IPAddr),
		color.New(color.FgCyan, color.Bold).Sprintf("%d", pkt.TTL),
		color.New(color.FgMagenta, color.Bold).Sprintf("%v", pkt.Rtt),
	)
}

func pingerOnFinish(stats *probing.Statistics) {
	color.New(color.FgWhite, color.Bold).Printf(
		"\n───────── %s ping statistics ─────────\n",
		stats.Addr,
	)
	fmt.Printf(
		"%s: %v transmitted => %v received (%v loss)\n",
		color.New(color.FgWhite, color.Bold).Sprintf("PACKET STATISTICS"),
		color.New(color.FgBlue, color.Bold).Sprintf("%d", stats.PacketsSent),
		color.New(color.FgGreen, color.Bold).Sprintf("%d", stats.PacketsRecv),
		color.New(color.FgRed, color.Bold).Sprintf("%v%%", stats.PacketLoss),
	)
	fmt.Printf(
		"%s: min=%v avg=%v max=%v stddev=%v\n",
		color.New(color.FgWhite, color.Bold).Sprintf("ROUND TRIP"),
		color.New(color.FgBlue, color.Bold).Sprintf("%v", stats.MinRtt),
		color.New(color.FgCyan, color.Bold).Sprintf("%v", stats.AvgRtt),
		color.New(color.FgGreen, color.Bold).Sprintf("%v", stats.MaxRtt),
		color.New(color.FgMagenta, color.Bold).Sprintf("%v", stats.StdDevRtt),
	)
}

func renderASCIIArt(idx int) string {
	if len(pingu) <= idx {
		idx %= len(pingu)
	}

	line := pingu[idx]

	line = colorize(line, 'R', color.New(color.FgRed, color.Bold))
	line = colorize(line, 'Y', color.New(color.FgYellow, color.Bold))
	line = colorize(line, 'B', color.New(color.FgBlack, color.Bold))
	line = colorize(line, 'W', color.New(color.FgWhite, color.Bold))

	return line
}

func colorize(text string, target rune, color color.PrinterFace) string {
	return strings.ReplaceAll(
		text,
		string(target),
		color.Sprint("#"),
	)
}
