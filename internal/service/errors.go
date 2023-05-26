package service

import "errors"

var (
	EmptyBalanceError = errors.New("empty balance")
	NotAuthError      = errors.New("user not authorized")
	decodingError     = errors.New("error url decoding")
	blockSizeError    = errors.New("ciphertext  is too short")
	AdminUsageErr     = errors.New("invalid command")
	PromoError        = errors.New("invalid Promo")
	PromoUsedError    = errors.New("promo already used")
)
