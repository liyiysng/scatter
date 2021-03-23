package handle

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/liyiysng/scatter/encoding"
	_ "github.com/liyiysng/scatter/encoding/json"
)

type TestReq struct {
	LValue int `json:"l_value"`
	RValue int `json:"r_value"`
}

type TestRes struct {
	Sum int `json:"sum"`
}

type fooSessionType int

func hookCall(ctx context.Context, session interface{}, srv interface{}, srvName string, methodName string, req interface{}, callee func(req interface{}) (res interface{}, err error)) error {

	beg := time.Now()

	res, err := callee(req)

	myLog.Infof("%s.%s(req:%v) (res:%v,err:%v) => %v", srvName, methodName, req, res, err, time.Now().Sub(beg))

	return err
}

func notifyCall(ctx context.Context, session interface{}, srv interface{}, srvName string, methodName string, req interface{}, callee func(req interface{}) (err error)) error {

	beg := time.Now()

	err := callee(req)

	myLog.Infof("%s.%s(req:%v) (err:%v) => %v", srvName, methodName, req, err, time.Now().Sub(beg))

	return err
}

type TestService struct {
	t *testing.T
}

func (t *TestService) Add(ctx context.Context, session fooSessionType, req *TestReq) (res *TestRes, err error) {
	time.Sleep(time.Millisecond)
	return &TestRes{
		Sum: req.LValue + req.RValue,
	}, nil
}

func (t *TestService) ErrTest(ctx context.Context, session fooSessionType, req *TestReq) (res *TestRes, err error) {
	return nil, errors.New("ErrTest")
}

func (t *TestService) Wron(ctx context.Context, session fooSessionType) (res *TestRes, err error) {
	return nil, errors.New("ErrTest")
}

func (t *TestService) WrongReqType(ctx context.Context, session fooSessionType, req TestReq) (res *TestRes, err error) {
	return nil, errors.New("ErrTest")
}

func (t *TestService) WrongResType(ctx context.Context, session fooSessionType, req *TestReq) (res TestRes, err error) {
	return TestRes{}, errors.New("ErrTest")
}

func (t *TestService) Notify(ctx context.Context, session fooSessionType, req *TestReq) (err error) {
	return nil
}

func (t *TestService) AddOptArgs(ctx context.Context, session fooSessionType, req *TestReq, subValue int, addtionValue int32) (res *TestRes, err error) {
	time.Sleep(time.Millisecond)
	return &TestRes{
		Sum: req.LValue + req.RValue - subValue + int(addtionValue),
	}, nil
}

func (t *TestService) AddOptArg(ctx context.Context, session fooSessionType, req *TestReq, subValue int, addtionValue int32) (res *TestRes, err error) {
	time.Sleep(time.Millisecond)
	return &TestRes{
		Sum: req.LValue + req.RValue - subValue + int(addtionValue),
	}, nil
}

func TestServiceHandler(t *testing.T) {

	otp := &Option{
		Codec:            encoding.GetCodec("json"),
		ReqTypeValidator: func(reqType reflect.Type) error { return nil },
		ResTypeValidator: func(reqType reflect.Type) error { return nil },
		SessionType:      reflect.TypeOf((*fooSessionType)(nil)).Elem(),
		HookCall:         hookCall,
		HookNofify:       notifyCall,
		OptArgs: &OptionalArgs{
			ArgsTypeValidator: func(srvName string, methodName string, argsType []reflect.Type) error {
				if strings.HasSuffix(methodName, "Args") {
					if len(argsType) != 2 {
						return errors.New("invalid optional args num")
					}
					if argsType[0] != reflect.TypeOf((*int)(nil)).Elem() {
						return errors.New("first optional args type must be int")
					}
					if argsType[1] != reflect.TypeOf((*int32)(nil)).Elem() {
						return errors.New("second optional args type must be int32")
					}
				}
				return nil
			},
			Call: func(session interface{}, srvName string, methodName string, callee func(argValues ...interface{}) error) error {

				if strings.HasSuffix(methodName, "Args") {
					return callee(10, int32(1))
				}
				return errors.New("optional args not support")
			},
		},
	}

	h := NewServiceHandle(otp)

	h.Register(&TestService{
		t: t,
	})

	req := &TestReq{
		LValue: 5,
		RValue: 10,
	}

	b, _ := json.Marshal(req)

	var session fooSessionType = 555

	resB, err := h.Call(context.Background(), session, "TestService", "Add", b)

	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(resB))

	resB, err = h.Call(context.Background(), session, "TestService", "ErrTest", b)

	if err != nil {
		t.Log(err)
	}

	t.Log(string(resB))

	err = h.Notify(context.Background(), session, "TestService", "Notify", b)

	if err != nil {
		t.Log(err)
	}

	resB, err = h.Call(context.Background(), session, "TestService", "AddOptArgs", b)

	if err != nil {
		t.Log(err)
	}

	t.Log(string(resB))

	resB, err = h.Call(context.Background(), session, "TestService", "AddOptArg", b)

	if err != nil {
		t.Log(err)
	}

	t.Log(string(resB))

}
