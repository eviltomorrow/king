// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.1
// 	protoc        v5.28.2
// source: assets.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
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

type Item_Type int32

const (
	Item_STOCK Item_Type = 0
	Item_ETF   Item_Type = 1
)

// Enum value maps for Item_Type.
var (
	Item_Type_name = map[int32]string{
		0: "STOCK",
		1: "ETF",
	}
	Item_Type_value = map[string]int32{
		"STOCK": 0,
		"ETF":   1,
	}
)

func (x Item_Type) Enum() *Item_Type {
	p := new(Item_Type)
	*p = x
	return p
}

func (x Item_Type) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (Item_Type) Descriptor() protoreflect.EnumDescriptor {
	return file_assets_proto_enumTypes[0].Descriptor()
}

func (Item_Type) Type() protoreflect.EnumType {
	return &file_assets_proto_enumTypes[0]
}

func (x Item_Type) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use Item_Type.Descriptor instead.
func (Item_Type) EnumDescriptor() ([]byte, []int) {
	return file_assets_proto_rawDescGZIP(), []int{1, 0}
}

type ItemListResp struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Items      []*Item `protobuf:"bytes,1,rep,name=items,proto3" json:"items,omitempty"`
	TotalCount int64   `protobuf:"varint,2,opt,name=total_count,json=totalCount,proto3" json:"total_count,omitempty"`
}

func (x *ItemListResp) Reset() {
	*x = ItemListResp{}
	mi := &file_assets_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *ItemListResp) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ItemListResp) ProtoMessage() {}

func (x *ItemListResp) ProtoReflect() protoreflect.Message {
	mi := &file_assets_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ItemListResp.ProtoReflect.Descriptor instead.
func (*ItemListResp) Descriptor() ([]byte, []int) {
	return file_assets_proto_rawDescGZIP(), []int{0}
}

func (x *ItemListResp) GetItems() []*Item {
	if x != nil {
		return x.Items
	}
	return nil
}

func (x *ItemListResp) GetTotalCount() int64 {
	if x != nil {
		return x.TotalCount
	}
	return 0
}

type Item struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	UserId           string    `protobuf:"bytes,1,opt,name=user_id,json=userId,proto3" json:"user_id,omitempty"`
	FundNo           string    `protobuf:"bytes,2,opt,name=fund_no,json=fundNo,proto3" json:"fund_no,omitempty"`
	Type             Item_Type `protobuf:"varint,3,opt,name=type,proto3,enum=account.Item_Type" json:"type,omitempty"`
	CashPosition     float64   `protobuf:"fixed64,4,opt,name=cash_position,json=cashPosition,proto3" json:"cash_position,omitempty"`
	Code             string    `protobuf:"bytes,5,opt,name=code,proto3" json:"code,omitempty"`
	Name             string    `protobuf:"bytes,6,opt,name=name,proto3" json:"name,omitempty"`
	OpenInterest     int64     `protobuf:"varint,7,opt,name=open_interest,json=openInterest,proto3" json:"open_interest,omitempty"`
	OpenId           string    `protobuf:"bytes,8,opt,name=open_id,json=openId,proto3" json:"open_id,omitempty"`
	FirstBuyDatetime int64     `protobuf:"varint,9,opt,name=first_buy_datetime,json=firstBuyDatetime,proto3" json:"first_buy_datetime,omitempty"`
}

func (x *Item) Reset() {
	*x = Item{}
	mi := &file_assets_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *Item) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Item) ProtoMessage() {}

func (x *Item) ProtoReflect() protoreflect.Message {
	mi := &file_assets_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Item.ProtoReflect.Descriptor instead.
func (*Item) Descriptor() ([]byte, []int) {
	return file_assets_proto_rawDescGZIP(), []int{1}
}

func (x *Item) GetUserId() string {
	if x != nil {
		return x.UserId
	}
	return ""
}

func (x *Item) GetFundNo() string {
	if x != nil {
		return x.FundNo
	}
	return ""
}

func (x *Item) GetType() Item_Type {
	if x != nil {
		return x.Type
	}
	return Item_STOCK
}

func (x *Item) GetCashPosition() float64 {
	if x != nil {
		return x.CashPosition
	}
	return 0
}

func (x *Item) GetCode() string {
	if x != nil {
		return x.Code
	}
	return ""
}

func (x *Item) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

func (x *Item) GetOpenInterest() int64 {
	if x != nil {
		return x.OpenInterest
	}
	return 0
}

func (x *Item) GetOpenId() string {
	if x != nil {
		return x.OpenId
	}
	return ""
}

func (x *Item) GetFirstBuyDatetime() int64 {
	if x != nil {
		return x.FirstBuyDatetime
	}
	return 0
}

var File_assets_proto protoreflect.FileDescriptor

var file_assets_proto_rawDesc = []byte{
	0x0a, 0x0c, 0x61, 0x73, 0x73, 0x65, 0x74, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x07,
	0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x1a, 0x1e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2f,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75, 0x66, 0x2f, 0x77, 0x72, 0x61, 0x70, 0x70, 0x65, 0x72,
	0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x54, 0x0a, 0x0c, 0x49, 0x74, 0x65, 0x6d, 0x4c,
	0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x12, 0x23, 0x0a, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73,
	0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0d, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x2e, 0x49, 0x74, 0x65, 0x6d, 0x52, 0x05, 0x69, 0x74, 0x65, 0x6d, 0x73, 0x12, 0x1f, 0x0a, 0x0b,
	0x74, 0x6f, 0x74, 0x61, 0x6c, 0x5f, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28,
	0x03, 0x52, 0x0a, 0x74, 0x6f, 0x74, 0x61, 0x6c, 0x43, 0x6f, 0x75, 0x6e, 0x74, 0x22, 0xb5, 0x02,
	0x0a, 0x04, 0x49, 0x74, 0x65, 0x6d, 0x12, 0x17, 0x0a, 0x07, 0x75, 0x73, 0x65, 0x72, 0x5f, 0x69,
	0x64, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x17, 0x0a, 0x07, 0x66, 0x75, 0x6e, 0x64, 0x5f, 0x6e, 0x6f, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09,
	0x52, 0x06, 0x66, 0x75, 0x6e, 0x64, 0x4e, 0x6f, 0x12, 0x26, 0x0a, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x18, 0x03, 0x20, 0x01, 0x28, 0x0e, 0x32, 0x12, 0x2e, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x2e, 0x49, 0x74, 0x65, 0x6d, 0x2e, 0x54, 0x79, 0x70, 0x65, 0x52, 0x04, 0x74, 0x79, 0x70, 0x65,
	0x12, 0x23, 0x0a, 0x0d, 0x63, 0x61, 0x73, 0x68, 0x5f, 0x70, 0x6f, 0x73, 0x69, 0x74, 0x69, 0x6f,
	0x6e, 0x18, 0x04, 0x20, 0x01, 0x28, 0x01, 0x52, 0x0c, 0x63, 0x61, 0x73, 0x68, 0x50, 0x6f, 0x73,
	0x69, 0x74, 0x69, 0x6f, 0x6e, 0x12, 0x12, 0x0a, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x04, 0x63, 0x6f, 0x64, 0x65, 0x12, 0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d,
	0x65, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x12, 0x23, 0x0a,
	0x0d, 0x6f, 0x70, 0x65, 0x6e, 0x5f, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x65, 0x73, 0x74, 0x18, 0x07,
	0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x6f, 0x70, 0x65, 0x6e, 0x49, 0x6e, 0x74, 0x65, 0x72, 0x65,
	0x73, 0x74, 0x12, 0x17, 0x0a, 0x07, 0x6f, 0x70, 0x65, 0x6e, 0x5f, 0x69, 0x64, 0x18, 0x08, 0x20,
	0x01, 0x28, 0x09, 0x52, 0x06, 0x6f, 0x70, 0x65, 0x6e, 0x49, 0x64, 0x12, 0x2c, 0x0a, 0x12, 0x66,
	0x69, 0x72, 0x73, 0x74, 0x5f, 0x62, 0x75, 0x79, 0x5f, 0x64, 0x61, 0x74, 0x65, 0x74, 0x69, 0x6d,
	0x65, 0x18, 0x09, 0x20, 0x01, 0x28, 0x03, 0x52, 0x10, 0x66, 0x69, 0x72, 0x73, 0x74, 0x42, 0x75,
	0x79, 0x44, 0x61, 0x74, 0x65, 0x74, 0x69, 0x6d, 0x65, 0x22, 0x1a, 0x0a, 0x04, 0x54, 0x79, 0x70,
	0x65, 0x12, 0x09, 0x0a, 0x05, 0x53, 0x54, 0x4f, 0x43, 0x4b, 0x10, 0x00, 0x12, 0x07, 0x0a, 0x03,
	0x45, 0x54, 0x46, 0x10, 0x01, 0x32, 0x4d, 0x0a, 0x06, 0x41, 0x73, 0x73, 0x65, 0x74, 0x73, 0x12,
	0x43, 0x0a, 0x0c, 0x4c, 0x69, 0x73, 0x74, 0x42, 0x79, 0x55, 0x73, 0x65, 0x72, 0x49, 0x64, 0x12,
	0x1c, 0x2e, 0x67, 0x6f, 0x6f, 0x67, 0x6c, 0x65, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x62, 0x75,
	0x66, 0x2e, 0x53, 0x74, 0x72, 0x69, 0x6e, 0x67, 0x56, 0x61, 0x6c, 0x75, 0x65, 0x1a, 0x15, 0x2e,
	0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x2e, 0x49, 0x74, 0x65, 0x6d, 0x4c, 0x69, 0x73, 0x74,
	0x52, 0x65, 0x73, 0x70, 0x42, 0x07, 0x5a, 0x05, 0x2e, 0x2f, 0x3b, 0x70, 0x62, 0x62, 0x06, 0x70,
	0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_assets_proto_rawDescOnce sync.Once
	file_assets_proto_rawDescData = file_assets_proto_rawDesc
)

func file_assets_proto_rawDescGZIP() []byte {
	file_assets_proto_rawDescOnce.Do(func() {
		file_assets_proto_rawDescData = protoimpl.X.CompressGZIP(file_assets_proto_rawDescData)
	})
	return file_assets_proto_rawDescData
}

var file_assets_proto_enumTypes = make([]protoimpl.EnumInfo, 1)
var file_assets_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_assets_proto_goTypes = []any{
	(Item_Type)(0),                 // 0: account.Item.Type
	(*ItemListResp)(nil),           // 1: account.ItemListResp
	(*Item)(nil),                   // 2: account.Item
	(*wrapperspb.StringValue)(nil), // 3: google.protobuf.StringValue
}
var file_assets_proto_depIdxs = []int32{
	2, // 0: account.ItemListResp.items:type_name -> account.Item
	0, // 1: account.Item.type:type_name -> account.Item.Type
	3, // 2: account.Assets.ListByUserId:input_type -> google.protobuf.StringValue
	1, // 3: account.Assets.ListByUserId:output_type -> account.ItemListResp
	3, // [3:4] is the sub-list for method output_type
	2, // [2:3] is the sub-list for method input_type
	2, // [2:2] is the sub-list for extension type_name
	2, // [2:2] is the sub-list for extension extendee
	0, // [0:2] is the sub-list for field type_name
}

func init() { file_assets_proto_init() }
func file_assets_proto_init() {
	if File_assets_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_assets_proto_rawDesc,
			NumEnums:      1,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_assets_proto_goTypes,
		DependencyIndexes: file_assets_proto_depIdxs,
		EnumInfos:         file_assets_proto_enumTypes,
		MessageInfos:      file_assets_proto_msgTypes,
	}.Build()
	File_assets_proto = out.File
	file_assets_proto_rawDesc = nil
	file_assets_proto_goTypes = nil
	file_assets_proto_depIdxs = nil
}