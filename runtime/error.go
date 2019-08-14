package runtime

import "github.com/golang/protobuf/ptypes/any"

type errorBody struct {
	Error string `protobuf:"bytes,1,name=error" json:"error"`
	// This is to make the error more compatible with users that expect errors to be Status objects:
	// https://github.com/grpc/grpc/blob/master/src/proto/grpc/status/status.proto
	// It should be the exact same message as the Error field.
	Message string     `protobuf:"bytes,1,name=message" json:"message"`
	Code    int32      `protobuf:"varint,2,name=code" json:"code"`
	Details []*any.Any `protobuf:"bytes,3,rep,name=details" json:"details,omitempty"`
}
