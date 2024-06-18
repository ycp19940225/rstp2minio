package monitorstorage

import (
	"context"
	"sync"
	"time"

	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/format/rtspv2"
	"github.com/pkg/errors"
)

type Stream struct {
	ops          StreamOptions
	lock         sync.Mutex
	RTSPClient   *rtspv2.RTSPClient
	CodecData    []av.CodecData
	StopSignal   chan bool
	UpdateSignal chan string
}

func NewStream(options ...func(*StreamOptions)) (st *Stream, err error) {
	ops := getStreamOptionsOrSetDefault(nil)
	for _, f := range options {
		f(ops)
	}
	st = &Stream{
		StopSignal:   make(chan bool, 1),
		UpdateSignal: make(chan string, 1),
	}
	if ops.uri == "" {
		err = errors.Errorf("invalid uri")
		return
	}
	client, err := dial(ops.uri)
	if err != nil {
		err = errors.Errorf("dail %s failed: %v", ops.uri, err)
		return
	}
	st.ops = *ops
	st.RTSPClient = client
	return
}

func (st *Stream) GetCodecData() (data []av.CodecData, err error) {
	st.lock.Lock()
	defer st.lock.Unlock()
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	ticker := time.NewTimer(100 * time.Millisecond)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				err = errors.Errorf("GetCodecData timeout")
				return
			}
			err = errors.Errorf("GetCodecData failed: %v", ctx.Err())
			return
		case <-ticker.C:
			data = st.RTSPClient.CodecData
			if !st.RTSPClient.WaitCodec {
				if data != nil {
					return
				}
			}
		}
	}
}

func dial(uri string) (client *rtspv2.RTSPClient, err error) {
	client, err = rtspv2.Dial(
		rtspv2.RTSPClientOptions{
			URL:                uri,
			DisableAudio:       true,
			DialTimeout:        3 * time.Second,
			ReadWriteTimeout:   5 * time.Second,
			Debug:              false,
			OutgoingProxy:      false,
			InsecureSkipVerify: true,
		},
	)
	if err != nil {
		err = errors.Errorf("dial rtsp uri %s failed", err)
		return
	}
	return
}
