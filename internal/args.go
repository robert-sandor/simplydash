package internal

import "github.com/alecthomas/kong"

type Args struct {
	Host          string `name:"host"        default:"0.0.0.0" help:"host to listen on"           short:"l"`
	Port          string `name:"port"        default:"8080"    help:"port to listen on"           short:"p"`
	ConfigFile    string `name:"config"      default:"./config/config.yml"                        help:"Path to config file"         short:"c"`
	ImageCacheDir string `name:"image-cache" default:"./images"                                   help:"Path to dir to store images" short:"i"`
	Log           struct {
		Level string `name:"level" default:"info" help:"log level" enum:"trace,debug,info,warn,error,fatal,panic"`
		Type  string `name:"type" default:"text" help:"log type" enum:"text,json"`
	} `embed:"" prefix:"log-"`
	AccessLogs bool `name:"access-logs" default:"false" help:"enable access logs" type:"boolean"`
}

func GetArgs() Args {
	args := Args{}
	kong.Parse(&args)
	return args
}
