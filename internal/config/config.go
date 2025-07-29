package config

import (
	"os"
	"strconv"
)

type Config struct {
	App      AppConfig
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	IPFS     IPFSConfig
	Ethereum EthereumConfig
}

type AppConfig struct {
	Environment string
}

type ServerConfig struct {
	Port    string
	GinMode string
}

type DatabaseConfig struct {
	URL string
}

type JWTConfig struct {
	Secret    string
	ExpiresIn string
}

type IPFSConfig struct {
	APIURL string
}

type EthereumConfig struct {
	RPCURL             string
	PrivateKey         string
	PaperContractAddr  string
	ReviewContractAddr string
}

func LoadConfig() *Config {
	return &Config{
		App: AppConfig{
			Environment: getEnv("ENVIRONMENT", "development"),
		},
		Server: ServerConfig{
			Port:    getEnv("PORT", "8080"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			URL: getEnv("DATABASE_URL", "postgres://user:password@localhost:5432/nft_platform?sslmode=disable"),
		},
		JWT: JWTConfig{
			Secret:    getEnv("JWT_SECRET", "your_jwt_secret_here"),
			ExpiresIn: getEnv("JWT_EXPIRES_IN", "24h"),
		},
		IPFS: IPFSConfig{
			APIURL: getEnv("IPFS_API_URL", "http://localhost:5001"),
		},
		Ethereum: EthereumConfig{
			RPCURL:             getEnv("ETHEREUM_RPC_URL", "http://localhost:8545"),
			PrivateKey:         getEnv("PRIVATE_KEY", ""),
			PaperContractAddr:  getEnv("CONTRACT_ADDRESS_PAPER", ""),
			ReviewContractAddr: getEnv("CONTRACT_ADDRESS_REVIEW", ""),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
