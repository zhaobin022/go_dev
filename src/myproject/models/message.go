package models

import (
	_ "github.com/go-sql-driver/mysql"
)

type ApiRequest struct {
	RequestType string
	Data        []*int
}
