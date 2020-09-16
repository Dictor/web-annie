package main

type (
	// TaskAddRequest is model for "POST /tasks" endpoint.
	TaskAddRequest struct {
		Address string `json:"address" validate:"required,url"`
	}

	// Config is model of configuration for "config.yaml"
	Config struct {
		HTTPProxy         bool   `yaml:"http_proxy"`
		HTTPProxyAddress  string `yaml:"http_proxy_address"`
		DownloadDirectory string `yaml:"download_path"`
		ListenAddress     string `yaml:"listen_address"`
		IgnoreExitError   bool   `yaml:"ignore_exit_error"`
	}
)

var (
	// DefaultConfig is definition of Config's default value
	DefaultConfig Config = Config{
		HTTPProxy:         false,
		DownloadDirectory: "./download",
		ListenAddress:     ":80",
		IgnoreExitError:   false,
	}
)
