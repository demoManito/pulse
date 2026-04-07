package logger

// Config holds logger configuration.
type Config struct {
	// Driver specifies the logging backend: "logrus" or "zap". Default: "logrus".
	Driver string `yaml:"driver"`
	// Level specifies the minimum log level. Default: InfoLevel.
	Level Level `yaml:"level"`
}
