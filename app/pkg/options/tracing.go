package options

import (
	"github.com/spf13/pflag"
	"minifast/pkg/errors"
)

type TelemetryOptions struct {
	Name     string  `mapstructure:"name" json:"name,omitempty"`
	Endpoint string  `mapstructure:"endpoint" json:"endpoint,omitempty"`
	Sampler  float64 `mapstructure:"sampler" json:"sampler,omitempty"`
	Batcher  string  `mapstructure:"batcher" json:"batcher,omitempty"`
}

func NewTelemetryOptions() *TelemetryOptions {
	return &TelemetryOptions{
		Name:     "minifast",
		Endpoint: "http://127.0.0.1:14268/api/traces",
		Sampler:  1.0,
		Batcher:  "jaeger",
	}
}

func (t *TelemetryOptions) Validate() []error {
	errs := []error{}
	if t.Batcher != "jaeger" && t.Batcher != "zipkin" {
		errs = append(errs, errors.New("opentelemetry batcher only support jaeger or zipkin"))
	}
	return errs
}

// AddFlags adds flags related to open telemetry for a specific tracing to the specified FlagSet.
func (t *TelemetryOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringVar(&t.Name, "telemetry.name", t.Name, "opentelemetry name")
	fs.StringVar(&t.Endpoint, "telemetry.endpoint", t.Endpoint, "opentelemetry endpoint")
	fs.Float64Var(&t.Sampler, "telemetry.sampler", t.Sampler, "telemetry sampler")
	fs.StringVar(&t.Batcher, "telemetry.batcher", t.Batcher, "telemetry batcher, only support jaeger and zipkin")
}
