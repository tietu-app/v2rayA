package conf

import (
	"fmt"
	"github.com/stevenroose/gonfig"
	"github.com/v2rayA/v2rayA/common"
	"github.com/v2rayA/v2rayA/pkg/util/log"
	log2 "log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type Params struct {
	Address                 string   `id:"address" short:"a" default:"0.0.0.0:2017" desc:"Listening address"`
	Config                  string   `id:"config" short:"c" desc:"v2rayA configuration directory"`
	V2rayBin                string   `id:"v2ray-bin" desc:"Executable v2ray binary path. Auto-detect if put it empty."`
	V2rayConfigDirectory    string   `id:"v2ray-confdir" desc:"Additional v2ray config directory, files in it will be combined with config generated by v2rayA"`
	V2rayAssetsDirectory    string   `id:"v2ray-assetsdir" desc:"v2ray-core assets directory for searching and downloading files like geoip.dat. This will override environment V2RAY_LOCATION_ASSET and XRAY_LOCATION_ASSET."`
	TransparentHook         string   `id:"transparent-hook" desc:"the executable file to run before the transparent proxy starting. v2rayA will pass in the --transparent-type (tproxy, redirect) and --stage (pre-start, post-start, pre-stop, post-stop) arguments."`
	WebDir                  string   `id:"webdir" desc:"v2rayA web files directory. use embedded files if not specify."`
	VlessGrpcInboundCertKey []string `id:"vless-grpc-inbound-cert-key" desc:"Specify the certification path instead of automatically generating a self-signed certificate. Example: /etc/v2raya/grpc_certificate.crt,/etc/v2raya/grpc_private.key"`
	IPV6Support             string   `id:"ipv6-support" default:"auto" desc:"Optional values: auto, on, off. Make sure your IPv6 network works fine before you turn it on."`
	PassCheckRoot           bool     `desc:"Skip privilege checking. Use it only when you cannot start v2raya but confirm you have root privilege"`
	ResetPassword           bool     `id:"reset-password"`
	LogLevel                string   `id:"log-level" default:"info" desc:"Optional values: trace, debug, info, warn or error"`
	LogFile                 string   `id:"log-file" desc:"The path of log file"`
	LogMaxDays              int64    `id:"log-max-days" default:"3" desc:"Maximum number of days to keep log files"`
	LogDisableColor         bool     `id:"log-disable-color"`
	LogDisableTimestamp     bool     `id:"log-disable-timestamp" desc:"Intended for use with systemd/journald to avoid duplicate timestamps in logs. This flag is ignored when using the --log-file flag or the V2RAYA_LOG_FILE environment variable."`
	Lite                    bool     `id:"lite" desc:"Lite mode for non-root and non-linux users"`
	ShowVersion             bool     `id:"version"`
	PrintReport             string   `id:"report" desc:"Print report"`
}

var params Params

func initFunc() {
	err := gonfig.Load(&params, gonfig.Conf{
		FileDisable:       true,
		FlagIgnoreUnknown: false,
		EnvPrefix:         "V2RAYA_",
	})
	if err != nil {
		if err.Error() != "unexpected word while parsing flags: '-test.v'" {
			log2.Fatal(err)
		}
	}
	if params.ShowVersion {
		fmt.Println(Version)
		os.Exit(0)
	}
	if params.Lite {
		params.PassCheckRoot = true
	}
	if params.Config == "" {
		if params.Lite {
			params.Config = "$HOME/.config/v2raya"
		} else {
			params.Config = "/etc/v2raya"
		}
	}
	// replace all dots of the filename with underlines
	params.Config = filepath.Join(
		filepath.Dir(params.Config),
		strings.ReplaceAll(filepath.Base(params.Config), ".", "_"),
	)
	// expand '~' with user home
	params.Config, err = common.HomeExpand(params.Config)
	if err != nil {
		log2.Fatal(err)
	}
	if strings.Contains(params.Config, "$HOME") {
		if h, err := os.UserHomeDir(); err == nil {
			params.Config = strings.ReplaceAll(params.Config, "$HOME", h)
		}
	}
	if _, err := os.Stat(params.Config); os.IsNotExist(err) {
		_ = os.MkdirAll(params.Config, os.ModeDir|0750)
	} else if err != nil {
		log.Warn("%v", err)
	}
	logWay := "console"
	if params.LogFile != "" {
		logWay = "file"
		_ = os.MkdirAll(filepath.Dir(params.LogFile), os.ModeDir|0750)
	}
	log.InitLog(logWay, params.LogFile, params.LogLevel, params.LogMaxDays, params.LogDisableColor, params.LogDisableTimestamp)

	// V2rayAssetsDirectory
	if params.V2rayAssetsDirectory != "" {
		if err = os.Setenv("V2RAY_LOCATION_ASSET", params.V2rayAssetsDirectory); err != nil {
			log.Fatal("failed to set V2rayAssetsDirectory: %v", err)
		}
		if err = os.Setenv("XRAY_LOCATION_ASSET", params.V2rayAssetsDirectory); err != nil {
			log.Fatal("failed to set V2rayAssetsDirectory: %v", err)
		}
	}
}

var once sync.Once

func GetEnvironmentConfig() *Params {
	once.Do(initFunc)
	return &params
}

func SetConfig(config Params) {
	params = config
}
