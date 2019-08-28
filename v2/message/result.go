package message

import (
	"fmt"
	"github.com/gojuukaze/YTask/v2/util"
	"github.com/gojuukaze/YTask/v2/yerrors"
	"github.com/tidwall/gjson"
)

type resultStatusChoice struct {
	Sent         int
	FirstRunning int
	WaitingRetry int
	Running      int
	Success      int
	Failure      int
}

var ResultStatus = resultStatusChoice{
	Sent:         0,
	FirstRunning: 1,
	WaitingRetry: 2,
	Running:      3,
	Success:      4,
	Failure:      5,
}

type Result struct {
	Id         string `json:"id"`
	Status     int    `json:"status"` // 0:sent , 1:first running , 2: waiting to retry , 3: running , 4: success , 5: Failure
	JsonResult string `json:"json_result"`
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
	}
}

func (r Result) get(index int) (gjson.Result, error) {
	gR := gjson.Get(r.JsonResult, fmt.Sprintf("%d", index))
	var err error
	if !gR.Exists() {
		err = yerrors.ErrOutOfRange{}
	} else {
		gR = gR.Get("value")
	}
	return gR, err
}

func (r Result) GetInterface(index int) (interface{}, error) {

	gR := gjson.Get(r.JsonResult, fmt.Sprintf("%d", index))
	if !gR.Exists() {
		return nil, yerrors.ErrOutOfRange{}
	} else {
		v, err := util.GetValueFromJson(gR.String())
		if err != nil {
			return nil, err
		}
		return v.Interface(), err
	}
}
func (r Result) GetInt64(index int) (int64, error) {
	gR, err := r.get(index)
	return gR.Int(), err
}

func (r Result) GetUint64(index int) (uint64, error) {
	gR, err := r.get(index)
	return gR.Uint(), err
}
func (r Result) GetFloat64(index int) (float64, error) {
	gR, err := r.get(index)
	return gR.Float(), err
}

func (r Result) GetBool(index int) (bool, error) {
	gR, err := r.get(index)
	return gR.Bool(), err
}

func (r Result) GetString(index int) (string, error) {
	gR, err := r.get(index)
	return gR.String(), err
}

func (r Result) IsSuccess() bool {
	return r.Status == ResultStatus.Success
}

func (r Result) IsFinish() bool {
	if r.Status == ResultStatus.Success || r.Status == ResultStatus.Failure {
		return true
	}
	return false
}
