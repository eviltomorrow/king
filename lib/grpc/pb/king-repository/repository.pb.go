// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.30.0-devel
// 	protoc        v4.23.2
// source: repository.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	wrapperspb "google.golang.org/protobuf/types/known/wrapperspb"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type QuoteRequest_Mode int32

const (
	QuoteRequest_Day  QuoteRequest_Mode = 0
	QuoteRequest_Week QuoteRequest_Mode = 1
)

// Enum value maps for QuoteRequest_Mode.
var (
	QuoteRequest_Mode_name = map[int32]string{
		0: "Day",
		1: "Week",
	}
	QuoteRequest_Mode_value = map[string]int32{
		"Day":  0,
		"Week": 1,
	}
)

func (x QuoteRequest_Mode) Enum() *QuoteRequest_Mode {
	p := new(QuoteRequest_Mode)
	*p = x
	return p
}

func (x QuoteRequest_Mode) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (QuoteRequest_Mode) Descriptor() protoreflect.EnumDescriptor {
	return file_repository_proto_enumTypes[0].Descriptor()
}

func (QuoteRequest_Mode) Type() protoreflect.EnumType {
	return &file_repository_proto_enumTypes[0]
}

func (x QuoteRequest_Mode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use QuoteRequest_Mode.Descriptor instead.
func (QuoteRequest_Mode) EnumDescriptor() ([]byte, []int) {
	return file_repository_proto_rawDescGZIP(), []int{2, 0}
}

type Counter struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	AffectedStock     int64 `protobuf:"varint,1,opt,name=affected_stock,json=affectedStock,proto3" json:"affected_stock,omitempty"`
	AffectedQuoteDay  int64 `protobuf:"varint,2,opt,name=affected_quote_day,json=affectedQuoteDay,proto3" json:"affected_quote_day,omitempty"`
	AffectedQuoteWeek int64 `protobuf:"varint,3,opt,name=affected_quote_week,json=affectedQuoteWeek,proto3" json:"affected_quote_week,omitempty"`
}

func (x *Counter) Reset() {
	*x = Counter{}
	if protoimpl.UnsafeEnabled {
		mi := &file_repository_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Counter) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Counter) ProtoMessage() {}

func (x *Counter) ProtoReflect() protoreflect.Message {
	mi := &file_repository_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Counter.ProtoReflect.Descriptor instead.
func (*Counter) Descriptor() ([]byte, []int) {
	return file_repository_proto_rawDescGZIP(), []int{0}
}

func (x *Counter) GetAffectedStock() int64 {
	if x != nil {
		return x.AffectedStock
	}
	return 0
}

func (x *Counter) GetAffectedQuoteDay() int64 {
	if x != nil {
		return x.AffectedQuoteDay
	}
	return 0
}

func (x *Counter) GetAffectedQuoteWeek() int64 {
	if x != nil {
		return x.AffectedQuoteWeek
	}
	return 0
}

type Stock struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code    string `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	Name    string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	Suspend string `protobuf:"bytes,3,opt,name=suspend,proto3" json:"suspend,omitempty"`
}

func (x *Stock) Reset() {
	*x = Stock{}
	if protoimpl.UnsafeEnabled {
		mi := &file_repository_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Stock) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Stock) ProtoMessage() {}

func (x *Stock) ProtoReflect() protoreflect.Message {
	mi := &file_repository_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Stock.ProtoReflect.Descriptor instead.
func (*Stock) Descriptor() ([]byte, []int) {
	return file_repository_proto_rawDescGZIP(), []int{1}
}

func (x *Stock) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *Stock) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Stock) GetSuspend() string {
	if x != nil {
		return x.Suspend
	}
	return ""
}

type QuoteRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code  string            `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	Date  string            `protobuf:"bytes,2,opt,name=date,proto3" json:"date,omitempty"`
	Limit int64             `protobuf:"varint,3,opt,name=limit,proto3" json:"limit,omitempty"`
	Mode  QuoteRequest_Mode `protobuf:"varint,4,opt,name=mode,proto3,enum=repository.QuoteRequest_Mode" json:"mode,omitempty"`
}

func (x *QuoteRequest) Reset() {
	*x = QuoteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_repository_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QuoteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QuoteRequest) ProtoMessage() {}

func (x *QuoteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_repository_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use QuoteRequest.ProtoReflect.Descriptor instead.
func (*QuoteRequest) Descriptor() ([]byte, []int) {
	return file_repository_proto_rawDescGZIP(), []int{2}
}

func (x *QuoteRequest) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *QuoteRequest) GetDate() string {
	if x != nil {
		return x.Date
	}
	return ""
}

func (x *QuoteRequest) GetLimit() int64 {
	if x != nil {
		return x.Limit
	}
	return 0
}

func (x *QuoteRequest) GetMode() QuoteRequest_Mode {
	if x != nil {
		return x.Mode
	}
	return QuoteRequest_Day
}

type Quote struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Code            string  `protobuf:"bytes,1,opt,name=code,proto3" json:"code,omitempty"`
	Open            float64 `protobuf:"fixed64,2,opt,name=open,proto3" json:"open,omitempty"`
	Close           float64 `protobuf:"fixed64,3,opt,name=close,proto3" json:"close,omitempty"`
	High            float64 `protobuf:"fixed64,4,opt,name=high,proto3" json:"high,omitempty"`
	Low             float64 `protobuf:"fixed64,5,opt,name=low,proto3" json:"low,omitempty"`
	YesterdayClosed float64 `protobuf:"fixed64,6,opt,name=yesterday_closed,json=yesterdayClosed,proto3" json:"yesterday_closed,omitempty"`
	Volume          uint64  `protobuf:"varint,7,opt,name=volume,proto3" json:"volume,omitempty"`
	Account         float64 `protobuf:"fixed64,8,opt,name=account,proto3" json:"account,omitempty"`
	Date            string  `protobuf:"bytes,9,opt,name=date,proto3" json:"date,omitempty"`
	NumOfYear       int32   `protobuf:"varint,10,opt,name=num_of_year,json=numOfYear,proto3" json:"num_of_year,omitempty"`
}

func (x *Quote) Reset() {
	*x = Quote{}
	if protoimpl.UnsafeEnabled {
		mi := &file_repository_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Quote) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Quote) ProtoMessage() {}

func (x *Quote) ProtoReflect() protoreflect.Message {
	mi := &file_repository_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Quote.ProtoReflect.Descriptor instead.
func (*Quote) Descriptor() ([]byte, []int) {
	return file_repository_proto_rawDescGZIP(), []int{3}
}

func (x *Quote) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *Quote) GetOpen() float64 {
	if x != nil {
		return x.Open
	}
	return 0
}

func (x *Quote) GetClose() float64 {
	if x != nil {
		return x.Close
	}
	return 0
}

func (x *Quote) GetHigh() float64 {
	if x != nil {
		return x.High
	}
	return 0
}

func (x *Quote) GetLow() float64 {
	if x != nil {
		return x.Low
	}
	return 0
}

func (x *Quote) GetYesterdayClosed() float64 {
	if x != nil {
		return x.YesterdayClosed
	}
	return 0
}

func (x *Quote) GetVolume() uint64 {
	if x != nil {
		return x.Volume
	}
	return 0
}

func (x *Quote) GetAccount() float64 {
	if x != nil {
		return x.Account
	}
	return 0
}

func (x *Quote) GetDate() string {
	if x != nil {
		return x.Date
	}
	return ""
}

func (x *Quote) GetNumOfYear() int32 {
	if x != nil {
		return x.NumOfYear
	}
	return 0
}

var File_repository_proto protoreflect.FileDescriptor

var file_repository_proto_rawDesc = []byte{
	0x0a, 0x10, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x0a, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x1a, 0x1b,
	0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f,
	0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72, 0x61,
	0x70, 0x70, 0x65, 0x72, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x8e, 0x01, 0x0a, 0x07,
	0x43, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x12, 0x25, 0x0a, 0x0e, 0x61, 0x66, 0x66, 0x65, 0x63,
	0x74, 0x65, 0x64, 0x5f, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x0d, 0x61, 0x66, 0x66, 0x65, 0x63, 0x74, 0x65, 0x64, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x12, 0x2c,
	0x0a, 0x12, 0x61, 0x66, 0x66, 0x65, 0x63, 0x74, 0x65, 0x64, 0x5f, 0x71, 0x75, 0x6f, 0x74, 0x65,
	0x5f, 0x64, 0x61, 0x79, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x10, 0x61, 0x66, 0x66, 0x65,
	0x63, 0x74, 0x65, 0x64, 0x51, 0x75, 0x6f, 0x74, 0x65, 0x44, 0x61, 0x79, 0x12, 0x2e, 0x0a, 0x13,
	0x61, 0x66, 0x66, 0x65, 0x63, 0x74, 0x65, 0x64, 0x5f, 0x71, 0x75, 0x6f, 0x74, 0x65, 0x5f, 0x77,
	0x65, 0x65, 0x6b, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x11, 0x61, 0x66, 0x66, 0x65, 0x63,
	0x74, 0x65, 0x64, 0x51, 0x75, 0x6f, 0x74, 0x65, 0x57, 0x65, 0x65, 0x6b, 0x22, 0x49, 0x0a, 0x05,
	0x53, 0x74, 0x6f, 0x63, 0x6b, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a,
	0x07, 0x73, 0x75, 0x73, 0x70, 0x65, 0x6e, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07,
	0x73, 0x75, 0x73, 0x70, 0x65, 0x6e, 0x64, 0x22, 0x9a, 0x01, 0x0a, 0x0c, 0x51, 0x75, 0x6f, 0x74,
	0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x64, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x65,
	0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52,
	0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x31, 0x0a, 0x04, 0x6d, 0x6f, 0x64, 0x65, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x0e, 0x32, 0x1d, 0x2e, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72,
	0x79, 0x2e, 0x51, 0x75, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x2e, 0x4d,
	0x6f, 0x64, 0x65, 0x52, 0x04, 0x6d, 0x6f, 0x64, 0x65, 0x22, 0x19, 0x0a, 0x04, 0x4d, 0x6f, 0x64,
	0x65, 0x12, 0x07, 0x0a, 0x03, 0x44, 0x61, 0x79, 0x10, 0x00, 0x12, 0x08, 0x0a, 0x04, 0x57, 0x65,
	0x65, 0x6b, 0x10, 0x01, 0x22, 0xfc, 0x01, 0x0a, 0x05, 0x51, 0x75, 0x6f, 0x74, 0x65, 0x12, 0x12,
	0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f,
	0x64, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6f, 0x70, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x01,
	0x52, 0x04, 0x6f, 0x70, 0x65, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x18,
	0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x12, 0x12, 0x0a, 0x04,
	0x68, 0x69, 0x67, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x04, 0x68, 0x69, 0x67, 0x68,
	0x12, 0x10, 0x0a, 0x03, 0x6c, 0x6f, 0x77, 0x18, 0x05, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x6c,
	0x6f, 0x77, 0x12, 0x29, 0x0a, 0x10, 0x79, 0x65, 0x73, 0x74, 0x65, 0x72, 0x64, 0x61, 0x79, 0x5f,
	0x63, 0x6c, 0x6f, 0x73, 0x65, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0f, 0x79, 0x65,
	0x73, 0x74, 0x65, 0x72, 0x64, 0x61, 0x79, 0x43, 0x6c, 0x6f, 0x73, 0x65, 0x64, 0x12, 0x16, 0x0a,
	0x06, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01, 0x28, 0x04, 0x52, 0x06, 0x76,
	0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x18, 0x08, 0x20, 0x01, 0x28, 0x01, 0x52, 0x07, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x12,
	0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x64,
	0x61, 0x74, 0x65, 0x12, 0x1e, 0x0a, 0x0b, 0x6e, 0x75, 0x6d, 0x5f, 0x6f, 0x66, 0x5f, 0x79, 0x65,
	0x61, 0x72, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x6e, 0x75, 0x6d, 0x4f, 0x66, 0x59,
	0x65, 0x61, 0x72, 0x32, 0xd6, 0x01, 0x0a, 0x0a, 0x52, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f,
	0x72, 0x79, 0x12, 0x46, 0x0a, 0x0f, 0x41, 0x72, 0x63, 0x68, 0x69, 0x76, 0x65, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61,
	0x6c, 0x75, 0x65, 0x1a, 0x13, 0x2e, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79,
	0x2e, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x22, 0x00, 0x12, 0x3d, 0x0a, 0x0c, 0x47, 0x65,
	0x74, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x46, 0x75, 0x6c, 0x6c, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f,
	0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70,
	0x74, 0x79, 0x1a, 0x11, 0x2e, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x2e,
	0x53, 0x74, 0x6f, 0x63, 0x6b, 0x22, 0x00, 0x30, 0x01, 0x12, 0x41, 0x0a, 0x0e, 0x47, 0x65, 0x74,
	0x51, 0x75, 0x6f, 0x74, 0x65, 0x4c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x12, 0x18, 0x2e, 0x72, 0x65,
	0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f, 0x72, 0x79, 0x2e, 0x51, 0x75, 0x6f, 0x74, 0x65, 0x52, 0x65,
	0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x11, 0x2e, 0x72, 0x65, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x6f,
	0x72, 0x79, 0x2e, 0x51, 0x75, 0x6f, 0x74, 0x65, 0x22, 0x00, 0x30, 0x01, 0x42, 0x07, 0x5a, 0x05,
	0x2e, 0x2f, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_repository_proto_rawDescOnce sync.Once
	file_repository_proto_rawDescData = file_repository_proto_rawDesc
)

func file_repository_proto_rawDescGZIP() []byte {
	file_repository_proto_rawDescOnce.Do(func() {
		file_repository_proto_rawDescData = protoimpl.X.CompressGZIP(file_repository_proto_rawDescData)
	})
	return file_repository_proto_rawDescData
}

var file_repository_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_repository_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_repository_proto_goTypes = []interface{}{
	(QuoteRequest_Mode)(0),         // 0: repository.QuoteRequest.Mode
	(*Counter)(nil),                // 1: repository.Counter
	(*Stock)(nil),                  // 2: repository.Stock
	(*QuoteRequest)(nil),           // 3: repository.QuoteRequest
	(*Quote)(nil),                  // 4: repository.Quote
	(*wrapperspb.StringValue)(nil), // 5: google.protobuf.StringValue
	(*emptypb.Empty)(nil),          // 6: google.protobuf.Empty
}
var file_repository_proto_depIdxs = []int32{
	0, // 0: repository.QuoteRequest.mode:type_name -> repository.QuoteRequest.Mode
	5, // 1: repository.Repository.ArchiveMetadata:input_type -> google.protobuf.StringValue
	6, // 2: repository.Repository.GetStockFull:input_type -> google.protobuf.Empty
	3, // 3: repository.Repository.GetQuoteLatest:input_type -> repository.QuoteRequest
	1, // 4: repository.Repository.ArchiveMetadata:output_type -> repository.Counter
	2, // 5: repository.Repository.GetStockFull:output_type -> repository.Stock
	4, // 6: repository.Repository.GetQuoteLatest:output_type -> repository.Quote
	4, // [4:7] is the sub-list for method output_type
	1, // [1:4] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_repository_proto_init() }
func file_repository_proto_init() {
	if File_repository_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_repository_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Counter); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_repository_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Stock); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_repository_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*QuoteRequest); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_repository_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Quote); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_repository_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_repository_proto_goTypes,
		DependencyIndexes: file_repository_proto_depIdxs,
		EnumInfos:         file_repository_proto_enumTypes,
		MessageInfos:      file_repository_proto_msgTypes,
	}.Build()
	File_repository_proto = out.File
	file_repository_proto_rawDesc = nil
	file_repository_proto_goTypes = nil
	file_repository_proto_depIdxs = nil
}
