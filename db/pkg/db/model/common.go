package model

import "google.golang.org/protobuf/encoding/protojson"

var protoJsonUnmarshalOptions = protojson.UnmarshalOptions{
	DiscardUnknown: true,
}
