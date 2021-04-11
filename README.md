viper consul remote 
----


使用方法
```
v := viper.New()
	v.AddSecureRemoteProvider("consul", "dev.env.devops.nx.cn:8500", "jzh/dev/log_test/test.yaml", "feec61ec-8d6f-c59e-7514-eb62e6635b6a")
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