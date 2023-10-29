package config

type Oss struct {
	Endpoint        string `mapstructure:"endpoint" json:"endpoint" yaml:"endpoint"`
	AccessKey       string `mapstructure:"access_key" json:"access_key" yaml:"access_key"`
	SecretAccessKey string `mapstructure:"secret_access_key" json:"secret_access_key" yaml:"secret_access_key"`
	UseSSL          bool   `mapstructure:"use_ssl" json:"use_ssl" yaml:"use_ssl"`
	BucketName      string `mapstructure:"bucket_name" json:"bucket_name" yaml:"bucket_name"`
}
