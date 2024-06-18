package monitorstorage

const (
	MSE Mode = iota
	WEBRTC
)

type Mode uint8

func (m Mode) String() string {
	switch m {
	case MSE:
		return "mse"
	case WEBRTC:
		return "webrtc"
	}
	return "unknown"
}

type StreamOptions struct {
	uri  string
	mode Mode
}

func WithStreamUri(uri string) func(*StreamOptions) {
	return func(options *StreamOptions) {
		getStreamOptionsOrSetDefault(options).uri = uri
	}
}

func WithStreamMode(mode Mode) func(*StreamOptions) {
	return func(options *StreamOptions) {
		getStreamOptionsOrSetDefault(options).mode = mode
	}
}

func getStreamOptionsOrSetDefault(options *StreamOptions) *StreamOptions {
	if options == nil {
		return &StreamOptions{}
	}
	return options
}
