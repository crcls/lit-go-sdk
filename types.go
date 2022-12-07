package lit

import "net/http"

type SendReqFuncType func(url string, body []byte) (*http.Response, error)
