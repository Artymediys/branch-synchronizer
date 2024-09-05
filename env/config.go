package env

import (
	"encoding/json"
	"errors"
	"os"
)

type Config struct {
	GitLab               GitLabConfig     `json:"gitlab"`
	Telegram             TelegramConfig   `json:"telegram,omitempty"`
	Mattermost           MattermostConfig `json:"mattermost,omitempty"`
	BranchPairs          []string         `json:"branch_pairs"`
	WhitelistProjects    []string         `json:"whitelist_projects,omitempty"`
	BlacklistProjects    []string         `json:"blacklist_projects,omitempty"`
	CheckIntervalInHours int              `json:"check_interval_in_hours,omitempty"`
}

type GitLabConfig struct {
	URL     string `json:"url"`
	Token   string `json:"token"`
	GroupID int    `json:"group_id"`
}

type TelegramConfig struct {
	Enabled   bool   `json:"enabled"`
	BotToken  string `json:"bot_token"`
	ChannelID int64  `json:"channel_id"`
}

type MattermostConfig struct {
	Enabled   bool   `json:"enabled"`
	Url       string `json:"url"`
	BotToken  string `json:"bot_token"`
	ChannelID string `json:"channel_id"`
}

func NewConfig(configPath string) (*Config, error) {
	if configPath == "" {
		return nil, errors.New("configuration file path is required")
	}

	var cfg Config
	err := cfg.loadConfig(configPath)
	if err != nil {
		return nil, errors.New("error loading config: " + err.Error())
	}

	err = cfg.validateConfig()
	if err != nil {
		return nil, errors.New("config validation error: " + err.Error())
	}

	return &cfg, err
}

func (cfg *Config) loadConfig(configPath string) error {
	file, err := os.ReadFile(configPath)
	if err != nil {
		return err
	}

	err = json.Unmarshal(file, cfg)
	if err != nil {
		return err
	}

	return nil
}

func (cfg *Config) validateConfig() error {
	if len(cfg.WhitelistProjects) > 0 && len(cfg.BlacklistProjects) > 0 {
		return errors.New("both whitelist and blacklist are specified; only one may be used")
	}

	if cfg.Telegram.BotToken != "" && cfg.Telegram.ChannelID != 0 {
		cfg.Telegram.Enabled = true
	} else {
		cfg.Telegram.Enabled = false
	}

	if cfg.Mattermost.Url != "" && cfg.Mattermost.BotToken != "" && cfg.Mattermost.ChannelID != "" {
		cfg.Mattermost.Enabled = true
	} else {
		cfg.Mattermost.Enabled = false
	}

	if cfg.CheckIntervalInHours <= 0 {
		cfg.CheckIntervalInHours = 24
	}

	return nil
}
