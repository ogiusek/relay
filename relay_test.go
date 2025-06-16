package relay_test

import (
	"errors"
	"testing"

	"github.com/ogiusek/relay"
)

func TestDefaultHandler(t *testing.T) {
	type EgReq struct{ relay.Req[int] }

	{
		var r, err = relay.TryNewRelay(relay.NewConfigBuilder().Build())
		if err != nil {
			t.Errorf("%s\n", err.Error())
		}

		var req = EgReq{}
		_, err = relay.Handle(r, req)
		if err != relay.ErrHandlerNotFound {
			t.Errorf("unexpected error.\ngot %s\nexpected %s\n", err.Error(), relay.ErrHandlerNotFound.Error())
		}
	}

	{
		var customErr = errors.New("")
		r, err := relay.TryNewRelay(relay.NewConfigBuilder().
			DefaultHandler(func(req any) (relay.Res, error) { return nil, customErr }).
			Build(),
		)
		if err != nil {
			t.Errorf("%s\n", err.Error())
		}

		req := EgReq{}
		_, err = relay.Handle(r, req)
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

	r, err := relay.TryNewRelay(relay.NewConfigBuilder().Build())
	if err != nil {
		t.Errorf("%s\n", err.Error())
	}
	if err := relay.TryRegister(r, handler); err != nil {
		t.Errorf("unexpected error %s\n", err.Error())
	}
	if err := relay.TryRegister(r, handler); err != relay.ErrHandlerAlreadyExists {
		t.Errorf("unexpected error %s\nexpected %s\n", err.Error(), relay.ErrHandlerAlreadyExists)
	}
}

func TestHandler(t *testing.T) {
	r := relay.NewRelay(relay.NewConfigBuilder().Build())
	type EgReq struct {
		relay.Req[int]
		EgArr []int // arrays are not composable
	}
	req := EgReq{}
	expectedRes := 10
	handler := func(req EgReq) (int, error) { return expectedRes, nil }
	relay.Register(r, handler)

	type EgReq2 struct{ relay.Req[int] }
	req2 := EgReq2{}
	expectedRes2 := 11
	handler2 := func(req EgReq2) (int, error) { return expectedRes2, nil }
	relay.Register(r, handler2)

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
