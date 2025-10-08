package config

type Config struct {
	RealmPath string
	Remote    string
	ChainID   string
	GasFee    string
	GasWanted string
	Account   string
}

func DefaultConfig() *Config {
	return &Config{
		RealmPath: "gno.land/r/example",
		Remote:    "tcp://127.0.0.1:26657",
		ChainID:   "dev",
		GasFee:    "1000000ugnot",
		GasWanted: "500000000",
		Account:   "test",
	}
}
