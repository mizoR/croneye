package croneye

import "time"

type Job struct {
	RunTime time.Time
	Script string
}

func NewJob(runsAt time.Time, script string) *Job {
	job := &Job{RunTime: runsAt, Script: script}
	return job
}

type JobList []Job

func (l JobList) Len() int {
	return len(l)
}

func (l JobList) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func (l JobList) Less(i, j int) bool {
	return l[i].RunTime.Before(l[j].RunTime)
}
