package util

import (
	"github.com/aws/aws-xray-sdk-go/xray"
)

var (
	initialized = false
)

func InitXray() {
	if initialized {
		return
	}

	initialized = true
	if err := xray.Configure(xray.Config{
		DaemonAddr:     "127.0.0.1:2000", // default
		ServiceVersion: "1.2.3",
	}); err != nil {
		panic(err)
	}
}
