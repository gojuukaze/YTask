package message

import (
	"fmt"
	"github.com/gojuukaze/YTask/v3/core/util/yjson"
)

type resultStatusChoice struct {
	Sent         int
	FirstRunning int
	WaitingRetry int
	Running      int
	Success      int
	Failure      int
	Expired      int
	Abort        int // 手动中止任务
}

var ResultStatus = resultStatusChoice{
	Sent:         0,
	FirstRunning: 1,
	WaitingRetry: 2,
	Running:      3,
	Success:      4,
	Failure:      5,
	Expired:      6,
	Abort:        7, // 手动中止任务
}

type workflowStatusChoice struct {
	Waiting string
	Running string
	Success string
	Failure string
	Expired string
	Abort   string
}

var WorkflowStatus = workflowStatusChoice{
	Waiting: "waiting",
	Running: "running",
	Success: "success",
	Failure: "failure",
	Expired: "expired",
	Abort:   "abort", // 手动中止任务
}

var StatusToWorkflowStatus = map[int]string{
	ResultStatus.Sent:         WorkflowStatus.Waiting,
	ResultStatus.FirstRunning: WorkflowStatus.Running,
	ResultStatus.WaitingRetry: WorkflowStatus.Running,
	ResultStatus.Running:      WorkflowStatus.Running,
	ResultStatus.Success:      WorkflowStatus.Success,
	ResultStatus.Failure:      WorkflowStatus.Failure,
	ResultStatus.Expired:      WorkflowStatus.Expired,
	ResultStatus.Abort:        WorkflowStatus.Abort,
}

type Result struct {
	Id         string      `json:"id"`
	Status     int         `json:"status"` // 0:sent , 1:first running , 2: waiting to retry , 3: running , 4: success , 5: Failure
	FuncReturn []string    `json:"func_return"`
	RetryCount int         `json:"retry_count"`
	Workflow   [][2]string `json:"workflow"` // [["workName","status"],] ;  status: waiting , running , success , failure , expired , abort
	Err        string      `json:"err"`
}

func NewResult(id string) Result {
	return Result{
		Id: id,
	}
}
func (r Result) GetBackendKey() string {
	return "YTask:Backend:" + r.Id
}

func (r *Result) SetStatusRunning() {
	if r.Status == ResultStatus.Sent {
		r.Status = ResultStatus.FirstRunning
	} else {
		r.Status = ResultStatus.Running
		r.RetryCount += 1
	}
}

func (r Result) Get(index int, v interface{}) error {

	err := yjson.YJson.UnmarshalFromString(r.FuncReturn[index], v)
	return err
}

func (r Result) Gets(args ...interface{}) error {
	for i, v := range args {
		err := yjson.YJson.UnmarshalFromString(r.FuncReturn[i], v)
		if err != nil {
			return err
		}
	}
	return nil
}

// 过时: 此方法只能用于v1.0.0，高版本中，如果值为int64,uint64类型，会导致获取的值不对
// Deprecated: only can use in v1.0.0
func (r Result) GetInterface(index int) (interface{}, error) {

	var result interface{}

	err := yjson.YJson.Unmarshal([]byte(r.FuncReturn[index]), &result)
	return result, err
}
func (r Result) GetInt64(index int) (int64, error) {
	var v int64
	err := r.Get(index, &v)
	return v, err
}

func (r Result) GetUint64(index int) (uint64, error) {
	var v uint64
	err := r.Get(index, &v)
	return v, err
}
func (r Result) GetFloat64(index int) (float64, error) {
	var v float64
	err := r.Get(index, &v)
	return v, err
}

func (r Result) GetBool(index int) (bool, error) {
	var v bool
	err := r.Get(index, &v)
	return v, err
}

func (r Result) GetString(index int) (string, error) {
	var v string
	err := r.Get(index, &v)
	return v, err
}

func (r Result) IsSuccess() bool {
	return r.Status == ResultStatus.Success
}

func (r Result) IsFailure() bool {
	if r.Status == ResultStatus.Failure || r.Status == ResultStatus.Expired || r.Status == ResultStatus.Abort {
		return true
	}
	return false
}

func (r Result) IsFinish() bool {
	if r.Status == ResultStatus.Success || r.IsFailure() {
		return true
	}
	return false
}

// 结束任务标志
func GetAbortKey(id string) string {
	return fmt.Sprintf("Abort:%s", id)

}

func NewAbortResult(id string) Result {
	return Result{
		Id: GetAbortKey(id),
	}
}
