package cronjob

type Icronjob interface {
	Run()
}

type CronJob struct {
	Interval int
}

func RunCronJobs(cronJobs []Icronjob) {
	for i := 0; i < len(cronJobs); i++ {
		go func(cronJob Icronjob) {
			for {
				cronJob.Run()
			}
		}(cronJobs[i])
	}
}
