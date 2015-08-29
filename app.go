package croneye

import (
	"fmt"
	"os"
	"sort"
	"time"
)

type App struct {
	Parser Parser
}

func NewApp(fromTime time.Time, toTime time.Time) *App {
	parser := NewParser(fromTime, toTime)
	app := &App{Parser: *parser}
	return app
}

func (a App) Run(file *os.File) {
	var jobList JobList = []Job{}
	jobList = a.Parser.Parse(file)

	sort.Sort(jobList)
	for i := 0; i < len(jobList); i++ {
		fmt.Printf("%v\t%s\n", jobList[i].RunTime, jobList[i].Script)
	}
}
