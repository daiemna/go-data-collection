// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: dataframebatch.proto

package services

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// TimeSeriesDataClient is the client API for TimeSeriesData service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TimeSeriesDataClient interface {
	StreamRecords(ctx context.Context, opts ...grpc.CallOption) (TimeSeriesData_StreamRecordsClient, error)
	StreamDataSeries(ctx context.Context, opts ...grpc.CallOption) (TimeSeriesData_StreamDataSeriesClient, error)
	StreamDataframes(ctx context.Context, opts ...grpc.CallOption) (TimeSeriesData_StreamDataframesClient, error)
}

type timeSeriesDataClient struct {
	cc grpc.ClientConnInterface
}

func NewTimeSeriesDataClient(cc grpc.ClientConnInterface) TimeSeriesDataClient {
	return &timeSeriesDataClient{cc}
}

func (c *timeSeriesDataClient) StreamRecords(ctx context.Context, opts ...grpc.CallOption) (TimeSeriesData_StreamRecordsClient, error) {
	stream, err := c.cc.NewStream(ctx, &TimeSeriesData_ServiceDesc.Streams[0], "/TimeSeriesData/StreamRecords", opts...)
	if err != nil {
		return nil, err
	}
	x := &timeSeriesDataStreamRecordsClient{stream}
	return x, nil
}

type TimeSeriesData_StreamRecordsClient interface {
	Send(*DataRecord) error
	Recv() (*Response, error)
	grpc.ClientStream
}

type timeSeriesDataStreamRecordsClient struct {
	grpc.ClientStream
}

func (x *timeSeriesDataStreamRecordsClient) Send(m *DataRecord) error {
	return x.ClientStream.SendMsg(m)
}

func (x *timeSeriesDataStreamRecordsClient) Recv() (*Response, error) {
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *timeSeriesDataClient) StreamDataSeries(ctx context.Context, opts ...grpc.CallOption) (TimeSeriesData_StreamDataSeriesClient, error) {
	stream, err := c.cc.NewStream(ctx, &TimeSeriesData_ServiceDesc.Streams[1], "/TimeSeriesData/StreamDataSeries", opts...)
	if err != nil {
		return nil, err
	}
	x := &timeSeriesDataStreamDataSeriesClient{stream}
	return x, nil
}

type TimeSeriesData_StreamDataSeriesClient interface {
	Send(*DataSeries) error
	Recv() (*Response, error)
	grpc.ClientStream
}

type timeSeriesDataStreamDataSeriesClient struct {
	grpc.ClientStream
}

func (x *timeSeriesDataStreamDataSeriesClient) Send(m *DataSeries) error {
	return x.ClientStream.SendMsg(m)
}

func (x *timeSeriesDataStreamDataSeriesClient) Recv() (*Response, error) {
	m := new(Response)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *timeSeriesDataClient) StreamDataframes(ctx context.Context, opts ...grpc.CallOption) (TimeSeriesData_StreamDataframesClient, error) {
	stream, err := c.cc.NewStream(ctx, &TimeSeriesData_ServiceDesc.Streams[2], "/TimeSeriesData/StreamDataframes", opts...)
	if err != nil {
		return nil, err
	}
	x := &timeSeriesDataStreamDataframesClient{stream}
	return x, nil
}

type TimeSeriesData_StreamDataframesClient interface {
	Send(*Dataframe) error
	Recv() (*DataframeResponse, error)
	grpc.ClientStream
}

type timeSeriesDataStreamDataframesClient struct {
	grpc.ClientStream
}

func (x *timeSeriesDataStreamDataframesClient) Send(m *Dataframe) error {
	return x.ClientStream.SendMsg(m)
}

func (x *timeSeriesDataStreamDataframesClient) Recv() (*DataframeResponse, error) {
	m := new(DataframeResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TimeSeriesDataServer is the server API for TimeSeriesData service.
// All implementations must embed UnimplementedTimeSeriesDataServer
// for forward compatibility
type TimeSeriesDataServer interface {
	StreamRecords(TimeSeriesData_StreamRecordsServer) error
	StreamDataSeries(TimeSeriesData_StreamDataSeriesServer) error
	StreamDataframes(TimeSeriesData_StreamDataframesServer) error
	mustEmbedUnimplementedTimeSeriesDataServer()
}

// UnimplementedTimeSeriesDataServer must be embedded to have forward compatible implementations.
type UnimplementedTimeSeriesDataServer struct {
}

func (UnimplementedTimeSeriesDataServer) StreamRecords(TimeSeriesData_StreamRecordsServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamRecords not implemented")
}
func (UnimplementedTimeSeriesDataServer) StreamDataSeries(TimeSeriesData_StreamDataSeriesServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamDataSeries not implemented")
}
func (UnimplementedTimeSeriesDataServer) StreamDataframes(TimeSeriesData_StreamDataframesServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamDataframes not implemented")
}
func (UnimplementedTimeSeriesDataServer) mustEmbedUnimplementedTimeSeriesDataServer() {}

// UnsafeTimeSeriesDataServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TimeSeriesDataServer will
// result in compilation errors.
type UnsafeTimeSeriesDataServer interface {
	mustEmbedUnimplementedTimeSeriesDataServer()
}

func RegisterTimeSeriesDataServer(s grpc.ServiceRegistrar, srv TimeSeriesDataServer) {
	s.RegisterService(&TimeSeriesData_ServiceDesc, srv)
}

func _TimeSeriesData_StreamRecords_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(TimeSeriesDataServer).StreamRecords(&timeSeriesDataStreamRecordsServer{stream})
}

type TimeSeriesData_StreamRecordsServer interface {
	Send(*Response) error
	Recv() (*DataRecord, error)
	grpc.ServerStream
}

type timeSeriesDataStreamRecordsServer struct {
	grpc.ServerStream
}

func (x *timeSeriesDataStreamRecordsServer) Send(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *timeSeriesDataStreamRecordsServer) Recv() (*DataRecord, error) {
	m := new(DataRecord)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _TimeSeriesData_StreamDataSeries_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(TimeSeriesDataServer).StreamDataSeries(&timeSeriesDataStreamDataSeriesServer{stream})
}

type TimeSeriesData_StreamDataSeriesServer interface {
	Send(*Response) error
	Recv() (*DataSeries, error)
	grpc.ServerStream
}

type timeSeriesDataStreamDataSeriesServer struct {
	grpc.ServerStream
}

func (x *timeSeriesDataStreamDataSeriesServer) Send(m *Response) error {
	return x.ServerStream.SendMsg(m)
}

func (x *timeSeriesDataStreamDataSeriesServer) Recv() (*DataSeries, error) {
	m := new(DataSeries)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func _TimeSeriesData_StreamDataframes_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(TimeSeriesDataServer).StreamDataframes(&timeSeriesDataStreamDataframesServer{stream})
}

type TimeSeriesData_StreamDataframesServer interface {
	Send(*DataframeResponse) error
	Recv() (*Dataframe, error)
	grpc.ServerStream
}

type timeSeriesDataStreamDataframesServer struct {
	grpc.ServerStream
}

func (x *timeSeriesDataStreamDataframesServer) Send(m *DataframeResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *timeSeriesDataStreamDataframesServer) Recv() (*Dataframe, error) {
	m := new(Dataframe)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// TimeSeriesData_ServiceDesc is the grpc.ServiceDesc for TimeSeriesData service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TimeSeriesData_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "TimeSeriesData",
	HandlerType: (*TimeSeriesDataServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamRecords",
			Handler:       _TimeSeriesData_StreamRecords_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "StreamDataSeries",
			Handler:       _TimeSeriesData_StreamDataSeries_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
		{
			StreamName:    "StreamDataframes",
			Handler:       _TimeSeriesData_StreamDataframes_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "dataframebatch.proto",
}
