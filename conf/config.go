package config

import (
	//"github.com/rmenn/GitWebhookProxy/pkg/proxy"

	"log"
	"os"
	"strings"

	"github.com/spf13/viper"
)

type ProxyConf struct {
	FrontEndURL  string   `mapstructure:"frontEndURL"`
	Provider     string   `mapstructure:"provider"`
	UpstreamURL  string   `mapstructure:"upstreamURL"`
	AllowedPaths []string `mapstructure:"allowedPaths"`
	Secret       string   `mapstructure:"secret"`
	IgnoredUsers []string `mapstructure:"ignoredUsers"`
	AllowedUsers []string `mapstructure:"allowedUsers"`
}

func Init() (string, []ProxyConf) {
	var Proxys []ProxyConf
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	err := viper.ReadInConfig()
	if err != nil {
		log.Fatalf("unable to read config : %v", err)
	}
	err = viper.UnmarshalKey("proxies", &Proxys)
	if err != nil {
		log.Fatalf("unable to unmarshal config : %v", err)
	}

	for i := range Proxys {
		if strings.HasPrefix(Proxys[i].Secret, "$") {
			Proxys[i].Secret = os.Getenv(strings.TrimPrefix(Proxys[i].Secret, "$"))
		}
	}
	return viper.GetString("port"), Proxys
}
