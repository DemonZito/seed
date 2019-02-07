package script

import (
	//Global ids.
	"encoding/base64"
	"math/big"
)

import qlova "github.com/qlova/script"
import "github.com/qlova/script/language"
import "github.com/qlova/script/language/javascript"

type setable interface {
	raw() string
	variable()
}

type Variable string

//All globals have a unique id.
var global_id int64 = 1;

func NewVariable() Variable {
	//global identification is compressed to base64 and prefixed with g_.
	var result = "g_"+base64.RawURLEncoding.EncodeToString(big.NewInt(global_id).Bytes())

	global_id++

	return Variable(result)
}

func (Variable) variable() {}
func (v Variable) raw() string {
	return string(v)
}

func (q Script) Get(name Variable) qlova.String {
	return q.wrap(`window.localStorage.getItem("`+string(name)+`")`)
}

func (q Script) Set(name setable, value qlova.Type) {
	var v = value.LanguageType()

	switch t := v.(type) {
	case Javascript.String:
		q.Javascript(`window.localStorage.setItem("`+name.raw()+`", `+string(t.Expression)+`);`)

	case Javascript.Integer:
		q.Javascript(`window.localStorage.setItem("`+name.raw()+`", (`+string(t.Expression)+`).toString());`)

	default:
		panic("Unimplemented")
	}

}

type Int struct {
	Variable
}

func NewInt() Int {
	return Int{NewVariable()}
}

func (i Int) Script(q Script) qlova.Int {
	return q.IntFromLanguageType(Javascript.Integer{
		Expression: language.Statement(`parseInt(window.localStorage.getItem("`+string(i.Variable)+`") || "0")`),
	})
}