package script

import qlova "github.com/qlova/script"

//Display a message to the user.
func (q Script) Alert(message qlova.String) {
	q.js.Run("alert", message)
}

//Display a confirmation box to the user, returns a bool indicating true for 'ok' false for 'cancel'.
func (q Script) Confirm(message qlova.String) qlova.Bool {
	return q.js.Call("confirm", message).Bool()
}

//Display a prompt that requests a string from the user.
func (q Script) Prompt(message qlova.String) qlova.String {
	return q.js.Call("prompt", message).String()
}