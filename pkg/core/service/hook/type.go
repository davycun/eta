package hook

const (
	CallbackBefore CallbackPosition = 1
	CallbackAfter  CallbackPosition = 2
)

type CallbackPosition int

type AuthFunc Callback
