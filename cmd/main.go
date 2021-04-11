package main

import (
	"fmt"
	"log"
	"time"

	"github.com/spf13/viper"

	_ "git.nxdev.cn/pkg/viper-consul-remote/remote"
)

func main() {
	v := viper.New()
	v.AddSecureRemoteProvider("consul", "<host>:8500", "jzh/dev/log_test/test.yaml", "<token>")
	v.SetConfigType("yaml") // Need to explicitly set this to json

	err := v.ReadRemoteConfig()
	if err != nil {
		log.Println("error:", err)
	}
	v.WatchRemoteConfigOnChannel()

	for {
		fmt.Println(v.Get("languages")) // 8080
		fmt.Println(v.Get("hostname"))  // myhostname.com
		time.Sleep(1 * time.Second)
	}

}
