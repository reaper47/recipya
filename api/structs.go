package api

// ErrorJson is an object that stores an
// Error as the value of the "error" field.
//
// For example: {
//	"error": ...
//}
type ErrorJson struct {
	Objects Error `json:"error"`
}

// Error is an object that stores information
// related the error at hand.
//
// For example: {
//	"code": 404,
//	"message": "oops",
//  "status": "Not Found"
//}
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Status  string `json:"status"`
}
