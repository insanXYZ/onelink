package message

import (
	"errors"
)

var (
	ERR_CREDENTIALS             = createError("credentials is incorrect")
	ERR_LOGIN                   = createError("email or password is incorrect")
	ERR_REGISTER                = ERR_CREDENTIALS
	ERR_REGISTER_EMAIL_USED     = createError("email is already in use")
	ERR_CREATE_SITE_DOMAIN_USED = createError("domain is already in use")
	ERR_CREATE_SITE             = createError("create site failed")
	ERR_PUBLISH_SITE_DOMAIN     = createError("domain site is incorrect")
	ERR_CLICK                   = ERR_CREDENTIALS
)

func createError(message string) error {
	return errors.New(message)
}
