// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.26.0
// 	protoc        v3.19.4
// source: gozab.proto

package __

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

// used as ACK_E and parameters of NEWLEADER()
type EpochHist struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Epoch int32      `protobuf:"varint,1,opt,name=Epoch,proto3" json:"Epoch,omitempty"`
	Hist  []*PropTxn `protobuf:"bytes,2,rep,name=hist,proto3" json:"hist,omitempty"`
}

func (x *EpochHist) Reset() {
	*x = EpochHist{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gozab_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *EpochHist) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*EpochHist) ProtoMessage() {}

func (x *EpochHist) ProtoReflect() protoreflect.Message {
	mi := &file_gozab_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use EpochHist.ProtoReflect.Descriptor instead.
func (*EpochHist) Descriptor() ([]byte, []int) {
	return file_gozab_proto_rawDescGZIP(), []int{0}
}

func (x *EpochHist) GetEpoch() int32 {
	if x != nil {
		return x.Epoch
	}
	return 0
}

func (x *EpochHist) GetHist() []*PropTxn {
	if x != nil {
		return x.Hist
	}
	return nil
}

type Epoch struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Epoch int32 `protobuf:"varint,1,opt,name=epoch,proto3" json:"epoch,omitempty"`
}

func (x *Epoch) Reset() {
	*x = Epoch{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gozab_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Epoch) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Epoch) ProtoMessage() {}

func (x *Epoch) ProtoReflect() protoreflect.Message {
	mi := &file_gozab_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Epoch.ProtoReflect.Descriptor instead.
func (*Epoch) Descriptor() ([]byte, []int) {
	return file_gozab_proto_rawDescGZIP(), []int{1}
}

func (x *Epoch) GetEpoch() int32 {
	if x != nil {
		return x.Epoch
	}
	return 0
}

type Vote struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Voted bool `protobuf:"varint,1,opt,name=voted,proto3" json:"voted,omitempty"`
}

func (x *Vote) Reset() {
	*x = Vote{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gozab_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Vote) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Vote) ProtoMessage() {}

func (x *Vote) ProtoReflect() protoreflect.Message {
	mi := &file_gozab_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Vote.ProtoReflect.Descriptor instead.
func (*Vote) Descriptor() ([]byte, []int) {
	return file_gozab_proto_rawDescGZIP(), []int{2}
}

func (x *Vote) GetVoted() bool {
	if x != nil {
		return x.Voted
	}
	return false
}

type Vec struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key   string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
	Value int32  `protobuf:"varint,2,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *Vec) Reset() {
	*x = Vec{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gozab_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Vec) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Vec) ProtoMessage() {}

func (x *Vec) ProtoReflect() protoreflect.Message {
	mi := &file_gozab_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Vec.ProtoReflect.Descriptor instead.
func (*Vec) Descriptor() ([]byte, []int) {
	return file_gozab_proto_rawDescGZIP(), []int{3}
}

func (x *Vec) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

func (x *Vec) GetValue() int32 {
	if x != nil {
		return x.Value
	}
	return 0
}

type Zxid struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Epoch   int32 `protobuf:"varint,1,opt,name=epoch,proto3" json:"epoch,omitempty"`
	Counter int32 `protobuf:"varint,2,opt,name=counter,proto3" json:"counter,omitempty"`
}

func (x *Zxid) Reset() {
	*x = Zxid{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gozab_proto_msgTypes[4]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Zxid) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Zxid) ProtoMessage() {}

func (x *Zxid) ProtoReflect() protoreflect.Message {
	mi := &file_gozab_proto_msgTypes[4]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Zxid.ProtoReflect.Descriptor instead.
func (*Zxid) Descriptor() ([]byte, []int) {
	return file_gozab_proto_rawDescGZIP(), []int{4}
}

func (x *Zxid) GetEpoch() int32 {
	if x != nil {
		return x.Epoch
	}
	return 0
}

func (x *Zxid) GetCounter() int32 {
	if x != nil {
		return x.Counter
	}
	return 0
}

type Txn struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	V *Vec  `protobuf:"bytes,1,opt,name=v,proto3" json:"v,omitempty"`
	Z *Zxid `protobuf:"bytes,2,opt,name=z,proto3" json:"z,omitempty"`
}

func (x *Txn) Reset() {
	*x = Txn{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gozab_proto_msgTypes[5]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Txn) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Txn) ProtoMessage() {}

func (x *Txn) ProtoReflect() protoreflect.Message {
	mi := &file_gozab_proto_msgTypes[5]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Txn.ProtoReflect.Descriptor instead.
func (*Txn) Descriptor() ([]byte, []int) {
	return file_gozab_proto_rawDescGZIP(), []int{5}
}

func (x *Txn) GetV() *Vec {
	if x != nil {
		return x.V
	}
	return nil
}

func (x *Txn) GetZ() *Zxid {
	if x != nil {
		return x.Z
	}
	return nil
}

type PropTxn struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	E           int32 `protobuf:"varint,1,opt,name=e,proto3" json:"e,omitempty"`
	Transaction *Txn  `protobuf:"bytes,2,opt,name=transaction,proto3" json:"transaction,omitempty"`
}

func (x *PropTxn) Reset() {
	*x = PropTxn{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gozab_proto_msgTypes[6]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *PropTxn) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*PropTxn) ProtoMessage() {}

func (x *PropTxn) ProtoReflect() protoreflect.Message {
	mi := &file_gozab_proto_msgTypes[6]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use PropTxn.ProtoReflect.Descriptor instead.
func (*PropTxn) Descriptor() ([]byte, []int) {
	return file_gozab_proto_rawDescGZIP(), []int{6}
}

func (x *PropTxn) GetE() int32 {
	if x != nil {
		return x.E
	}
	return 0
}

func (x *PropTxn) GetTransaction() *Txn {
	if x != nil {
		return x.Transaction
	}
	return nil
}

type AckTxn struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Content string `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *AckTxn) Reset() {
	*x = AckTxn{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gozab_proto_msgTypes[7]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *AckTxn) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*AckTxn) ProtoMessage() {}

func (x *AckTxn) ProtoReflect() protoreflect.Message {
	mi := &file_gozab_proto_msgTypes[7]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use AckTxn.ProtoReflect.Descriptor instead.
func (*AckTxn) Descriptor() ([]byte, []int) {
	return file_gozab_proto_rawDescGZIP(), []int{7}
}

func (x *AckTxn) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

type CommitTxn struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Content string `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
	Epoch   int32  `protobuf:"varint,2,opt,name=epoch,proto3" json:"epoch,omitempty"`
}

func (x *CommitTxn) Reset() {
	*x = CommitTxn{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gozab_proto_msgTypes[8]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CommitTxn) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CommitTxn) ProtoMessage() {}

func (x *CommitTxn) ProtoReflect() protoreflect.Message {
	mi := &file_gozab_proto_msgTypes[8]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CommitTxn.ProtoReflect.Descriptor instead.
func (*CommitTxn) Descriptor() ([]byte, []int) {
	return file_gozab_proto_rawDescGZIP(), []int{8}
}

func (x *CommitTxn) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

func (x *CommitTxn) GetEpoch() int32 {
	if x != nil {
		return x.Epoch
	}
	return 0
}

type GetTxn struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Key string `protobuf:"bytes,1,opt,name=key,proto3" json:"key,omitempty"`
}

func (x *GetTxn) Reset() {
	*x = GetTxn{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gozab_proto_msgTypes[9]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *GetTxn) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*GetTxn) ProtoMessage() {}

func (x *GetTxn) ProtoReflect() protoreflect.Message {
	mi := &file_gozab_proto_msgTypes[9]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use GetTxn.ProtoReflect.Descriptor instead.
func (*GetTxn) Descriptor() ([]byte, []int) {
	return file_gozab_proto_rawDescGZIP(), []int{9}
}

func (x *GetTxn) GetKey() string {
	if x != nil {
		return x.Key
	}
	return ""
}

type ResultTxn struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Value int32 `protobuf:"varint,1,opt,name=value,proto3" json:"value,omitempty"`
}

func (x *ResultTxn) Reset() {
	*x = ResultTxn{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gozab_proto_msgTypes[10]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ResultTxn) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ResultTxn) ProtoMessage() {}

func (x *ResultTxn) ProtoReflect() protoreflect.Message {
	mi := &file_gozab_proto_msgTypes[10]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ResultTxn.ProtoReflect.Descriptor instead.
func (*ResultTxn) Descriptor() ([]byte, []int) {
	return file_gozab_proto_rawDescGZIP(), []int{10}
}

func (x *ResultTxn) GetValue() int32 {
	if x != nil {
		return x.Value
	}
	return 0
}

type Empty struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Content string `protobuf:"bytes,1,opt,name=content,proto3" json:"content,omitempty"`
}

func (x *Empty) Reset() {
	*x = Empty{}
	if protoimpl.UnsafeEnabled {
		mi := &file_gozab_proto_msgTypes[11]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *Empty) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*Empty) ProtoMessage() {}

func (x *Empty) ProtoReflect() protoreflect.Message {
	mi := &file_gozab_proto_msgTypes[11]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use Empty.ProtoReflect.Descriptor instead.
func (*Empty) Descriptor() ([]byte, []int) {
	return file_gozab_proto_rawDescGZIP(), []int{11}
}

func (x *Empty) GetContent() string {
	if x != nil {
		return x.Content
	}
	return ""
}

var File_gozab_proto protoreflect.FileDescriptor

var file_gozab_proto_rawDesc = []byte{
	0x0a, 0x0b, 0x67, 0x6f, 0x7a, 0x61, 0x62, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x05, 0x67,
	0x6f, 0x7a, 0x61, 0x62, 0x22, 0x45, 0x0a, 0x09, 0x45, 0x70, 0x6f, 0x63, 0x68, 0x48, 0x69, 0x73,
	0x74, 0x12, 0x14, 0x0a, 0x05, 0x45, 0x70, 0x6f, 0x63, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x05, 0x45, 0x70, 0x6f, 0x63, 0x68, 0x12, 0x22, 0x0a, 0x04, 0x68, 0x69, 0x73, 0x74, 0x18,
	0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x67, 0x6f, 0x7a, 0x61, 0x62, 0x2e, 0x50, 0x72,
	0x6f, 0x70, 0x54, 0x78, 0x6e, 0x52, 0x04, 0x68, 0x69, 0x73, 0x74, 0x22, 0x1d, 0x0a, 0x05, 0x45,
	0x70, 0x6f, 0x63, 0x68, 0x12, 0x14, 0x0a, 0x05, 0x65, 0x70, 0x6f, 0x63, 0x68, 0x18, 0x01, 0x20,
	0x01, 0x28, 0x05, 0x52, 0x05, 0x65, 0x70, 0x6f, 0x63, 0x68, 0x22, 0x1c, 0x0a, 0x04, 0x56, 0x6f,
	0x74, 0x65, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x6f, 0x74, 0x65, 0x64, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x08, 0x52, 0x05, 0x76, 0x6f, 0x74, 0x65, 0x64, 0x22, 0x2d, 0x0a, 0x03, 0x56, 0x65, 0x63, 0x12,
	0x10, 0x0a, 0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65,
	0x79, 0x12, 0x14, 0x0a, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05,
	0x52, 0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x22, 0x36, 0x0a, 0x04, 0x5a, 0x78, 0x69, 0x64, 0x12,
	0x14, 0x0a, 0x05, 0x65, 0x70, 0x6f, 0x63, 0x68, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05,
	0x65, 0x70, 0x6f, 0x63, 0x68, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72,
	0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x07, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x65, 0x72, 0x22,
	0x3a, 0x0a, 0x03, 0x54, 0x78, 0x6e, 0x12, 0x18, 0x0a, 0x01, 0x76, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x0b, 0x32, 0x0a, 0x2e, 0x67, 0x6f, 0x7a, 0x61, 0x62, 0x2e, 0x56, 0x65, 0x63, 0x52, 0x01, 0x76,
	0x12, 0x19, 0x0a, 0x01, 0x7a, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0b, 0x2e, 0x67, 0x6f,
	0x7a, 0x61, 0x62, 0x2e, 0x5a, 0x78, 0x69, 0x64, 0x52, 0x01, 0x7a, 0x22, 0x45, 0x0a, 0x07, 0x50,
	0x72, 0x6f, 0x70, 0x54, 0x78, 0x6e, 0x12, 0x0c, 0x0a, 0x01, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x05, 0x52, 0x01, 0x65, 0x12, 0x2c, 0x0a, 0x0b, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0a, 0x2e, 0x67, 0x6f, 0x7a, 0x61,
	0x62, 0x2e, 0x54, 0x78, 0x6e, 0x52, 0x0b, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69,
	0x6f, 0x6e, 0x22, 0x22, 0x0a, 0x06, 0x41, 0x63, 0x6b, 0x54, 0x78, 0x6e, 0x12, 0x18, 0x0a, 0x07,
	0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63,
	0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x22, 0x3b, 0x0a, 0x09, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74,
	0x54, 0x78, 0x6e, 0x12, 0x18, 0x0a, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x01,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x12, 0x14, 0x0a,
	0x05, 0x65, 0x70, 0x6f, 0x63, 0x68, 0x18, 0x02, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x65, 0x70,
	0x6f, 0x63, 0x68, 0x22, 0x1a, 0x0a, 0x06, 0x47, 0x65, 0x74, 0x54, 0x78, 0x6e, 0x12, 0x10, 0x0a,
	0x03, 0x6b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x22,
	0x21, 0x0a, 0x09, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x54, 0x78, 0x6e, 0x12, 0x14, 0x0a, 0x05,
	0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x01, 0x20, 0x01, 0x28, 0x05, 0x52, 0x05, 0x76, 0x61, 0x6c,
	0x75, 0x65, 0x22, 0x21, 0x0a, 0x05, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x12, 0x18, 0x0a, 0x07, 0x63,
	0x6f, 0x6e, 0x74, 0x65, 0x6e, 0x74, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x07, 0x63, 0x6f,
	0x6e, 0x74, 0x65, 0x6e, 0x74, 0x32, 0x95, 0x01, 0x0a, 0x0e, 0x46, 0x6f, 0x6c, 0x6c, 0x6f, 0x77,
	0x65, 0x72, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x2c, 0x0a, 0x09, 0x42, 0x72, 0x6f, 0x61,
	0x64, 0x63, 0x61, 0x73, 0x74, 0x12, 0x0e, 0x2e, 0x67, 0x6f, 0x7a, 0x61, 0x62, 0x2e, 0x50, 0x72,
	0x6f, 0x70, 0x54, 0x78, 0x6e, 0x1a, 0x0d, 0x2e, 0x67, 0x6f, 0x7a, 0x61, 0x62, 0x2e, 0x41, 0x63,
	0x6b, 0x54, 0x78, 0x6e, 0x22, 0x00, 0x12, 0x2a, 0x0a, 0x06, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74,
	0x12, 0x10, 0x2e, 0x67, 0x6f, 0x7a, 0x61, 0x62, 0x2e, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x54,
	0x78, 0x6e, 0x1a, 0x0c, 0x2e, 0x67, 0x6f, 0x7a, 0x61, 0x62, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79,
	0x22, 0x00, 0x12, 0x29, 0x0a, 0x09, 0x48, 0x65, 0x61, 0x72, 0x74, 0x42, 0x65, 0x61, 0x74, 0x12,
	0x0c, 0x2e, 0x67, 0x6f, 0x7a, 0x61, 0x62, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0c, 0x2e,
	0x67, 0x6f, 0x7a, 0x61, 0x62, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22, 0x00, 0x32, 0x89, 0x01,
	0x0a, 0x0a, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x55, 0x73, 0x65, 0x72, 0x12, 0x23, 0x0a, 0x05,
	0x53, 0x74, 0x6f, 0x72, 0x65, 0x12, 0x0a, 0x2e, 0x67, 0x6f, 0x7a, 0x61, 0x62, 0x2e, 0x56, 0x65,
	0x63, 0x1a, 0x0c, 0x2e, 0x67, 0x6f, 0x7a, 0x61, 0x62, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22,
	0x00, 0x12, 0x2d, 0x0a, 0x08, 0x52, 0x65, 0x74, 0x72, 0x69, 0x65, 0x76, 0x65, 0x12, 0x0d, 0x2e,
	0x67, 0x6f, 0x7a, 0x61, 0x62, 0x2e, 0x47, 0x65, 0x74, 0x54, 0x78, 0x6e, 0x1a, 0x10, 0x2e, 0x67,
	0x6f, 0x7a, 0x61, 0x62, 0x2e, 0x52, 0x65, 0x73, 0x75, 0x6c, 0x74, 0x54, 0x78, 0x6e, 0x22, 0x00,
	0x12, 0x27, 0x0a, 0x08, 0x49, 0x64, 0x65, 0x6e, 0x74, 0x69, 0x66, 0x79, 0x12, 0x0c, 0x2e, 0x67,
	0x6f, 0x7a, 0x61, 0x62, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x1a, 0x0b, 0x2e, 0x67, 0x6f, 0x7a,
	0x61, 0x62, 0x2e, 0x56, 0x6f, 0x74, 0x65, 0x22, 0x00, 0x32, 0xc5, 0x01, 0x0a, 0x0e, 0x56, 0x6f,
	0x74, 0x65, 0x72, 0x43, 0x61, 0x6e, 0x64, 0x69, 0x64, 0x61, 0x74, 0x65, 0x12, 0x26, 0x0a, 0x07,
	0x41, 0x73, 0x6b, 0x56, 0x6f, 0x74, 0x65, 0x12, 0x0c, 0x2e, 0x67, 0x6f, 0x7a, 0x61, 0x62, 0x2e,
	0x45, 0x70, 0x6f, 0x63, 0x68, 0x1a, 0x0b, 0x2e, 0x67, 0x6f, 0x7a, 0x61, 0x62, 0x2e, 0x56, 0x6f,
	0x74, 0x65, 0x22, 0x00, 0x12, 0x2c, 0x0a, 0x08, 0x4e, 0x65, 0x77, 0x45, 0x70, 0x6f, 0x63, 0x68,
	0x12, 0x0c, 0x2e, 0x67, 0x6f, 0x7a, 0x61, 0x62, 0x2e, 0x45, 0x70, 0x6f, 0x63, 0x68, 0x1a, 0x10,
	0x2e, 0x67, 0x6f, 0x7a, 0x61, 0x62, 0x2e, 0x45, 0x70, 0x6f, 0x63, 0x68, 0x48, 0x69, 0x73, 0x74,
	0x22, 0x00, 0x12, 0x2c, 0x0a, 0x09, 0x4e, 0x65, 0x77, 0x4c, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12,
	0x10, 0x2e, 0x67, 0x6f, 0x7a, 0x61, 0x62, 0x2e, 0x45, 0x70, 0x6f, 0x63, 0x68, 0x48, 0x69, 0x73,
	0x74, 0x1a, 0x0b, 0x2e, 0x67, 0x6f, 0x7a, 0x61, 0x62, 0x2e, 0x56, 0x6f, 0x74, 0x65, 0x22, 0x00,
	0x12, 0x2f, 0x0a, 0x0f, 0x43, 0x6f, 0x6d, 0x6d, 0x69, 0x74, 0x4e, 0x65, 0x77, 0x4c, 0x65, 0x61,
	0x64, 0x65, 0x72, 0x12, 0x0c, 0x2e, 0x67, 0x6f, 0x7a, 0x61, 0x62, 0x2e, 0x45, 0x70, 0x6f, 0x63,
	0x68, 0x1a, 0x0c, 0x2e, 0x67, 0x6f, 0x7a, 0x61, 0x62, 0x2e, 0x45, 0x6d, 0x70, 0x74, 0x79, 0x22,
	0x00, 0x42, 0x04, 0x5a, 0x02, 0x2e, 0x2f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_gozab_proto_rawDescOnce sync.Once
	file_gozab_proto_rawDescData = file_gozab_proto_rawDesc
)

func file_gozab_proto_rawDescGZIP() []byte {
	file_gozab_proto_rawDescOnce.Do(func() {
		file_gozab_proto_rawDescData = protoimpl.X.CompressGZIP(file_gozab_proto_rawDescData)
	})
	return file_gozab_proto_rawDescData
}

var file_gozab_proto_msgTypes = make([]protoimpl.MessageInfo, 12)
var file_gozab_proto_goTypes = []interface{}{
	(*EpochHist)(nil), // 0: gozab.EpochHist
	(*Epoch)(nil),     // 1: gozab.Epoch
	(*Vote)(nil),      // 2: gozab.Vote
	(*Vec)(nil),       // 3: gozab.Vec
	(*Zxid)(nil),      // 4: gozab.Zxid
	(*Txn)(nil),       // 5: gozab.Txn
	(*PropTxn)(nil),   // 6: gozab.PropTxn
	(*AckTxn)(nil),    // 7: gozab.AckTxn
	(*CommitTxn)(nil), // 8: gozab.CommitTxn
	(*GetTxn)(nil),    // 9: gozab.GetTxn
	(*ResultTxn)(nil), // 10: gozab.ResultTxn
	(*Empty)(nil),     // 11: gozab.Empty
}
var file_gozab_proto_depIdxs = []int32{
	6,  // 0: gozab.EpochHist.hist:type_name -> gozab.PropTxn
	3,  // 1: gozab.Txn.v:type_name -> gozab.Vec
	4,  // 2: gozab.Txn.z:type_name -> gozab.Zxid
	5,  // 3: gozab.PropTxn.transaction:type_name -> gozab.Txn
	6,  // 4: gozab.FollowerLeader.Broadcast:input_type -> gozab.PropTxn
	8,  // 5: gozab.FollowerLeader.Commit:input_type -> gozab.CommitTxn
	11, // 6: gozab.FollowerLeader.HeartBeat:input_type -> gozab.Empty
	3,  // 7: gozab.LeaderUser.Store:input_type -> gozab.Vec
	9,  // 8: gozab.LeaderUser.Retrieve:input_type -> gozab.GetTxn
	11, // 9: gozab.LeaderUser.Identify:input_type -> gozab.Empty
	1,  // 10: gozab.VoterCandidate.AskVote:input_type -> gozab.Epoch
	1,  // 11: gozab.VoterCandidate.NewEpoch:input_type -> gozab.Epoch
	0,  // 12: gozab.VoterCandidate.NewLeader:input_type -> gozab.EpochHist
	1,  // 13: gozab.VoterCandidate.CommitNewLeader:input_type -> gozab.Epoch
	7,  // 14: gozab.FollowerLeader.Broadcast:output_type -> gozab.AckTxn
	11, // 15: gozab.FollowerLeader.Commit:output_type -> gozab.Empty
	11, // 16: gozab.FollowerLeader.HeartBeat:output_type -> gozab.Empty
	11, // 17: gozab.LeaderUser.Store:output_type -> gozab.Empty
	10, // 18: gozab.LeaderUser.Retrieve:output_type -> gozab.ResultTxn
	2,  // 19: gozab.LeaderUser.Identify:output_type -> gozab.Vote
	2,  // 20: gozab.VoterCandidate.AskVote:output_type -> gozab.Vote
	0,  // 21: gozab.VoterCandidate.NewEpoch:output_type -> gozab.EpochHist
	2,  // 22: gozab.VoterCandidate.NewLeader:output_type -> gozab.Vote
	11, // 23: gozab.VoterCandidate.CommitNewLeader:output_type -> gozab.Empty
	14, // [14:24] is the sub-list for method output_type
	4,  // [4:14] is the sub-list for method input_type
	4,  // [4:4] is the sub-list for extension type_name
	4,  // [4:4] is the sub-list for extension extendee
	0,  // [0:4] is the sub-list for field type_name
}

func init() { file_gozab_proto_init() }
func file_gozab_proto_init() {
	if File_gozab_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_gozab_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*EpochHist); i {
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
		file_gozab_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Epoch); i {
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
		file_gozab_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Vote); i {
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
		file_gozab_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Vec); i {
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
		file_gozab_proto_msgTypes[4].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Zxid); i {
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
		file_gozab_proto_msgTypes[5].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Txn); i {
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
		file_gozab_proto_msgTypes[6].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*PropTxn); i {
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
		file_gozab_proto_msgTypes[7].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*AckTxn); i {
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
		file_gozab_proto_msgTypes[8].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CommitTxn); i {
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
		file_gozab_proto_msgTypes[9].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*GetTxn); i {
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
		file_gozab_proto_msgTypes[10].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ResultTxn); i {
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
		file_gozab_proto_msgTypes[11].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*Empty); i {
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
			RawDescriptor: file_gozab_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   12,
			NumExtensions: 0,
			NumServices:   3,
		},
		GoTypes:           file_gozab_proto_goTypes,
		DependencyIndexes: file_gozab_proto_depIdxs,
		MessageInfos:      file_gozab_proto_msgTypes,
	}.Build()
	File_gozab_proto = out.File
	file_gozab_proto_rawDesc = nil
	file_gozab_proto_goTypes = nil
	file_gozab_proto_depIdxs = nil
}
