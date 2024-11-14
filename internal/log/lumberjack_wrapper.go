package log

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/natefinch/lumberjack"
	"github.com/smallnest/ringbuffer"
)

const padding = 256

// LumberjackWrapper ...
type LumberjackWrapper struct {
	*lumberjack.Logger
	buffer *ringbuffer.RingBuffer
}

// LumberjackWrapperConfig ...
type LumberjackWrapperConfig struct {
	Path       string // 日志本地存储目录
	MaxSize    int    // 单个日志文件大小，单位：MB。
	MaxBackups int    // 最多存储日志文件数目
	MaxAge     int    // 保留旧日志文件最大天数
	BufferSize int    // 日志记录写到磁盘缓存区大小，单位：字节。0 表示立即写入磁盘
	IsCompress bool   // 日志文件是否压缩，true gzip压缩、false 不压缩。
	FileClose  bool
	StdClose   bool
}

// NewLumberjackWrapper ...
func NewLumberjackWrapper(cfg *LumberjackWrapperConfig) *LumberjackWrapper {
	l := LumberjackWrapper{Logger: &lumberjack.Logger{
		Filename:   filepath.Join(cfg.Path, "info.log"),
		MaxSize:    cfg.MaxSize,
		MaxBackups: cfg.MaxBackups,
		MaxAge:     cfg.MaxAge,
		Compress:   cfg.IsCompress,
	}}

	if cfg.BufferSize == 0 {
		return &l
	}
	l.buffer = ringbuffer.New(cfg.BufferSize + padding)

	return &l
}

func (l *LumberjackWrapper) Write(p []byte) (n int, err error) {
	if l.buffer == nil {
		return l.Logger.Write(p)
	}

	if l.buffer.Free() < len(p) {
		l.Sync()
	}

	if l.buffer.IsEmpty() && l.buffer.Free() < len(p) {
		return l.Logger.Write(p)
	}

	for {
		n, err = l.buffer.Write(p)
		if err != nil {
			l.Logger.Write([]byte(fmt.Sprintf("LumberjackWrapper Write err:%v", err)))
			return n, err
		}

		if n == len(p) {
			break
		}

		p = p[n:]

		l.Sync()
	}

	if l.buffer.Free() < padding {
		return l.Sync()
	}

	return n, nil
}

// Sync ...
func (l *LumberjackWrapper) Sync() (n int, err error) {
	bs := l.buffer.Bytes(nil)
	if _, err := l.Logger.Write(bs); err != nil {
		log.Printf("LumberjackWrapper Sync err:%v", err)
		return 0, err
	}

	l.buffer.Reset()
	return len(bs), nil
}
