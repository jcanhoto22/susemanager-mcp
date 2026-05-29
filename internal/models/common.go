package models

type ApiVersion struct {
	Version string `xmlrpc:"version"`
}

type Error struct {
	Code    int    `xmlrpc:"faultCode"`
	Message string `xmlrpc:"faultString"`
}
