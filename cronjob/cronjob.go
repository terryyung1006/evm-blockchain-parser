package cronjob

type Icronjob interface {
	Run()
}

type CronJob struct {
	Interval int
}

func (cjg CronjobGroup) RunCronJobs() {
	for i := 0; i < len(cjg); i++ {
		go cjg[i].Run(10)
	}
}
