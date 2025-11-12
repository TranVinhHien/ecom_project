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

func (j *JobScheduler) start() {
	fmt.Println("Starting job scheduler...")
	j.scheduler.Start()
}

func (j *JobScheduler) NewJob(hours int, job func()) error {
	newJob, err := j.scheduler.NewJob(
		gocron.DurationJob(
			time.Duration(hours)*time.Hour,
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

	// start the scheduler
	j.start()
	return nil
}
