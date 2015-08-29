package croneye

import (
	"bufio"
	"github.com/gorhill/cronexpr"
	"io"
	"regexp"
	"sync"
	"time"
)

type Parser struct {
	FromTime time.Time
	ToTime   time.Time
}

func NewParser(fromTime time.Time, toTime time.Time) *Parser {
	parser := &Parser{FromTime: fromTime, ToTime: toTime}
	return parser
}

func (p Parser) ParseLine(line string) JobList {
	var jobList JobList = []Job{}
	re := regexp.MustCompile("^([^ ]+ +[^ ]+ +[^ ]+ +[^ ]+ +[^ ]+) +(.+)$")

	if !re.MatchString(line) {
		return jobList
	}

	matched := re.FindStringSubmatch(line)
	time := matched[1]
	script := matched[2]

	expr, err := cronexpr.Parse(time)
	if err != nil {
		return jobList
	}

	for nextTime := expr.Next(p.FromTime); !nextTime.After(p.ToTime); nextTime = expr.Next(nextTime) {
		job := NewJob(nextTime, script)
		jobList = append(jobList, *job)
	}

	return jobList
}

func (p Parser) Parse(r io.Reader) JobList {
	var jobList JobList = []Job{}
	var wait sync.WaitGroup

	channel := make(chan JobList)
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		wait.Add(1)

		go func(line string) {
			defer wait.Done()
			channel <- p.ParseLine(line)
		}(scanner.Text())

		go func(list *JobList) {
			*list = append(*list, <-channel...)
		}(&jobList)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}

	wait.Wait()

	return jobList
}
