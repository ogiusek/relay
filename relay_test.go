package relay_test

import (
	"errors"
	"testing"

	"github.com/ogiusek/relay/v2"
)

func afterPanic() {
	print("\033[1A") // go 1 line up
	print("\033[2K") // clear line
}

func TestDefaultHandler(t *testing.T) {
	type EgReq struct{ relay.Req[int] }

	{
		r := relay.NewBuilder().Build()

		var req = EgReq{}
		_, err := relay.Handle(r, req)
		if err != relay.ErrHandlerNotFound {
			t.Errorf("unexpected error.\ngot %s\nexpected %s\n", err.Error(), relay.ErrHandlerNotFound.Error())
		}
	}

	{
		var customErr = errors.New("")
		r := relay.NewBuilder().
			DefaultHandler(func(req any) (relay.Res, error) { return nil, customErr }).
			Build()

		req := EgReq{}
		_, err := relay.Handle(r, req)
		if err != customErr {
			t.Errorf("unexpected error.\ngot %s\nexpected %s\n", err.Error(), customErr.Error())
		}
	}
}

func TestRegisteringTwice(t *testing.T) {
	type EgReq struct{ relay.Req[int] }
	handler := func(req EgReq) (int, error) {
		return 0, nil
	}

	t.Run("panics", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				afterPanic()
			} else {
				t.Errorf("relay do not panics when handler is registered twice")
			}
		}()
		relay.NewBuilder().
			Wrap(func(b relay.Builder) relay.Builder { return relay.Register(b, handler) }).
			Wrap(func(b relay.Builder) relay.Builder { return relay.Register(b, handler) }).
			Build()
	})
}

func TestHandler(t *testing.T) {
	type EgReq struct {
		relay.Req[int]
		EgArr []int // arrays are not composable
	}
	req := EgReq{}
	expectedRes := 10
	handler := func(req EgReq) (int, error) { return expectedRes, nil }

	type EgReq2 struct{ relay.Req[int] }
	req2 := EgReq2{}
	expectedRes2 := 11
	handler2 := func(req EgReq2) (int, error) { return expectedRes2, nil }

	r := relay.NewBuilder().
		Wrap(func(b relay.Builder) relay.Builder { return relay.Register(b, handler) }).
		Wrap(func(b relay.Builder) relay.Builder { return relay.Register(b, handler2) }).
		Build()

	res, err := relay.Handle(r, req)
	if res != expectedRes {
		t.Errorf("unexpected response.\ngot %d\nexpected %d\n", res, expectedRes)
	}

	if err != nil {
		t.Errorf("unexpected response error %s\n", err.Error())
	}

	res2, err := relay.Handle(r, req2)
	if res2 != expectedRes2 {
		t.Errorf("unexpected response.\ngot %d\nexpected %d\n", res2, expectedRes2)
	}

	if err != nil {
		t.Errorf("unexpected response error %s\n", err.Error())
	}

}
