package config

type ObservabilityConfig struct {
	MeterName     string `mapstructure:"meter_name"`
	TracerName    string `mapstructure:"tracer_name"`
	MeterExporter string `mapstructure:"meter_exporter"`
	TraceExporter string `mapstructure:"trace_exporter"`
	LogExporter   string `mapstructure:"logger_exporter"`
}
