package client

import (
	"bytes"
	"sync"

	"github.com/Koshsky/Intellivue-api/pkg/intellivue/packages"
)

// PollDataResult определяет интерфейс для обработки результатов опроса данных
type PollDataResult interface {
	UnmarshalBinary(r *bytes.Reader) error
	ShowInfo(mu *sync.Mutex, indent int)
	GetPollInfoList() interface{}
}

// pollDataResultWrapper обертка для SinglePollDataResult
type pollDataResultWrapper struct {
	result *packages.SinglePollDataResult
}

func (w *pollDataResultWrapper) UnmarshalBinary(r *bytes.Reader) error {
	return w.result.UnmarshalBinary(r)
}

func (w *pollDataResultWrapper) ShowInfo(mu *sync.Mutex, indent int) {
	w.result.ShowInfo(mu, indent)
}

func (w *pollDataResultWrapper) GetPollInfoList() interface{} {
	return w.result.PollMdibDataReply.PollInfoList
}

type pollDataResultLinkedWrapper struct {
	result *packages.SinglePollDataResultLinked
}

func (w *pollDataResultLinkedWrapper) UnmarshalBinary(r *bytes.Reader) error {
	return w.result.UnmarshalBinary(r)
}

func (w *pollDataResultLinkedWrapper) ShowInfo(mu *sync.Mutex, indent int) {
	w.result.ShowInfo(mu, indent)
}

func (w *pollDataResultLinkedWrapper) GetPollInfoList() interface{} {
	return w.result.PollMdibDataReply.PollInfoList
}
