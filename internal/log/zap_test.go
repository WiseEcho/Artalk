package log

import (
	"fmt"
	"testing"
)

func TestZapBuffer(t *testing.T) {
	InitLogger(&LumberjackWrapperConfig{
		Path:       "/Users/yptian/log",
		MaxSize:    500,
		MaxBackups: 10,
		MaxAge:     60,
		BufferSize: 10,
		IsCompress: false,
	})
	defer Sync()

	Infof("aaaaaaa")
	Infof("asdfghjklzasdfghjklzasdfghjklzasdfghjklzaasdfasdfsdfasdfghjklzasdfghjklzasdfghjk")
}

func TestBuffer(t *testing.T) {
	buffer := NewLumberjackWrapper(&LumberjackWrapperConfig{
		Path:       "/Users/yptian/log",
		MaxSize:    500,
		MaxBackups: 10,
		MaxAge:     60,
		BufferSize: 10,
		IsCompress: false,
	})
	a := "aaaaaaa\n"
	b := "asdfghjklzasdfghjklzasdfghjklasdfghjzasdfghjklzasdfghjklzasdfghjklzasdfghjklzasdfghjklzaasdfasdfsdfasdfghjklzasdfghjklzasdfghjkasdf\n"
	fmt.Printf("len(a):%d,len(b):%d \n", len(a), len(b))
	n, err := buffer.Write([]byte(a))
	if err != nil {
		fmt.Printf("write b failed:%s", err.Error())
		return
	}
	fmt.Printf("wa:%d \n", n)
	n, err = buffer.Write([]byte(b))
	if err != nil {
		fmt.Printf("write b failed:%s", err.Error())
		return
	}
	fmt.Printf("wb:%d \n", n)

	buffer.Sync()
}
