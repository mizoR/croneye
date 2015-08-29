package main

import (
	"flag"
	"fmt"
	"github.com/jinzhu/now"
	"github.com/mizoR/croneye"
	"golang.org/x/crypto/ssh/terminal"
	"os"
	"runtime"
	"time"
)

var fromTime, toTime time.Time
var timeFormat = "2006-01-02 15:04:05"

func init() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage
  $ %s [OPTIONS]
Options
`, os.Args[0])
		flag.PrintDefaults()
		os.Exit(0)
	}
	from := flag.String("from", time.Now().Format(timeFormat), "From time (default: Current time)")
	to := flag.String("to", time.Now().Add(24*time.Hour).Format(timeFormat), "To time (default: After 1 day since current time)")
	flag.Parse()

	fromTime = now.MustParse(*from)
	toTime = now.MustParse(*to)
}

func main() {
	if terminal.IsTerminal(int(os.Stdin.Fd())) {
		fmt.Fprintf(os.Stderr, `Error:
  You must supply cron definition via stdin
Example:
  $ crontab -l | %s [OPTIONS]

`, os.Args[0])
		os.Exit(1)
	}

	app := croneye.NewApp(fromTime, toTime)
	app.Run(os.Stdin)
}
