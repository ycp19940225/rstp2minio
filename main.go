package main

import (
	"os"
	minio2 "rstp2minio/minio"
	"rstp2minio/steam"
	monitorstorage "rstp2minio/steam/monitor_storage"
)

func main() {
	// 设置 TZ 环境变量
	os.Setenv("TZ", "Asia/Shanghai")

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
		IP: "rtsp://localhost:8554/stream",
		ID: 8950183874114224129,
		No: "test1",
	}
	streamsData = append(streamsData, testStreamsData)
	server.SaveMonitorStream(streamsData)
	for true {

	}
}
