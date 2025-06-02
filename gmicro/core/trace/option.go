package trace

type Options struct {
	Name     string  `json:"name,omitempty"`
	Endpoint string  `json:"endpoint,omitempty"`
	Sampler  float64 `json:"sampler,omitempty"`
	Batcher  string  `json:"batcher,omitempty"`
}
