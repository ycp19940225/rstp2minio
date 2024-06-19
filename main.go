package main

import (
	"fmt"
	"os"
	"os/signal"
	minio2 "rstp2minio/minio"
	"rstp2minio/steam"
	monitorstorage "rstp2minio/steam/monitor_storage"
	"syscall"
)

func main() {
	//    endpoint: 'localhost:9000'
	//    bucket: 'test'
	//    key: 'admin'
	//    secret: 'admin123'
	minio, _ := minio2.New(
		minio2.WithEndpoint("localhost:9000"),
		minio2.WithKey("admin"),
		minio2.WithSecret("admin123"),
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
		IP: "rtsp://rtspstream:c32b55f44bdd4a0a711e993c8a143e3f@zephyr.rtsp.stream/pattern",
		ID: 8950183874114224129,
		No: "test1",
	}
	streamsData = append(streamsData, testStreamsData)
	server.SaveMonitorStream(streamsData)

	c := make(chan os.Signal)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	// 阻塞，直到接受到退出信号，才停止进程
	waitElegantExit(c)
}

// 优雅退出（退出信号）
func waitElegantExit(c chan os.Signal) {
	for i := range c {
		switch i {
		case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT:
			// 这里做一些清理操作或者输出相关说明，比如 断开数据库连接
			fmt.Println("receive exit signal ", i.String(), ",exit...")
			os.Exit(0)
		}
	}
}
