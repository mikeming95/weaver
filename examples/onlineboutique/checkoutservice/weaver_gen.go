// go:build !ignoreWeaverGen

package checkoutservice

// Code generated by "weaver generate". DO NOT EDIT.
import (
	"context"
	"errors"
	"fmt"
	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/examples/onlineboutique/types"
	"github.com/ServiceWeaver/weaver/runtime/codegen"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"reflect"
	"time"
)

func init() {
	codegen.Register(codegen.Registration{
		Name:        "github.com/ServiceWeaver/weaver/examples/onlineboutique/checkoutservice/T",
		Iface:       reflect.TypeOf((*T)(nil)).Elem(),
		New:         func() any { return &impl{} },
		LocalStubFn: func(impl any, tracer trace.Tracer) any { return impl_local_stub{impl: impl.(T), tracer: tracer} },
		ClientStubFn: func(stub codegen.Stub, caller string) any {
			return impl_client_stub{stub: stub, placeOrderMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "github.com/ServiceWeaver/weaver/examples/onlineboutique/checkoutservice/T", Method: "PlaceOrder"})}
		},
		ServerStubFn: func(impl any, addLoad func(uint64, float64)) codegen.Server {
			return impl_server_stub{impl: impl.(T), addLoad: addLoad}
		},
	})
}

// Local stub implementations.

type impl_local_stub struct {
	impl   T
	tracer trace.Tracer
}

func (s impl_local_stub) PlaceOrder(ctx context.Context, a0 PlaceOrderRequest) (r0 types.Order, err error) {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.tracer.Start(ctx, "github.com/ServiceWeaver/weaver/examples/onlineboutique/checkoutservice/T.PlaceOrder", trace.WithSpanKind(trace.SpanKindInternal))
		defer func() {
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			span.End()
		}()
	}

	return s.impl.PlaceOrder(ctx, a0)
}

// Client stub implementations.

type impl_client_stub struct {
	stub              codegen.Stub
	placeOrderMetrics *codegen.MethodMetrics
}

func (s impl_client_stub) PlaceOrder(ctx context.Context, a0 PlaceOrderRequest) (r0 types.Order, err error) {
	// Update metrics.
	start := time.Now()
	s.placeOrderMetrics.Count.Add(1)

	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.stub.Tracer().Start(ctx, "github.com/ServiceWeaver/weaver/examples/onlineboutique/checkoutservice/T.PlaceOrder", trace.WithSpanKind(trace.SpanKindClient))
	}

	defer func() {
		// Catch and return any panics detected during encoding/decoding/rpc.
		if err == nil {
			err = codegen.CatchPanics(recover())
			if err != nil {
				err = errors.Join(weaver.RemoteCallError, err)
			}
		}

		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			s.placeOrderMetrics.ErrorCount.Add(1)
		}
		span.End()

		s.placeOrderMetrics.Latency.Put(float64(time.Since(start).Microseconds()))
	}()

	// Encode arguments.
	enc := codegen.NewEncoder()
	(a0).WeaverMarshal(enc)
	var shardKey uint64

	// Call the remote method.
	s.placeOrderMetrics.BytesRequest.Put(float64(len(enc.Data())))
	var results []byte
	results, err = s.stub.Run(ctx, 0, enc.Data(), shardKey)
	if err != nil {
		err = errors.Join(weaver.RemoteCallError, err)
		return
	}
	s.placeOrderMetrics.BytesReply.Put(float64(len(results)))

	// Decode the results.
	dec := codegen.NewDecoder(results)
	(&r0).WeaverUnmarshal(dec)
	err = dec.Error()
	return
}

// Server stub implementations.

type impl_server_stub struct {
	impl    T
	addLoad func(key uint64, load float64)
}

// GetStubFn implements the stub.Server interface.
func (s impl_server_stub) GetStubFn(method string) func(ctx context.Context, args []byte) ([]byte, error) {
	switch method {
	case "PlaceOrder":
		return s.placeOrder
	default:
		return nil
	}
}

func (s impl_server_stub) placeOrder(ctx context.Context, args []byte) (res []byte, err error) {
	// Catch and return any panics detected during encoding/decoding/rpc.
	defer func() {
		if err == nil {
			err = codegen.CatchPanics(recover())
		}
	}()

	// Decode arguments.
	dec := codegen.NewDecoder(args)
	var a0 PlaceOrderRequest
	(&a0).WeaverUnmarshal(dec)

	// TODO(rgrandl): The deferred function above will recover from panics in the
	// user code: fix this.
	// Call the local method.
	r0, appErr := s.impl.PlaceOrder(ctx, a0)

	// Encode the results.
	enc := codegen.NewEncoder()
	(r0).WeaverMarshal(enc)
	enc.Error(appErr)
	return enc.Data(), nil
}

// AutoMarshal implementations.

var _ codegen.AutoMarshal = &PlaceOrderRequest{}

func (x *PlaceOrderRequest) WeaverMarshal(enc *codegen.Encoder) {
	if x == nil {
		panic(fmt.Errorf("PlaceOrderRequest.WeaverMarshal: nil receiver"))
	}
	enc.String(x.UserID)
	enc.String(x.UserCurrency)
	(x.Address).WeaverMarshal(enc)
	enc.String(x.Email)
	(x.CreditCard).WeaverMarshal(enc)
}

func (x *PlaceOrderRequest) WeaverUnmarshal(dec *codegen.Decoder) {
	if x == nil {
		panic(fmt.Errorf("PlaceOrderRequest.WeaverUnmarshal: nil receiver"))
	}
	x.UserID = dec.String()
	x.UserCurrency = dec.String()
	(&x.Address).WeaverUnmarshal(dec)
	x.Email = dec.String()
	(&x.CreditCard).WeaverUnmarshal(dec)
}
