package azurerm

type Mapping struct {
	ResourceType         string `json:"resourceType"`
	ExampleConfiguration string `json:"exampleConfiguration,omitempty"`
	IdPattern            string `json:"idPattern"`
}

type DependencyLoader interface {
	Load() ([]Mapping, error)
}
