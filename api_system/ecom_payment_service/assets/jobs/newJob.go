package assets_jobs

import (
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
)

type JobScheduler struct {
	scheduler gocron.Scheduler
}

func NewJobScheduler() (*JobScheduler, error) {
	scheduler, err := gocron.NewScheduler()
	if err != nil {
		// handle error
		return nil, fmt.Errorf("failed to create job scheduler: %w", err)
	}
	return &JobScheduler{
		scheduler: scheduler,
	}, nil
}

func (j *JobScheduler) Start() {
	fmt.Println("Starting job scheduler...")
	j.scheduler.Start()
}

func (j *JobScheduler) NewJob(minute, hours, date int, job func()) error {
	interval := time.Duration(hours)*time.Hour + time.Duration(minute)*time.Minute + time.Duration(date)*time.Hour*24

	newJob, err := j.scheduler.NewJob(
		gocron.DurationJob(
			interval,
		// time.Duration(1)*time.Minute,
		// time.Duration(30)*time.Second,
		),
		gocron.NewTask(
			func() {
				fmt.Println("run task")
				go job()
			},
		),
	)

	if err != nil {
		return err
		// handle error
	}
	// each job has a unique id
	fmt.Println(newJob.ID())
	fmt.Printf("Job %s scheduled to run every %v\n", newJob.ID(), interval)
	// start the scheduler
	// j.start()
	return nil
}
