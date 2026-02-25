package enum

type Code int

const (
	CodeSuccess Code = 0
	CodeFailed       = 500

	CodeInvalidParameter        = 400
	CodeTokenInvalid            = 401
	CodeInsufficientPermissions = 403
)
