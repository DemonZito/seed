package forms

import (
	"qlova.org/seed"
	"qlova.org/seed/script"
)

func init() {
	script.RegisterRenderer(func(c seed.Seed) []byte {
		return []byte(`s.form = {};

		s.form.reportValidity = function(el) {
			while(!el.reportValidity) {
				el = el.parentElement;
				if (!el) {
					return true;
				}
			}
			return el.reportValidity();
		};
		
		`)
	})
}
