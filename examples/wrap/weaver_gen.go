// go:build !ignoreWeaverGen

package main

// Code generated by "weaver generate". DO NOT EDIT.
import (
	"context"
	"errors"
	"github.com/ServiceWeaver/weaver"
	"github.com/ServiceWeaver/weaver/runtime/codegen"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"reflect"
	"time"
)

func init() {
	codegen.Register(codegen.Registration{
		Name:  "github.com/ServiceWeaver/weaver/examples/wrap/Wrapper",
		Iface: reflect.TypeOf((*Wrapper)(nil)).Elem(),
		New:   func() any { return &wrapper{} },
		LocalStubFn: func(impl any, tracer trace.Tracer) any {
			return wrapper_local_stub{impl: impl.(Wrapper), tracer: tracer}
		},
		ClientStubFn: func(stub codegen.Stub, caller string) any {
			return wrapper_client_stub{stub: stub, wrapMetrics: codegen.MethodMetricsFor(codegen.MethodLabels{Caller: caller, Component: "github.com/ServiceWeaver/weaver/examples/wrap/Wrapper", Method: "Wrap"})}
		},
		ServerStubFn: func(impl any, addLoad func(uint64, float64)) codegen.Server {
			return wrapper_server_stub{impl: impl.(Wrapper), addLoad: addLoad}
		},
	})
}

// Local stub implementations.

type wrapper_local_stub struct {
	impl   Wrapper
	tracer trace.Tracer
}

func (s wrapper_local_stub) Wrap(ctx context.Context, a0 string, a1 int) (r0 string, err error) {
	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.tracer.Start(ctx, "github.com/ServiceWeaver/weaver/examples/wrap/Wrapper.Wrap", trace.WithSpanKind(trace.SpanKindInternal))
		defer func() {
			if err != nil {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
			}
			span.End()
		}()
	}

	return s.impl.Wrap(ctx, a0, a1)
}

// Client stub implementations.

type wrapper_client_stub struct {
	stub        codegen.Stub
	wrapMetrics *codegen.MethodMetrics
}

func (s wrapper_client_stub) Wrap(ctx context.Context, a0 string, a1 int) (r0 string, err error) {
	// Update metrics.
	start := time.Now()
	s.wrapMetrics.Count.Add(1)

	span := trace.SpanFromContext(ctx)
	if span.SpanContext().IsValid() {
		// Create a child span for this method.
		ctx, span = s.stub.Tracer().Start(ctx, "github.com/ServiceWeaver/weaver/examples/wrap/Wrapper.Wrap", trace.WithSpanKind(trace.SpanKindClient))
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
			s.wrapMetrics.ErrorCount.Add(1)
		}
		span.End()

		s.wrapMetrics.Latency.Put(float64(time.Since(start).Microseconds()))
	}()

	// Preallocate a buffer of the right size.
	size := 0
	size += (4 + len(a0))
	size += 8
	enc := codegen.NewEncoder()
	enc.Reset(size)

	// Encode arguments.
	enc.String(a0)
	enc.Int(a1)
	var shardKey uint64

	// Call the remote method.
	s.wrapMetrics.BytesRequest.Put(float64(len(enc.Data())))
	var results []byte
	results, err = s.stub.Run(ctx, 0, enc.Data(), shardKey)
	if err != nil {
		err = errors.Join(weaver.RemoteCallError, err)
		return
	}
	s.wrapMetrics.BytesReply.Put(float64(len(results)))

	// Decode the results.
	dec := codegen.NewDecoder(results)
	r0 = dec.String()
	err = dec.Error()
	return
}

// Server stub implementations.

type wrapper_server_stub struct {
	impl    Wrapper
	addLoad func(key uint64, load float64)
}

// GetStubFn implements the stub.Server interface.
func (s wrapper_server_stub) GetStubFn(method string) func(ctx context.Context, args []byte) ([]byte, error) {
	switch method {
	case "Wrap":
		return s.wrap
	default:
		return nil
	}
}

func (s wrapper_server_stub) wrap(ctx context.Context, args []byte) (res []byte, err error) {
	// Catch and return any panics detected during encoding/decoding/rpc.
	defer func() {
		if err == nil {
			err = codegen.CatchPanics(recover())
		}
	}()

	// Decode arguments.
	dec := codegen.NewDecoder(args)
	var a0 string
	a0 = dec.String()
	var a1 int
	a1 = dec.Int()

	// TODO(rgrandl): The deferred function above will recover from panics in the
	// user code: fix this.
	// Call the local method.
	r0, appErr := s.impl.Wrap(ctx, a0, a1)

	// Encode the results.
	enc := codegen.NewEncoder()
	enc.String(r0)
	enc.Error(appErr)
	return enc.Data(), nil
}
