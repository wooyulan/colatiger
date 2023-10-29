package config

type Log struct {
	Level        string `mapstructure:"level" json:"level" yaml:"level"`
	RootDir      string `mapstructure:"root_dir" json:"root_dir" yaml:"root_dir"`
	Filename     string `mapstructure:"filename" json:"filename" yaml:"filename"`
	MaxBackups   int    `mapstructure:"max_backups" json:"max_backups" yaml:"max_backups"`
	MaxSize      int    `mapstructure:"max_size" json:"max_size" yaml:"max_size"` // MB
	MaxAge       int    `mapstructure:"max_age" json:"max_age" yaml:"max_age"`    // 保存的天数
	Compress     bool   `mapstructure:"compress" json:"compress" yaml:"compress"`
	ShowLine     bool   `mapstructure:"show_line" json:"show_line" yaml:"show_line"` // 显示行号
	LogInConsole bool   `mapstructure:"log_in_console" json:"log_in_console" yaml:"log_in_console"`
}
