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
	type EgReq struct{ relay.Req[int] }
	req := EgReq{}
	expectedRes := 10
	handler := func(req EgReq) (int, error) {
		return expectedRes, nil
	}

	var r, err = relay.TryNewRelay(relay.NewConfigBuilder().Build())
	if err != nil {
		t.Errorf("unexpected error %s\n", err.Error())
	}

	relay.Register(r, handler)

	res, err := relay.Handle(r, req)
	if res != expectedRes {
		t.Errorf("unexpected response.\ngot %d\nexpected %d\n", res, expectedRes)
	}

	if err != nil {
		t.Errorf("unexpected response error %s\n", err.Error())
	}

}
