package swagger

type ApiPath struct {
	Path           string
	ResourceType   string
	ApiVersion     string
	ExampleMap     map[string]string
	OperationIdMap map[string]string
	Methods        []string
	ApiType        ApiType
}

type ApiType string

const (
	ApiTypeUnknown        ApiType = "unknown"
	ApiTypeList           ApiType = "list"
	ApiTypeResource       ApiType = "resource"
	ApiTypeResourceAction ApiType = "resourceAction"
	ApiTypeProviderAction ApiType = "providerAction"
)
