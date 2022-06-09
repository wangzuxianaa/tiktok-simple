package conf

import (
	"github.com/spf13/viper"
	"log"
)

type Configuration struct {
	MysqlConfig struct {
		User     string `yaml:"user"`
		Password string `yaml:"password"`
		Host     string `yaml:"host"`
		Port     int    `yaml:"port"`
		Name     string `yaml:"name"`
	} `yaml:"mysqlConfig"`
	RedisConfig struct {
		Addr     string `yaml:"addr"`
		Password string `yaml:"password"`
		DB       int    `yaml:"db"`
	} `yaml:"redisConfig"`
	VideoAddr string `yaml:"videoAddr"`
	CoverAddr string `yaml:"coverAddr"`
}

var Conf *Configuration

func Config() error {
	if Conf != nil {
		return nil
	}

	var err error
	viper.SetConfigName("conf")
	viper.AddConfigPath("./conf")
	viper.SetConfigType("yaml")

	// 读取配置文件
	if err = viper.ReadInConfig(); err != nil {
		log.Printf("读取配置文件失败: %v", err)
		return err
	}

	if err = viper.Unmarshal(&Conf); err != nil {
		log.Printf("配置失败：%v", err)
		return err
	}

	return nil
}
