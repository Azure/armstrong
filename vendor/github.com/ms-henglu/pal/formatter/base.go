package formatter

import (
	"github.com/ms-henglu/pal/types"
)

type Formatter interface {
	Format(r types.RequestTrace) string
}
