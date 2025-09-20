package config

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/BurntSushi/toml"
)

// ParseConfig 从文件路径解析 TOML 配置文件
func ParseConfig(filePath string) (*Config, error) {
	// 检查文件是否存在
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("配置文件不存在: %s", filePath)
	}

	// 读取文件内容
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %v", err)
	}

	// 解析 TOML 内容
	var config Config
	if _, err := toml.Decode(string(data), &config); err != nil {
		return nil, fmt.Errorf("解析 TOML 配置失败: %v", err)
	}

	return &config, nil
}

// ParseConfigFromString 从字符串解析 TOML 配置
func ParseConfigFromString(data string) (*Config, error) {
	var config Config
	if _, err := toml.Decode(data, &config); err != nil {
		return nil, fmt.Errorf("解析 TOML 配置失败: %v", err)
	}

	return &config, nil
}
