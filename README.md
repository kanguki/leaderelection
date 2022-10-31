# What

I create the repo to solve the problem of 1 cron job running on multiple instances

# Usage

```
type Job1 struct {
	Le
	timeoutDecideLeader int
	canRunLater bool
	idempotent bool
	...other attribute
} 

func NewJob1() (Job1, err) {
	le, err := NewXXXLe(....) //e.g. NewNatsLe
	if err != nil {
		//log error and save failed job for instance
		return err
	}
	return Job1{Le: le} 
}

func (j Job1) Run() {
	defer j.Le.Close()
	if j.Le.AmITheLeader(j.timeoutDecideLeader) {
		//check job's state pending
		//do job
		//save job's state done
	}
}
```
