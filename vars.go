package xcli

type ServiceConfig struct {
	ServiceName        string
	ServiceDisplayName string
	ServiceDescription string
	EnvVars            map[string]string
}
