// Package response defines the standard API response envelope used by all
// HTTP handlers. The canonical format is:
//
//	{ "status": bool, "message": string, "data": any }
//
// Use Success() and Error() helpers to construct consistent responses.
package response
