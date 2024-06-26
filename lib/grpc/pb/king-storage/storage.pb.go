// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.34.2
// 	protoc        v5.27.1
// source: storage.proto

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
	return file_storage_proto_enumTypes[0].Descriptor()
}

func (QuoteRequest_Mode) Type() protoreflect.EnumType {
	return &file_storage_proto_enumTypes[0]
}

func (x QuoteRequest_Mode) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use QuoteRequest_Mode.Descriptor instead.
func (QuoteRequest_Mode) EnumDescriptor() ([]byte, []int) {
	return file_storage_proto_rawDescGZIP(), []int{2, 0}
}

type Stats struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	StockAffected     int64 `protobuf:"varint,1,opt,name=stock_affected,json=stockAffected,proto3" json:"stock_affected,omitempty"`
	QuoteDayAffected  int64 `protobuf:"varint,2,opt,name=quote_day_affected,json=quoteDayAffected,proto3" json:"quote_day_affected,omitempty"`
	QuoteWeekAffected int64 `protobuf:"varint,3,opt,name=quote_week_affected,json=quoteWeekAffected,proto3" json:"quote_week_affected,omitempty"`
}

func (x *Stats) Reset() {
	*x = Stats{}
	if protoimpl.UnsafeEnabled {
		mi := &file_storage_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Stats) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Stats) ProtoMessage() {}

func (x *Stats) ProtoReflect() protoreflect.Message {
	mi := &file_storage_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Stats.ProtoReflect.Descriptor instead.
func (*Stats) Descriptor() ([]byte, []int) {
	return file_storage_proto_rawDescGZIP(), []int{0}
}

func (x *Stats) GetStockAffected() int64 {
	if x != nil {
		return x.StockAffected
	}
	return 0
}

func (x *Stats) GetQuoteDayAffected() int64 {
	if x != nil {
		return x.QuoteDayAffected
	}
	return 0
}

func (x *Stats) GetQuoteWeekAffected() int64 {
	if x != nil {
		return x.QuoteWeekAffected
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
		mi := &file_storage_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Stock) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Stock) ProtoMessage() {}

func (x *Stock) ProtoReflect() protoreflect.Message {
	mi := &file_storage_proto_msgTypes[1]
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
	return file_storage_proto_rawDescGZIP(), []int{1}
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
	Mode  QuoteRequest_Mode `protobuf:"varint,4,opt,name=mode,proto3,enum=storage.QuoteRequest_Mode" json:"mode,omitempty"`
}

func (x *QuoteRequest) Reset() {
	*x = QuoteRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_storage_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *QuoteRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*QuoteRequest) ProtoMessage() {}

func (x *QuoteRequest) ProtoReflect() protoreflect.Message {
	mi := &file_storage_proto_msgTypes[2]
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
	return file_storage_proto_rawDescGZIP(), []int{2}
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
		mi := &file_storage_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Quote) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Quote) ProtoMessage() {}

func (x *Quote) ProtoReflect() protoreflect.Message {
	mi := &file_storage_proto_msgTypes[3]
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
	return file_storage_proto_rawDescGZIP(), []int{3}
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

type Metadata struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Source          string  `protobuf:"bytes,1,opt,name=source,proto3" json:"source,omitempty"`
	Code            string  `protobuf:"bytes,2,opt,name=code,proto3" json:"code,omitempty"`
	Name            string  `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
	Open            float64 `protobuf:"fixed64,4,opt,name=open,proto3" json:"open,omitempty"`
	YesterdayClosed float64 `protobuf:"fixed64,5,opt,name=yesterday_closed,json=yesterdayClosed,proto3" json:"yesterday_closed,omitempty"`
	Latest          float64 `protobuf:"fixed64,6,opt,name=latest,proto3" json:"latest,omitempty"`
	High            float64 `protobuf:"fixed64,7,opt,name=high,proto3" json:"high,omitempty"`
	Low             float64 `protobuf:"fixed64,8,opt,name=low,proto3" json:"low,omitempty"`
	Volume          uint64  `protobuf:"varint,9,opt,name=volume,proto3" json:"volume,omitempty"`
	Account         float64 `protobuf:"fixed64,10,opt,name=account,proto3" json:"account,omitempty"`
	Date            string  `protobuf:"bytes,11,opt,name=date,proto3" json:"date,omitempty"`
	Time            string  `protobuf:"bytes,12,opt,name=time,proto3" json:"time,omitempty"`
	Suspend         string  `protobuf:"bytes,13,opt,name=suspend,proto3" json:"suspend,omitempty"`
}

func (x *Metadata) Reset() {
	*x = Metadata{}
	if protoimpl.UnsafeEnabled {
		mi := &file_storage_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Metadata) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Metadata) ProtoMessage() {}

func (x *Metadata) ProtoReflect() protoreflect.Message {
	mi := &file_storage_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Metadata.ProtoReflect.Descriptor instead.
func (*Metadata) Descriptor() ([]byte, []int) {
	return file_storage_proto_rawDescGZIP(), []int{4}
}

func (x *Metadata) GetSource() string {
	if x != nil {
		return x.Source
	}
	return ""
}

func (x *Metadata) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *Metadata) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Metadata) GetOpen() float64 {
	if x != nil {
		return x.Open
	}
	return 0
}

func (x *Metadata) GetYesterdayClosed() float64 {
	if x != nil {
		return x.YesterdayClosed
	}
	return 0
}

func (x *Metadata) GetLatest() float64 {
	if x != nil {
		return x.Latest
	}
	return 0
}

func (x *Metadata) GetHigh() float64 {
	if x != nil {
		return x.High
	}
	return 0
}

func (x *Metadata) GetLow() float64 {
	if x != nil {
		return x.Low
	}
	return 0
}

func (x *Metadata) GetVolume() uint64 {
	if x != nil {
		return x.Volume
	}
	return 0
}

func (x *Metadata) GetAccount() float64 {
	if x != nil {
		return x.Account
	}
	return 0
}

func (x *Metadata) GetDate() string {
	if x != nil {
		return x.Date
	}
	return ""
}

func (x *Metadata) GetTime() string {
	if x != nil {
		return x.Time
	}
	return ""
}

func (x *Metadata) GetSuspend() string {
	if x != nil {
		return x.Suspend
	}
	return ""
}

var File_storage_proto protoreflect.FileDescriptor

var file_storage_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x1a, 0x1b, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65,
	0x2f, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x65, 0x6d, 0x70, 0x74, 0x79, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72, 0x73, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x8c, 0x01, 0x0a, 0x05, 0x53, 0x74, 0x61, 0x74, 0x73, 0x12,
	0x25, 0x0a, 0x0e, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x5f, 0x61, 0x66, 0x66, 0x65, 0x63, 0x74, 0x65,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0d, 0x73, 0x74, 0x6f, 0x63, 0x6b, 0x41, 0x66,
	0x66, 0x65, 0x63, 0x74, 0x65, 0x64, 0x12, 0x2c, 0x0a, 0x12, 0x71, 0x75, 0x6f, 0x74, 0x65, 0x5f,
	0x64, 0x61, 0x79, 0x5f, 0x61, 0x66, 0x66, 0x65, 0x63, 0x74, 0x65, 0x64, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x03, 0x52, 0x10, 0x71, 0x75, 0x6f, 0x74, 0x65, 0x44, 0x61, 0x79, 0x41, 0x66, 0x66, 0x65,
	0x63, 0x74, 0x65, 0x64, 0x12, 0x2e, 0x0a, 0x13, 0x71, 0x75, 0x6f, 0x74, 0x65, 0x5f, 0x77, 0x65,
	0x65, 0x6b, 0x5f, 0x61, 0x66, 0x66, 0x65, 0x63, 0x74, 0x65, 0x64, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x11, 0x71, 0x75, 0x6f, 0x74, 0x65, 0x57, 0x65, 0x65, 0x6b, 0x41, 0x66, 0x66, 0x65,
	0x63, 0x74, 0x65, 0x64, 0x22, 0x49, 0x0a, 0x05, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x12, 0x12, 0x0a,
	0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x73, 0x75, 0x73, 0x70, 0x65, 0x6e, 0x64,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73, 0x75, 0x73, 0x70, 0x65, 0x6e, 0x64, 0x22,
	0x97, 0x01, 0x0a, 0x0c, 0x51, 0x75, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74,
	0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04,
	0x63, 0x6f, 0x64, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x65, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x6c, 0x69, 0x6d, 0x69,
	0x74, 0x18, 0x03, 0x20, 0x01, 0x28, 0x03, 0x52, 0x05, 0x6c, 0x69, 0x6d, 0x69, 0x74, 0x12, 0x2e,
	0x0a, 0x04, 0x6d, 0x6f, 0x64, 0x65, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x1a, 0x2e, 0x73,
	0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x51, 0x75, 0x6f, 0x74, 0x65, 0x52, 0x65, 0x71, 0x75,
	0x65, 0x73, 0x74, 0x2e, 0x4d, 0x6f, 0x64, 0x65, 0x52, 0x04, 0x6d, 0x6f, 0x64, 0x65, 0x22, 0x19,
	0x0a, 0x04, 0x4d, 0x6f, 0x64, 0x65, 0x12, 0x07, 0x0a, 0x03, 0x44, 0x61, 0x79, 0x10, 0x00, 0x12,
	0x08, 0x0a, 0x04, 0x57, 0x65, 0x65, 0x6b, 0x10, 0x01, 0x22, 0xfc, 0x01, 0x0a, 0x05, 0x51, 0x75,
	0x6f, 0x74, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6f, 0x70, 0x65, 0x6e, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x01, 0x52, 0x04, 0x6f, 0x70, 0x65, 0x6e, 0x12, 0x14, 0x0a, 0x05, 0x63,
	0x6c, 0x6f, 0x73, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x01, 0x52, 0x05, 0x63, 0x6c, 0x6f, 0x73,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x68, 0x69, 0x67, 0x68, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52,
	0x04, 0x68, 0x69, 0x67, 0x68, 0x12, 0x10, 0x0a, 0x03, 0x6c, 0x6f, 0x77, 0x18, 0x05, 0x20, 0x01,
	0x28, 0x01, 0x52, 0x03, 0x6c, 0x6f, 0x77, 0x12, 0x29, 0x0a, 0x10, 0x79, 0x65, 0x73, 0x74, 0x65,
	0x72, 0x64, 0x61, 0x79, 0x5f, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x64, 0x18, 0x06, 0x20, 0x01, 0x28,
	0x01, 0x52, 0x0f, 0x79, 0x65, 0x73, 0x74, 0x65, 0x72, 0x64, 0x61, 0x79, 0x43, 0x6c, 0x6f, 0x73,
	0x65, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x18, 0x07, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x06, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x08, 0x20, 0x01, 0x28, 0x01, 0x52, 0x07, 0x61, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x65, 0x18, 0x09, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x65, 0x12, 0x1e, 0x0a, 0x0b, 0x6e, 0x75, 0x6d, 0x5f,
	0x6f, 0x66, 0x5f, 0x79, 0x65, 0x61, 0x72, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x05, 0x52, 0x09, 0x6e,
	0x75, 0x6d, 0x4f, 0x66, 0x59, 0x65, 0x61, 0x72, 0x22, 0xbb, 0x02, 0x0a, 0x08, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x12, 0x16, 0x0a, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x18,
	0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x73, 0x6f, 0x75, 0x72, 0x63, 0x65, 0x12, 0x12, 0x0a,
	0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64,
	0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52,
	0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6f, 0x70, 0x65, 0x6e, 0x18, 0x04, 0x20,
	0x01, 0x28, 0x01, 0x52, 0x04, 0x6f, 0x70, 0x65, 0x6e, 0x12, 0x29, 0x0a, 0x10, 0x79, 0x65, 0x73,
	0x74, 0x65, 0x72, 0x64, 0x61, 0x79, 0x5f, 0x63, 0x6c, 0x6f, 0x73, 0x65, 0x64, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x01, 0x52, 0x0f, 0x79, 0x65, 0x73, 0x74, 0x65, 0x72, 0x64, 0x61, 0x79, 0x43, 0x6c,
	0x6f, 0x73, 0x65, 0x64, 0x12, 0x16, 0x0a, 0x06, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x18, 0x06,
	0x20, 0x01, 0x28, 0x01, 0x52, 0x06, 0x6c, 0x61, 0x74, 0x65, 0x73, 0x74, 0x12, 0x12, 0x0a, 0x04,
	0x68, 0x69, 0x67, 0x68, 0x18, 0x07, 0x20, 0x01, 0x28, 0x01, 0x52, 0x04, 0x68, 0x69, 0x67, 0x68,
	0x12, 0x10, 0x0a, 0x03, 0x6c, 0x6f, 0x77, 0x18, 0x08, 0x20, 0x01, 0x28, 0x01, 0x52, 0x03, 0x6c,
	0x6f, 0x77, 0x12, 0x16, 0x0a, 0x06, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x18, 0x09, 0x20, 0x01,
	0x28, 0x04, 0x52, 0x06, 0x76, 0x6f, 0x6c, 0x75, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07, 0x61, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x0a, 0x20, 0x01, 0x28, 0x01, 0x52, 0x07, 0x61, 0x63, 0x63,
	0x6f, 0x75, 0x6e, 0x74, 0x12, 0x12, 0x0a, 0x04, 0x64, 0x61, 0x74, 0x65, 0x18, 0x0b, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x04, 0x64, 0x61, 0x74, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x74, 0x69, 0x6d, 0x65,
	0x18, 0x0c, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x74, 0x69, 0x6d, 0x65, 0x12, 0x18, 0x0a, 0x07,
	0x73, 0x75, 0x73, 0x70, 0x65, 0x6e, 0x64, 0x18, 0x0d, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x73,
	0x75, 0x73, 0x70, 0x65, 0x6e, 0x64, 0x32, 0xf0, 0x01, 0x0a, 0x07, 0x53, 0x74, 0x6f, 0x72, 0x61,
	0x67, 0x65, 0x12, 0x33, 0x0a, 0x0c, 0x50, 0x75, 0x73, 0x68, 0x4d, 0x65, 0x74, 0x61, 0x64, 0x61,
	0x74, 0x61, 0x12, 0x11, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x4d, 0x65, 0x74,
	0x61, 0x64, 0x61, 0x74, 0x61, 0x1a, 0x0e, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e,
	0x53, 0x74, 0x61, 0x74, 0x73, 0x28, 0x01, 0x12, 0x3b, 0x0a, 0x0b, 0x47, 0x65, 0x74, 0x53, 0x74,
	0x6f, 0x63, 0x6b, 0x4f, 0x6e, 0x65, 0x12, 0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56,
	0x61, 0x6c, 0x75, 0x65, 0x1a, 0x0e, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x53,
	0x74, 0x6f, 0x63, 0x6b, 0x12, 0x38, 0x0a, 0x0c, 0x47, 0x65, 0x74, 0x53, 0x74, 0x6f, 0x63, 0x6b,
	0x46, 0x75, 0x6c, 0x6c, 0x12, 0x16, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0e, 0x2e, 0x73,
	0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x53, 0x74, 0x6f, 0x63, 0x6b, 0x30, 0x01, 0x12, 0x39,
	0x0a, 0x0e, 0x47, 0x65, 0x74, 0x51, 0x75, 0x6f, 0x74, 0x65, 0x4c, 0x61, 0x74, 0x65, 0x73, 0x74,
	0x12, 0x15, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67, 0x65, 0x2e, 0x51, 0x75, 0x6f, 0x74, 0x65,
	0x52, 0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x0e, 0x2e, 0x73, 0x74, 0x6f, 0x72, 0x61, 0x67,
	0x65, 0x2e, 0x51, 0x75, 0x6f, 0x74, 0x65, 0x30, 0x01, 0x42, 0x07, 0x5a, 0x05, 0x2e, 0x2f, 0x3b,
	0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_storage_proto_rawDescOnce sync.Once
	file_storage_proto_rawDescData = file_storage_proto_rawDesc
)

func file_storage_proto_rawDescGZIP() []byte {
	file_storage_proto_rawDescOnce.Do(func() {
		file_storage_proto_rawDescData = protoimpl.X.CompressGZIP(file_storage_proto_rawDescData)
	})
	return file_storage_proto_rawDescData
}

var file_storage_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_storage_proto_msgTypes = make([]protoimpl.MessageInfo, 5)
var file_storage_proto_goTypes = []any{
	(QuoteRequest_Mode)(0),         // 0: storage.QuoteRequest.Mode
	(*Stats)(nil),                  // 1: storage.Stats
	(*Stock)(nil),                  // 2: storage.Stock
	(*QuoteRequest)(nil),           // 3: storage.QuoteRequest
	(*Quote)(nil),                  // 4: storage.Quote
	(*Metadata)(nil),               // 5: storage.Metadata
	(*wrapperspb.StringValue)(nil), // 6: google.protobuf.StringValue
	(*emptypb.Empty)(nil),          // 7: google.protobuf.Empty
}
var file_storage_proto_depIdxs = []int32{
	0, // 0: storage.QuoteRequest.mode:type_name -> storage.QuoteRequest.Mode
	5, // 1: storage.Storage.PushMetadata:input_type -> storage.Metadata
	6, // 2: storage.Storage.GetStockOne:input_type -> google.protobuf.StringValue
	7, // 3: storage.Storage.GetStockFull:input_type -> google.protobuf.Empty
	3, // 4: storage.Storage.GetQuoteLatest:input_type -> storage.QuoteRequest
	1, // 5: storage.Storage.PushMetadata:output_type -> storage.Stats
	2, // 6: storage.Storage.GetStockOne:output_type -> storage.Stock
	2, // 7: storage.Storage.GetStockFull:output_type -> storage.Stock
	4, // 8: storage.Storage.GetQuoteLatest:output_type -> storage.Quote
	5, // [5:9] is the sub-list for method output_type
	1, // [1:5] is the sub-list for method input_type
	1, // [1:1] is the sub-list for extension type_name
	1, // [1:1] is the sub-list for extension extendee
	0, // [0:1] is the sub-list for field type_name
}

func init() { file_storage_proto_init() }
func file_storage_proto_init() {
	if File_storage_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_storage_proto_msgTypes[0].Exporter = func(v any, i int) any {
			switch v := v.(*Stats); i {
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
		file_storage_proto_msgTypes[1].Exporter = func(v any, i int) any {
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
		file_storage_proto_msgTypes[2].Exporter = func(v any, i int) any {
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
		file_storage_proto_msgTypes[3].Exporter = func(v any, i int) any {
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
		file_storage_proto_msgTypes[4].Exporter = func(v any, i int) any {
			switch v := v.(*Metadata); i {
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
			RawDescriptor: file_storage_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   5,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_storage_proto_goTypes,
		DependencyIndexes: file_storage_proto_depIdxs,
		EnumInfos:         file_storage_proto_enumTypes,
		MessageInfos:      file_storage_proto_msgTypes,
	}.Build()
	File_storage_proto = out.File
	file_storage_proto_rawDesc = nil
	file_storage_proto_goTypes = nil
	file_storage_proto_depIdxs = nil
}
