package config

// LivenessConfig ...
// Configuration used to check liveness for HTTPServer
//
// Note: All env variables are prefixed with LIVENESS_ (see httpServerConfig.go)
type LivenessConfig struct {
	MaxGoRoutines int `env:"MAX_GO_ROUTINES,default=100"`
}

// ReadinessConfig ...
// Configuration used to check readiness for HTTPServer
//
// Note: All env variables are prefixed with READINESS_ (see httpServerConfig.go)
type ReadinessConfig struct {
}
