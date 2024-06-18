# rtsp2minio
摄像头视频流分段保存到minio、本地为MP4格式文件

## useage

```go

package main

import (
	minio2 "rstp2minio/minio"
	"rstp2minio/steam"
	monitorstorage "rstp2minio/steam/monitor_storage"
)

func main() {

	minio, _ := minio2.New(
		minio2.WithEndpoint(""),
		minio2.WithKey(""),
		minio2.WithSecret(""),
		minio2.WithBucket("test"),
	)
	steamClient := monitorstorage.NewMonitorStorage()

	// 存储配置
	config := steam.MonitorStorageConfig{
		Duration:       "1m",    //存储的视频分割时长，也就是每分钟生成一个一分钟时长的mp4文件, 1h 一小时
		SavePath:       "video", // 保存地址
		NoVideoTimeout: 120,     // 无视频流的时候，超时设置 120s
	}
	server := steam.NewRtspSteamStorageUseCase(minio, steamClient, config)
	// 摄像头信息
	var streamsData []*steam.Camera
	testStreamsData := &steam.Camera{
		IP: "rtsp://wowzaec2demo.streamlock.net/vod/mp4:BigBuckBunny_115k.mov",
		ID: 8950183874114224129,
		No: "test1",
	}
	streamsData = append(streamsData, testStreamsData)
	server.SaveMonitorStream(streamsData)
}

```