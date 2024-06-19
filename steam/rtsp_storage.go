package steam

import (
	"context"
	"fmt"
	minioclient "rstp2minio/steam/minio_client"
	monitorstorage "rstp2minio/steam/monitor_storage"
	"time"

	"github.com/deepch/vdk/av"
	"github.com/deepch/vdk/format/mp4"
	"github.com/golang-module/carbon/v2"
	"rstp2minio/minio"
)

type Camera struct {
	No        string          `gorm:"column:no;not null;comment:camera no" json:"no"`
	IP        string          `gorm:"column:ip;not null;comment:camera ip address" json:"ip"`
	Position  string          `gorm:"column:position;not null;comment:位置信息" json:"position"`
	Status    int32           `gorm:"column:status;not null;comment:status 0 default 1normal 9 broken" json:"status"`
	NetStatus int32           `gorm:"column:net_status;not null;comment:net_status 0 default 1ing 9 broken" json:"netStatus"`
	CreatedAt carbon.DateTime `gorm:"column:created_at;comment:create time" json:"createdAt"`
	UpdatedAt carbon.DateTime `gorm:"column:updated_at;comment:update time" json:"updatedAt"`
	DeletedAt carbon.DateTime `gorm:"column:deleted_at;comment:delete time" json:"deletedAt"`
	ID        int
}

type CreateVideo struct {
	CameraId    int64     `json:"cameraId,omitempty"`
	VideoNo     *string   `json:"videoNo,omitempty"`
	FileName    *string   `json:"fileName,omitempty"`
	FileSize    uint64    `json:"fileSize,string,omitempty"`
	FileType    uint32    `json:"fileType,omitempty"`
	FileRealURL *string   `json:"fileRealUrl,omitempty"`
	StartDate   time.Time `json:"startDate,omitempty"`
	EndDate     time.Time `json:"endDate,omitempty"`
}

type MonitorStorageConfig struct {
	Duration       string `protobuf:"bytes,1,opt,name=duration,proto3" json:"duration,omitempty"`
	SavePath       string `protobuf:"bytes,3,opt,name=savePath,proto3" json:"savePath,omitempty"`
	NoVideoTimeout int64  `protobuf:"varint,4,opt,name=noVideoTimeout,proto3" json:"noVideoTimeout,omitempty"`
}

type fileType int32

const (
	FileTypeMp4 = iota + 1
)

func (m fileType) String() string {
	switch m {
	case FileTypeMp4:
		return "mp4"
	}
	return "unknown"
}

type RtspStorageUseCase struct {
	minio          *minio.Minio
	monitorStorage *monitorstorage.MonitorStorage
	config         MonitorStorageConfig
}

// nolint: revive
func NewRtspSteamStorageUseCase(minio *minio.Minio, monitorStorage *monitorstorage.MonitorStorage, config MonitorStorageConfig) *RtspStorageUseCase {
	// prefix rule = project name + _ + business name, example: layout_rtsp
	return &RtspStorageUseCase{
		minio:          minio,
		monitorStorage: monitorStorage,
		config:         config,
	}
}

func (c *RtspStorageUseCase) SaveMonitorStream(streamsData []*Camera) {
	ctx := context.Background()
	// 保存每个摄像头的streams
	c.serveStreams(ctx, streamsData)
}

func (c *RtspStorageUseCase) serveStreams(ctx context.Context, streamsData []*Camera) {
	for _, stream := range streamsData {
		go c.HandleRTSPWorker(ctx, stream)
		// 平滑启动
		time.Sleep(1 * time.Second)
	}
}

func (c *RtspStorageUseCase) HandleRTSPWorker(ctx context.Context, camera *Camera) {
	fmt.Println("start save :", camera)
	err := c.RTSPWorkerLoopV2(ctx, camera)
	if err != nil {
		fmt.Println("err:", err)
	}
}

func (c *RtspStorageUseCase) RTSPWorkerLoopV2(ctx context.Context, camera *Camera) (err error) {
	rootPath := c.config.SavePath
	noVideoTimeout := c.config.NoVideoTimeout
	noVideoTimeoutFormat := time.Duration(noVideoTimeout) * time.Second
	name := fmt.Sprintf("%d", camera.ID)

	stream, err := c.monitorStorage.AddStream(camera.IP, name, monitorstorage.MSE)
	if err != nil {
		return ErrorInitClientErr
	}
	var tempPacket *av.Packet
	var tempDuration time.Duration
	for {
		currentClient := stream.RTSPClient
		codecs := stream.CodecData
		if len(codecs) == 0 {
			return ErrorStreamNoCodec
		}
		nowTime := carbon.Now()
		date := nowTime.ToDateString()
		videoInfo := CreateVideo{}
		videoInfo.CameraId = int64(camera.ID)
		videoInfo.FileType = FileTypeMp4
		fileName := nowTime.Format("H_i_s")
		videoInfo.StartDate = carbon.Now().ToStdTime()

		f := minioclient.NewMinIOObject(c.minio)
		muxer := mp4.NewMuxer(f)
		err = muxer.WriteHeader(codecs)
		if err != nil {
			return err
		}
		var loop = true
		noVideo := time.NewTimer(noVideoTimeoutFormat)
		dur, err := time.ParseDuration(c.config.Duration)
		saveLimit := time.NewTimer(dur)
		// 寻找关键帧
		checkLimit := false
		for loop {
			select {
			case <-stream.StopSignal:
				return ErrorVideoStopByDeleted
			case reconnectIp := <-stream.UpdateSignal:
				loop = false
				stream, err = c.monitorStorage.ReconnectStream(reconnectIp, name)
				if err != nil {
					return err
				}
			case <-noVideo.C:
				loop = false
				stream, err = c.reconnectTry(ctx, camera, name)
				if err != nil {
					return err
				}
			case <-saveLimit.C:
				checkLimit = true
				saveLimit.Stop()
			case packetAV := <-currentClient.OutgoingPacketQueue:
				if tempPacket != nil {
					if err = muxer.WritePacket(*tempPacket); err != nil {
						return err
					}
					dur += tempDuration
					tempPacket = nil
				}
				if checkLimit {
					if packetAV.IsKeyFrame {
						tempPacket = packetAV
						tempDuration = packetAV.Duration
						loop = false
						continue
					}
					dur += packetAV.Duration
				}
				if packetAV.IsKeyFrame {
					noVideo.Reset(noVideoTimeoutFormat)
				}
				if err = muxer.WritePacket(*packetAV); err != nil {
					return err
				}
				//fmt.Println("av:", packetAV)
			}
		}
		fileName += "-" + nowTime.AddDuration(dur.String()).Format("H_i_s")
		filePath := fmt.Sprintf("%s/%s/%s/%s.mp4", rootPath, name, date, fileName)
		videoInfo.EndDate = nowTime.AddDuration(dur.String()).ToStdTime()
		fileNameURL := fmt.Sprintf("%s-%s.mp4", date, fileName)
		videoInfo.FileName = &fileNameURL
		videoInfo.FileRealURL = &filePath
		videoInfo.FileSize = uint64(f.CurrentPosition)
		err = muxer.WriteTrailer()
		if err != nil {
			return err
		}
		go func() {
			f.ObjectName = filePath
			err = f.Close()
			if err != nil {
				fmt.Println("minio 保存失败，转为本地保存")
				// 如果mimio报出错的情况下，保存到本地
				err := f.LocalStorage()
				if err != nil {

				}
			}
		}()
	}
}

func (c *RtspStorageUseCase) reconnectTry(ctx context.Context, camera *Camera, name string) (*monitorstorage.Stream, error) {
	tryTime := 15
	for tryTime > 0 {
		stream, err := c.monitorStorage.ReconnectStream(camera.IP, name)
		if err != nil {
			time.Sleep(1)
			tryTime--
			continue
		}
		return stream, nil
	}
	c.monitorStorage.DeleteStream(name)
	return nil, ErrorNoVideoErr
}
