viper consul remote 
----

viper 使用了 `https://github.com/bketelsen/crypt` 库作为远程配置中心。对 `etcd`的支持比较好，`consul` 的支持并不完整。

`viper-consul-remote`适用于使用 `consul` 作为配置中心，并使用 `viper` 作为配置库。 

## example
引入 remote 库
```
_ "git.nxdev.cn/pkg/viper-consul-remote/remote"
```

```
	v := viper.New()
	v.AddSecureRemoteProvider("consul", "<host>:8500", "<consul_path>", "<token>")
	v.SetConfigType("yaml") // Need to explicitly set this to json

	err := v.ReadRemoteConfig()
	if err != nil {
		log.Println("error:", err)
	}
	v.WatchRemoteConfigOnChannel() // viper 的全局没有实现这个方法。

	for {
		fmt.Println(v.Get("languages")) // 8080
		fmt.Println(v.Get("hostname"))  // myhostname.com
		time.Sleep(1 * time.Second)
	}
```