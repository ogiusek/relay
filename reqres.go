package relay

// e.g.
//
//	type EgRes struct {
//		ExampleField int
//	}
type Res interface{}

// every interface should use this
// e.g.
//
//	type EgReq struct {
//		Req[EgRes] // it can be nil its just a marker
//		Field1     int
//		Field2     int
//	}
type Req[Response Res] interface {
	// this is never called
	// this is just request marker
	Response() Response
}

type Handler[Request Req[Response], Response any] func(Request) (Response, error)

type DefaultHandler func(req any) (Res, error)
