package marathon

import "fmt"

type Task struct {
	AppId              AppId               `json:"appId"`
	SlaveId            string              `json:"slaveId"`
	Id                 string              `json:"id"`
	Ports              []int               `json:"ports"`
	ServicePorts       []int               `json:"servicePorts"`
	IpAddresses        []IpAddress         `json:"ipAddresses"`
	Host               string              `json:"host"`
	HealthCheckResults []HealthCheckResult `json:"healthCheckResults"`
	stagedAt           string              `json:"stagedAt"`
	startedAt          string              `json:"startedAt"`
	version            string              `json:"version"`
}
type TaskList []Task

type TaskStatus string
type AppId string

const (
	TaskStaging  TaskStatus = "TASK_STAGING"
	TaskStarting            = "TASK_STARTING"
	TaskRunning             = "TASK_RUNNING"
	TaskFinished            = "TASK_FINISHED"
	TaskFailed              = "TASK_FAILED"
	TaskKilled              = "TASK_KILLED"
	TaskLost                = "TASK_LOST"
	// this is kinda wack, because TASK_ANY isnt a real task status in marathon
	// but its used for client.GetAllTasks to signify any status
	TaskAny = "TASK_ANY"
)

func (a TaskList) Len() int           { return len(a) }
func (a TaskList) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a TaskList) Less(i, j int) bool { return a[i].Id < a[j].Id }
func (t *Task) String() string {
	return fmt.Sprintf("%s on %s:%v", t.Id, t.Host, t.Ports)
}
