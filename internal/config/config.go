package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

const configFileName = "spotifycli.json"

type Config struct {
	ClientID     string `json:"client_id"`
	RedirectPath string `json:"redirect_path"`
	Port         string `json:"port"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	TokenType    string `json:"token_type"`
	TokenExpiry  int64  `json:"token_expiry"`
	LastSaved    int64  `json:"last_saved"`
}

func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", configFileName), nil
}

func LoadConfig() (*Config, error) {
	path, err := getConfigPath()
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		// Default config
		return &Config{
			RedirectPath: "callback",
		}, nil
	}

	data, err := os.ReadFile(path)
	if err != nil && err != io.EOF {
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	// Decrypt tokens if present
	key := os.Getenv("SPOTIFYCLI_KEY")
	if key == "" {
		fmt.Println(key)
	}

	if cfg.AccessToken != "" {
		plain, err := DecryptToken(cfg.AccessToken, key)
		if err == nil {
			cfg.AccessToken = plain
		}
	}

	if cfg.RefreshToken != "" {
		plain, err := DecryptToken(cfg.RefreshToken, key)
		if err == nil {
			cfg.RefreshToken = plain
		}
	}

	return &cfg, nil
}

func (c *Config) Save() error {
	path, err := getConfigPath()
	if err != nil {
		return err
	}

	key := os.Getenv("SPOTIFYCLI_KEY")
	if key == "" {
		key = generateEncryptionKey()
	}

	if err := os.Setenv("SPOTIFYCLI_KEY", key); err != nil {
		return err
	}

	encCfg := *c
	if encCfg.AccessToken != "" {
		enc, err := EncryptToken(encCfg.AccessToken, key)
		if err == nil {
			encCfg.AccessToken = enc
		}
	}

	if encCfg.RefreshToken != "" {
		enc, err := EncryptToken(encCfg.RefreshToken, key)
		if err == nil {
			encCfg.RefreshToken = enc
		}
	}

	encCfg.LastSaved = time.Now().Unix()
	data, err := json.MarshalIndent(encCfg, "", "  ")
	if err != nil {
		return err
	}

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0700); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

func (c *Config) IsAuthenticated() bool {
	return c.AccessToken != "" && c.RefreshToken != ""
}

func (c *Config) IsTokenExpired() bool {
	if c.TokenExpiry == 0 {
		return true
	}

	return time.Now().Unix() > c.TokenExpiry
}

func (c *Config) SetTokens(access, refresh, tokenType string, expiresIn int64) {
	c.AccessToken = access
	c.RefreshToken = refresh
	c.TokenType = tokenType
	c.TokenExpiry = time.Now().Unix() + expiresIn
}

func (c *Config) ClearTokens() {
	c.AccessToken = ""
	c.RefreshToken = ""
	c.TokenType = ""
	c.TokenExpiry = 0
}

func (c *Config) GetAccessToken() string {
	return c.AccessToken
}

func (c *Config) GetRefreshToken() string {
	return c.RefreshToken
}

func (c *Config) GetTokenType() string {
	return c.TokenType
}
