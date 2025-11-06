package types

import "net/http"

type Password string

type RequestAugment func(*http.Request) error
