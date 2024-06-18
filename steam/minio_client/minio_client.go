package minioclient

import (
	"bufio"
	"bytes"
	"context"
	"io"
	"os"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	zenwellminio "rstp2minio/minio"
)

type MinIOObject struct {
	client          *zenwellminio.Minio
	ObjectName      string
	buffer          bytes.Buffer
	CurrentPosition int64
}

func NewMinIOObject(client *zenwellminio.Minio) *MinIOObject {
	return &MinIOObject{
		client:          client,
		buffer:          bytes.Buffer{}, // 初始化 buffer
		CurrentPosition: 0,              // 初始化当前位置
	}
}
func (m *MinIOObject) Write(p []byte) (n int, err error) {
	if m.CurrentPosition < int64(m.buffer.Len()) {
		// 如果当前位置在缓冲区内，写入覆盖数据
		copy(m.buffer.Bytes()[m.CurrentPosition:], p)
		n = len(p)
	} else {
		// 如果当前位置超出缓冲区，直接追加数据到缓冲区末尾
		n, err = m.buffer.Write(p)
	}
	m.CurrentPosition += int64(n)
	return n, err
}

func (m *MinIOObject) Seek(offset int64, whence int) (int64, error) {
	var newPosition int64

	switch whence {
	case io.SeekStart:
		newPosition = offset
	case io.SeekCurrent:
		newPosition = m.CurrentPosition + offset
	case io.SeekEnd:
		newPosition = int64(m.buffer.Len()) + offset
	}

	// 防止位置超出缓冲区的范围
	if newPosition < 0 {
		newPosition = 0
	} else if newPosition > int64(m.buffer.Len()) {
		newPosition = int64(m.buffer.Len())
	}

	m.CurrentPosition = newPosition
	return newPosition, nil
}

func (m *MinIOObject) Close() error {
	subCtx := context.Background()
	// 设置上传选项
	options := minio.PutObjectOptions{
		ContentType: "video/mp4", // 设置 Content-Type
	}
	_, err := m.client.PutObject(subCtx, m.ObjectName, &m.buffer, int64(m.buffer.Len()), options)
	return err
}

func (m *MinIOObject) LocalStorage() (err error) {
	dir := filepath.Dir(m.ObjectName)
	err = os.MkdirAll(dir, os.ModePerm)
	if err != nil {
		return
	}
	// 打开文件以用于写入
	file, err := os.Create(m.ObjectName)
	if err != nil {
		return
	}
	defer file.Close()
	// 使用 bufio.Writer 进行文件缓冲，以提高效率
	writer := bufio.NewWriter(file)

	// 将 bytes.Buffer 的内容写入文件缓冲
	_, err = m.buffer.WriteTo(writer)
	if err != nil {
		return
	}
	// 刷新文件缓冲并检查是否有错误
	err = writer.Flush()
	if err != nil {
		return
	}
	// 再次尝试上传
	subCtx := context.Background()
	_, err = m.client.FPutObject(subCtx, m.ObjectName, m.ObjectName)
	if err != nil {
		return err
	}
	// 如果成功 移除本地文件
	err = os.Remove(m.ObjectName)
	if err != nil {
		return
	}
	return err
}
