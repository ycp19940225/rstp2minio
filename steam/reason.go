package steam

import "github.com/pkg/errors"

const (
	IdempotentMissingToken = "idempotent.token.missing"

	TooManyRequests = "too.many.requests"
	DataNotChange   = "data.not.change"
	DuplicateField  = "duplicate.field"
	RecordNotFound  = "record.not.found"
	NoPermission    = "no.permission"

	ErrorNoVideoErrStr             = "视频流缺失"
	ErrorInitConfigErrStr          = "stream init config err"
	ErrorInitClientErrStr          = "网络中断"
	ErrorStreamNoCodecStr          = "stream no codec"
	ErrorStreamSaveErrStr          = "stream minio save err"
	ErrorStreamStorageRecordErrStr = "stream minio save record err"
	ErrorVideoCompressErrStr       = "video compress err"
	ErrorVideoStopByDeletedStr     = "video stop by deleted"
)

// Default stream errors
var (
	ErrorNoVideoErr             = errors.New(ErrorNoVideoErrStr)
	ErrorInitConfigErr          = errors.New(ErrorInitConfigErrStr)
	ErrorInitClientErr          = errors.New(ErrorInitClientErrStr)
	ErrorStreamNoCodec          = errors.New(ErrorStreamNoCodecStr)
	ErrorStreamSaveErr          = errors.New(ErrorStreamSaveErrStr)
	ErrorStreamStorageRecordErr = errors.New(ErrorStreamStorageRecordErrStr)
	ErrorVideoCompressErr       = errors.New(ErrorVideoCompressErrStr)
	ErrorVideoStopByDeleted     = errors.New(ErrorVideoStopByDeletedStr)
)
