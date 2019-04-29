package element

type ElementMADL struct {
	ElemId    string
	ElemType  string
}

type ElementGo struct {
	ElemId string
	ElemType interface{}
	CSP string
}