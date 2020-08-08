package main

type (
	TaskAddRequest struct {
		Address string `json:"address" validate:"required,url"`
	}

	Config struct {
		HttpProxy         bool   `yaml:"http_proxy"`
		HttpProxyAddress  string `yaml:"http_proxy_address"`
		DownloadDirectory string `yaml:"download_path"`
	}
)

var (
	DefaultConfig Config = Config{
		HttpProxy:         false,
		DownloadDirectory: "./download",
	}
)
