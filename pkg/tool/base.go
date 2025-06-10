package tool

type Tool interface {
	Name() string
	Description() string
	InputSchema() map[string]string
	Run(params map[string]string) (string, error)
}
