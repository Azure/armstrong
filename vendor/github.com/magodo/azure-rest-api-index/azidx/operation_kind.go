package azidx

import "github.com/go-openapi/spec"

type OperationKind string

const (
	OperationKindGet     OperationKind = "GET"
	OperationKindPut                   = "PUT"
	OperationKindPost                  = "POST"
	OperationKindDelete                = "DELETE"
	OperationKindOptions               = "OPTIONS"
	OperationKindHead                  = "HEAD"
	OperationKindPatch                 = "PATCH"
)

var PossibleOperationKinds = []OperationKind{
	OperationKindGet,
	OperationKindPut,
	OperationKindPost,
	OperationKindDelete,
	OperationKindOptions,
	OperationKindHead,
	OperationKindPatch,
}

func PathItemOperation(pathItem spec.PathItem, op OperationKind) *spec.Operation {
	switch op {
	case OperationKindGet:
		return pathItem.Get
	case OperationKindPut:
		return pathItem.Put
	case OperationKindPost:
		return pathItem.Post
	case OperationKindDelete:
		return pathItem.Delete
	case OperationKindOptions:
		return pathItem.Options
	case OperationKindHead:
		return pathItem.Head
	case OperationKindPatch:
		return pathItem.Patch
	}
	return nil
}
