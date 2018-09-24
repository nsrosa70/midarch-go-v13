package message

type RequestHeader struct{
	Magic string
}

type RequestBody struct {
	Op string
	Args interface{}
}

type ReplyHeader struct {
	Status int
}

type ReplyBody struct {
	Reply interface{}
}

type MIOP struct {
	RequestHeader RequestHeader
	RequestBody RequestBody
	ReplyHeader ReplyHeader
	ReplyBody ReplyBody
}
