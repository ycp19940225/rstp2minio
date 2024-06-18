package monitorstorage

import (
	"fmt"
	"github.com/pkg/errors"
	"sync"
)

type MonitorStorage struct {
	lock    sync.Mutex
	streams map[string]*Stream
}

func NewMonitorStorage() (m *MonitorStorage) {
	m = &MonitorStorage{
		streams: make(map[string]*Stream),
	}
	return
}

func (m *MonitorStorage) AddStream(uri, clientId string, mode Mode) (stream *Stream, err error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if v, ok := m.streams[clientId]; ok {
		stream = v
	} else {
		stream, err = NewStream(
			WithStreamUri(uri),
			WithStreamMode(mode),
		)
		if err != nil {
			return
		}
		stream.CodecData, err = stream.GetCodecData()
		if err != nil {
			return
		}
	}
	m.streams[clientId] = stream
	return
}

func (m *MonitorStorage) RemoveStream(name string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	stream, ok := m.streams[name]
	if !ok {
		return
	}
	stream.StopSignal <- true
	delete(m.streams, name)
	return
}

func (m *MonitorStorage) UpdateStream(name string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	stream, ok := m.streams[name]
	if !ok {
		return
	}
	stream.UpdateSignal <- name
	return
}

func (m *MonitorStorage) ReconnectStream(uri, clientId string) (*Stream, error) {
	m.lock.Lock()
	defer m.lock.Unlock()
	if stream, ok := m.streams[clientId]; ok {
		rtspClient, err := dial(uri)
		stream.RTSPClient = rtspClient
		if err != nil {
			return nil, err
		}
		stream.CodecData, err = stream.GetCodecData()
		if err != nil {
			return nil, err
		}
		return stream, nil
	}
	return nil, errors.New(fmt.Sprintf("no uri:%s clientId:%s stream", uri, clientId))
}

func (m *MonitorStorage) DeleteStream(name string) {
	m.lock.Lock()
	defer m.lock.Unlock()
	_, ok := m.streams[name]
	if !ok {
		return
	}
	delete(m.streams, name)
	return
}

func (m *MonitorStorage) CheckStream(name string) bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	_, ok := m.streams[name]
	if ok {
		return true
	}
	return false
}
