package drive

import (
	"context"
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

type brokerItem struct {
	Msg [][]byte `json:"msg"`
}
type brokerStruct map[string]brokerItem
type backendItem struct {
	Data   []byte    `json:"data"`
	ExTime time.Time `json:"ex_time"`
}
type backendStruct map[string]backendItem

type LocalDrive struct {
	brokerLock  UnsafeFileLock
	backendLock UnsafeFileLock
	brokerPath  string
	backendPath string
	isBroker    bool
}

func NewLocalDrive(isBroker bool) LocalDrive {
	var d = LocalDrive{}
	if isBroker {
		d.isBroker = true
		d.brokerLock = NewFileLock("ytask_local_broker.lock")
		d.brokerPath = filepath.Join(os.TempDir(), "ytask_local_broker.json")
	} else {
		d.backendLock = NewFileLock("ytask_local_backend.lock")
		d.backendPath = filepath.Join(os.TempDir(), "ytask_local_backend.json")
	}

	d.Init()
	return d
}

func (d LocalDrive) Init() {
	if d.isBroker {
		d.brokerLock.Init()
		os.WriteFile(d.brokerPath, []byte("{}"), os.FileMode(0600))

	} else {
		d.backendLock.Init()
		os.WriteFile(d.backendPath, []byte("{}"), os.FileMode(0600))
	}
}
func (d LocalDrive) Close() {
	d.brokerLock.Unlock()
	d.backendLock.Unlock()
	os.Remove(d.brokerPath)
	os.Remove(d.backendPath)
}

func (d LocalDrive) getBrokerData() brokerStruct {
	var data brokerStruct
	b, _ := os.ReadFile(d.brokerPath)
	json.Unmarshal(b, &data)
	return data
}
func (d LocalDrive) setBrokerData(data brokerStruct) {
	b, _ := json.Marshal(data)
	os.WriteFile(d.brokerPath, b, os.FileMode(0600))
}

func (d LocalDrive) getBackendData() backendStruct {
	var data backendStruct
	b, _ := os.ReadFile(d.backendPath)
	json.Unmarshal(b, &data)
	return data
}

func (d LocalDrive) setBackendData(data backendStruct) {
	b, _ := json.Marshal(data)
	os.WriteFile(d.backendPath, b, os.FileMode(0600))
}

func (d LocalDrive) Set(key string, value []byte, exTime int) error {
	err := d.backendLock.Lock()
	if err != nil {
		return err
	}
	defer d.backendLock.Unlock()
	var t = time.Time{}
	if exTime > 0 {
		t = time.Now().Add(time.Duration(exTime) * time.Second)
	}
	data := d.getBackendData()
	data[key] = backendItem{value, t}
	d.setBackendData(data)
	return nil
}

func (d LocalDrive) Get(key string) ([]byte, error) {
	err := d.backendLock.Lock()
	if err != nil {
		return nil, err
	}
	defer d.backendLock.Unlock()
	data := d.getBackendData()
	r, ok := data[key]
	if !ok {
		return nil, NilResultError
	}
	if !r.ExTime.IsZero() && r.ExTime.Before(time.Now()) {
		return nil, NilResultError
	}
	return r.Data, nil
}

func (d LocalDrive) push(queueName string, value []byte, isRight bool) error {
	err := d.brokerLock.Lock()
	if err != nil {
		return err
	}
	defer d.brokerLock.Unlock()
	data := d.getBrokerData()
	item, _ := data[queueName]
	if isRight {
		item.Msg = rPush(item.Msg, value)
	} else {
		item.Msg = lPush(item.Msg, value)
	}
	data[queueName] = item
	d.setBrokerData(data)
	return nil
}

func (d LocalDrive) RPush(queueName string, value []byte) error {
	return d.push(queueName, value, true)

}

func (d LocalDrive) LPush(queueName string, value []byte) error {
	return d.push(queueName, value, false)
}
func (d LocalDrive) lPop(queueName string) ([]byte, error) {
	err := d.brokerLock.Lock()
	if err != nil {
		return nil, err
	}
	defer d.brokerLock.Unlock()
	data := d.getBrokerData()
	item, ok := data[queueName]
	if !ok {
		return nil, EmptyQueueError
	}
	b, msg := lPop(item.Msg)
	if b == nil {
		return nil, EmptyQueueError
	}
	item.Msg = msg
	data[queueName] = item
	d.setBrokerData(data)
	return b, nil

}

func (d LocalDrive) LPop(queueName string) ([]byte, error) {
	// 由于会有多个协程执行这个操作，这里超时时间短一点，尽快让出锁
	ctx, _ := context.WithTimeout(context.Background(), time.Millisecond*300)
	for {
		b, err := d.lPop(queueName)
		if err == nil {
			return b, nil
		}
		select {
		case <-ctx.Done():
			return nil, EmptyQueueError
		case <-time.After(time.Millisecond * 100):
		}
	}

}
