package spacer

import (
	"github.com/qlova/seed"
	"github.com/qlova/seed/css"
	"github.com/qlova/seed/style"

	"github.com/qlova/seed/s/html/div"
)

//New returns a new row.
func New(space style.Unit, options ...seed.Option) seed.Seed {
	return div.New(css.SetFlexBasis(space.Unit()), css.Set("display", "flex").And(options...))
}