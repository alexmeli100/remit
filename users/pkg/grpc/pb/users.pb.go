// Code generated by protoc-gen-go. DO NOT EDIT.
// source: users.proto

package pb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

type User struct {
	FirstName            string   `protobuf:"bytes,1,opt,name=firstName,proto3" json:"firstName,omitempty"`
	LastName             string   `protobuf:"bytes,2,opt,name=lastName,proto3" json:"lastName,omitempty"`
	Email                string   `protobuf:"bytes,3,opt,name=email,proto3" json:"email,omitempty"`
	Password             string   `protobuf:"bytes,4,opt,name=password,proto3" json:"password,omitempty"`
	Id                   int64    `protobuf:"varint,5,opt,name=id,proto3" json:"id,omitempty"`
	Address              string   `protobuf:"bytes,6,opt,name=address,proto3" json:"address,omitempty"`
	Confirmed            bool     `protobuf:"varint,7,opt,name=confirmed,proto3" json:"confirmed,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *User) Reset()         { *m = User{} }
func (m *User) String() string { return proto.CompactTextString(m) }
func (*User) ProtoMessage()    {}
func (*User) Descriptor() ([]byte, []int) {
	return fileDescriptor_030765f334c86cea, []int{0}
}

func (m *User) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_User.Unmarshal(m, b)
}
func (m *User) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_User.Marshal(b, m, deterministic)
}
func (m *User) XXX_Merge(src proto.Message) {
	xxx_messageInfo_User.Merge(m, src)
}
func (m *User) XXX_Size() int {
	return xxx_messageInfo_User.Size(m)
}
func (m *User) XXX_DiscardUnknown() {
	xxx_messageInfo_User.DiscardUnknown(m)
}

var xxx_messageInfo_User proto.InternalMessageInfo

func (m *User) GetFirstName() string {
	if m != nil {
		return m.FirstName
	}
	return ""
}

func (m *User) GetLastName() string {
	if m != nil {
		return m.LastName
	}
	return ""
}

func (m *User) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

func (m *User) GetPassword() string {
	if m != nil {
		return m.Password
	}
	return ""
}

func (m *User) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

func (m *User) GetAddress() string {
	if m != nil {
		return m.Address
	}
	return ""
}

func (m *User) GetConfirmed() bool {
	if m != nil {
		return m.Confirmed
	}
	return false
}

type CreateRequest struct {
	User                 *User    `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateRequest) Reset()         { *m = CreateRequest{} }
func (m *CreateRequest) String() string { return proto.CompactTextString(m) }
func (*CreateRequest) ProtoMessage()    {}
func (*CreateRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_030765f334c86cea, []int{1}
}

func (m *CreateRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateRequest.Unmarshal(m, b)
}
func (m *CreateRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateRequest.Marshal(b, m, deterministic)
}
func (m *CreateRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateRequest.Merge(m, src)
}
func (m *CreateRequest) XXX_Size() int {
	return xxx_messageInfo_CreateRequest.Size(m)
}
func (m *CreateRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateRequest.DiscardUnknown(m)
}

var xxx_messageInfo_CreateRequest proto.InternalMessageInfo

func (m *CreateRequest) GetUser() *User {
	if m != nil {
		return m.User
	}
	return nil
}

type CreateReply struct {
	Err                  string   `protobuf:"bytes,1,opt,name=err,proto3" json:"err,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *CreateReply) Reset()         { *m = CreateReply{} }
func (m *CreateReply) String() string { return proto.CompactTextString(m) }
func (*CreateReply) ProtoMessage()    {}
func (*CreateReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_030765f334c86cea, []int{2}
}

func (m *CreateReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_CreateReply.Unmarshal(m, b)
}
func (m *CreateReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_CreateReply.Marshal(b, m, deterministic)
}
func (m *CreateReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_CreateReply.Merge(m, src)
}
func (m *CreateReply) XXX_Size() int {
	return xxx_messageInfo_CreateReply.Size(m)
}
func (m *CreateReply) XXX_DiscardUnknown() {
	xxx_messageInfo_CreateReply.DiscardUnknown(m)
}

var xxx_messageInfo_CreateReply proto.InternalMessageInfo

func (m *CreateReply) GetErr() string {
	if m != nil {
		return m.Err
	}
	return ""
}

type GetUserByIDRequest struct {
	Id                   int64    `protobuf:"varint,1,opt,name=id,proto3" json:"id,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetUserByIDRequest) Reset()         { *m = GetUserByIDRequest{} }
func (m *GetUserByIDRequest) String() string { return proto.CompactTextString(m) }
func (*GetUserByIDRequest) ProtoMessage()    {}
func (*GetUserByIDRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_030765f334c86cea, []int{3}
}

func (m *GetUserByIDRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetUserByIDRequest.Unmarshal(m, b)
}
func (m *GetUserByIDRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetUserByIDRequest.Marshal(b, m, deterministic)
}
func (m *GetUserByIDRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetUserByIDRequest.Merge(m, src)
}
func (m *GetUserByIDRequest) XXX_Size() int {
	return xxx_messageInfo_GetUserByIDRequest.Size(m)
}
func (m *GetUserByIDRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetUserByIDRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetUserByIDRequest proto.InternalMessageInfo

func (m *GetUserByIDRequest) GetId() int64 {
	if m != nil {
		return m.Id
	}
	return 0
}

type GetUserByIDReply struct {
	User                 *User    `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	Err                  string   `protobuf:"bytes,2,opt,name=err,proto3" json:"err,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetUserByIDReply) Reset()         { *m = GetUserByIDReply{} }
func (m *GetUserByIDReply) String() string { return proto.CompactTextString(m) }
func (*GetUserByIDReply) ProtoMessage()    {}
func (*GetUserByIDReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_030765f334c86cea, []int{4}
}

func (m *GetUserByIDReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetUserByIDReply.Unmarshal(m, b)
}
func (m *GetUserByIDReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetUserByIDReply.Marshal(b, m, deterministic)
}
func (m *GetUserByIDReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetUserByIDReply.Merge(m, src)
}
func (m *GetUserByIDReply) XXX_Size() int {
	return xxx_messageInfo_GetUserByIDReply.Size(m)
}
func (m *GetUserByIDReply) XXX_DiscardUnknown() {
	xxx_messageInfo_GetUserByIDReply.DiscardUnknown(m)
}

var xxx_messageInfo_GetUserByIDReply proto.InternalMessageInfo

func (m *GetUserByIDReply) GetUser() *User {
	if m != nil {
		return m.User
	}
	return nil
}

func (m *GetUserByIDReply) GetErr() string {
	if m != nil {
		return m.Err
	}
	return ""
}

type GetUserByEmailRequest struct {
	Email                string   `protobuf:"bytes,1,opt,name=email,proto3" json:"email,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetUserByEmailRequest) Reset()         { *m = GetUserByEmailRequest{} }
func (m *GetUserByEmailRequest) String() string { return proto.CompactTextString(m) }
func (*GetUserByEmailRequest) ProtoMessage()    {}
func (*GetUserByEmailRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_030765f334c86cea, []int{5}
}

func (m *GetUserByEmailRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetUserByEmailRequest.Unmarshal(m, b)
}
func (m *GetUserByEmailRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetUserByEmailRequest.Marshal(b, m, deterministic)
}
func (m *GetUserByEmailRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetUserByEmailRequest.Merge(m, src)
}
func (m *GetUserByEmailRequest) XXX_Size() int {
	return xxx_messageInfo_GetUserByEmailRequest.Size(m)
}
func (m *GetUserByEmailRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_GetUserByEmailRequest.DiscardUnknown(m)
}

var xxx_messageInfo_GetUserByEmailRequest proto.InternalMessageInfo

func (m *GetUserByEmailRequest) GetEmail() string {
	if m != nil {
		return m.Email
	}
	return ""
}

type GetUserByEmailReply struct {
	User                 *User    `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	Err                  string   `protobuf:"bytes,2,opt,name=err,proto3" json:"err,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *GetUserByEmailReply) Reset()         { *m = GetUserByEmailReply{} }
func (m *GetUserByEmailReply) String() string { return proto.CompactTextString(m) }
func (*GetUserByEmailReply) ProtoMessage()    {}
func (*GetUserByEmailReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_030765f334c86cea, []int{6}
}

func (m *GetUserByEmailReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_GetUserByEmailReply.Unmarshal(m, b)
}
func (m *GetUserByEmailReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_GetUserByEmailReply.Marshal(b, m, deterministic)
}
func (m *GetUserByEmailReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_GetUserByEmailReply.Merge(m, src)
}
func (m *GetUserByEmailReply) XXX_Size() int {
	return xxx_messageInfo_GetUserByEmailReply.Size(m)
}
func (m *GetUserByEmailReply) XXX_DiscardUnknown() {
	xxx_messageInfo_GetUserByEmailReply.DiscardUnknown(m)
}

var xxx_messageInfo_GetUserByEmailReply proto.InternalMessageInfo

func (m *GetUserByEmailReply) GetUser() *User {
	if m != nil {
		return m.User
	}
	return nil
}

func (m *GetUserByEmailReply) GetErr() string {
	if m != nil {
		return m.Err
	}
	return ""
}

type UpdateEmailRequest struct {
	User                 *User    `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateEmailRequest) Reset()         { *m = UpdateEmailRequest{} }
func (m *UpdateEmailRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateEmailRequest) ProtoMessage()    {}
func (*UpdateEmailRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_030765f334c86cea, []int{7}
}

func (m *UpdateEmailRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateEmailRequest.Unmarshal(m, b)
}
func (m *UpdateEmailRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateEmailRequest.Marshal(b, m, deterministic)
}
func (m *UpdateEmailRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateEmailRequest.Merge(m, src)
}
func (m *UpdateEmailRequest) XXX_Size() int {
	return xxx_messageInfo_UpdateEmailRequest.Size(m)
}
func (m *UpdateEmailRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateEmailRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateEmailRequest proto.InternalMessageInfo

func (m *UpdateEmailRequest) GetUser() *User {
	if m != nil {
		return m.User
	}
	return nil
}

type UpdateEmailReply struct {
	Err                  string   `protobuf:"bytes,1,opt,name=err,proto3" json:"err,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateEmailReply) Reset()         { *m = UpdateEmailReply{} }
func (m *UpdateEmailReply) String() string { return proto.CompactTextString(m) }
func (*UpdateEmailReply) ProtoMessage()    {}
func (*UpdateEmailReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_030765f334c86cea, []int{8}
}

func (m *UpdateEmailReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateEmailReply.Unmarshal(m, b)
}
func (m *UpdateEmailReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateEmailReply.Marshal(b, m, deterministic)
}
func (m *UpdateEmailReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateEmailReply.Merge(m, src)
}
func (m *UpdateEmailReply) XXX_Size() int {
	return xxx_messageInfo_UpdateEmailReply.Size(m)
}
func (m *UpdateEmailReply) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateEmailReply.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateEmailReply proto.InternalMessageInfo

func (m *UpdateEmailReply) GetErr() string {
	if m != nil {
		return m.Err
	}
	return ""
}

type UpdatePasswordRequest struct {
	User                 *User    `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdatePasswordRequest) Reset()         { *m = UpdatePasswordRequest{} }
func (m *UpdatePasswordRequest) String() string { return proto.CompactTextString(m) }
func (*UpdatePasswordRequest) ProtoMessage()    {}
func (*UpdatePasswordRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_030765f334c86cea, []int{9}
}

func (m *UpdatePasswordRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdatePasswordRequest.Unmarshal(m, b)
}
func (m *UpdatePasswordRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdatePasswordRequest.Marshal(b, m, deterministic)
}
func (m *UpdatePasswordRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdatePasswordRequest.Merge(m, src)
}
func (m *UpdatePasswordRequest) XXX_Size() int {
	return xxx_messageInfo_UpdatePasswordRequest.Size(m)
}
func (m *UpdatePasswordRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdatePasswordRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UpdatePasswordRequest proto.InternalMessageInfo

func (m *UpdatePasswordRequest) GetUser() *User {
	if m != nil {
		return m.User
	}
	return nil
}

type UpdatePasswordReply struct {
	Err                  string   `protobuf:"bytes,1,opt,name=err,proto3" json:"err,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdatePasswordReply) Reset()         { *m = UpdatePasswordReply{} }
func (m *UpdatePasswordReply) String() string { return proto.CompactTextString(m) }
func (*UpdatePasswordReply) ProtoMessage()    {}
func (*UpdatePasswordReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_030765f334c86cea, []int{10}
}

func (m *UpdatePasswordReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdatePasswordReply.Unmarshal(m, b)
}
func (m *UpdatePasswordReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdatePasswordReply.Marshal(b, m, deterministic)
}
func (m *UpdatePasswordReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdatePasswordReply.Merge(m, src)
}
func (m *UpdatePasswordReply) XXX_Size() int {
	return xxx_messageInfo_UpdatePasswordReply.Size(m)
}
func (m *UpdatePasswordReply) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdatePasswordReply.DiscardUnknown(m)
}

var xxx_messageInfo_UpdatePasswordReply proto.InternalMessageInfo

func (m *UpdatePasswordReply) GetErr() string {
	if m != nil {
		return m.Err
	}
	return ""
}

type UpdateStatusRequest struct {
	User                 *User    `protobuf:"bytes,1,opt,name=user,proto3" json:"user,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateStatusRequest) Reset()         { *m = UpdateStatusRequest{} }
func (m *UpdateStatusRequest) String() string { return proto.CompactTextString(m) }
func (*UpdateStatusRequest) ProtoMessage()    {}
func (*UpdateStatusRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_030765f334c86cea, []int{11}
}

func (m *UpdateStatusRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateStatusRequest.Unmarshal(m, b)
}
func (m *UpdateStatusRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateStatusRequest.Marshal(b, m, deterministic)
}
func (m *UpdateStatusRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateStatusRequest.Merge(m, src)
}
func (m *UpdateStatusRequest) XXX_Size() int {
	return xxx_messageInfo_UpdateStatusRequest.Size(m)
}
func (m *UpdateStatusRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateStatusRequest.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateStatusRequest proto.InternalMessageInfo

func (m *UpdateStatusRequest) GetUser() *User {
	if m != nil {
		return m.User
	}
	return nil
}

type UpdateStatusReply struct {
	Err                  string   `protobuf:"bytes,1,opt,name=err,proto3" json:"err,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *UpdateStatusReply) Reset()         { *m = UpdateStatusReply{} }
func (m *UpdateStatusReply) String() string { return proto.CompactTextString(m) }
func (*UpdateStatusReply) ProtoMessage()    {}
func (*UpdateStatusReply) Descriptor() ([]byte, []int) {
	return fileDescriptor_030765f334c86cea, []int{12}
}

func (m *UpdateStatusReply) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_UpdateStatusReply.Unmarshal(m, b)
}
func (m *UpdateStatusReply) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_UpdateStatusReply.Marshal(b, m, deterministic)
}
func (m *UpdateStatusReply) XXX_Merge(src proto.Message) {
	xxx_messageInfo_UpdateStatusReply.Merge(m, src)
}
func (m *UpdateStatusReply) XXX_Size() int {
	return xxx_messageInfo_UpdateStatusReply.Size(m)
}
func (m *UpdateStatusReply) XXX_DiscardUnknown() {
	xxx_messageInfo_UpdateStatusReply.DiscardUnknown(m)
}

var xxx_messageInfo_UpdateStatusReply proto.InternalMessageInfo

func (m *UpdateStatusReply) GetErr() string {
	if m != nil {
		return m.Err
	}
	return ""
}

func init() {
	proto.RegisterType((*User)(nil), "pb.User")
	proto.RegisterType((*CreateRequest)(nil), "pb.CreateRequest")
	proto.RegisterType((*CreateReply)(nil), "pb.CreateReply")
	proto.RegisterType((*GetUserByIDRequest)(nil), "pb.GetUserByIDRequest")
	proto.RegisterType((*GetUserByIDReply)(nil), "pb.GetUserByIDReply")
	proto.RegisterType((*GetUserByEmailRequest)(nil), "pb.GetUserByEmailRequest")
	proto.RegisterType((*GetUserByEmailReply)(nil), "pb.GetUserByEmailReply")
	proto.RegisterType((*UpdateEmailRequest)(nil), "pb.UpdateEmailRequest")
	proto.RegisterType((*UpdateEmailReply)(nil), "pb.UpdateEmailReply")
	proto.RegisterType((*UpdatePasswordRequest)(nil), "pb.UpdatePasswordRequest")
	proto.RegisterType((*UpdatePasswordReply)(nil), "pb.UpdatePasswordReply")
	proto.RegisterType((*UpdateStatusRequest)(nil), "pb.UpdateStatusRequest")
	proto.RegisterType((*UpdateStatusReply)(nil), "pb.UpdateStatusReply")
}

func init() {
	proto.RegisterFile("users.proto", fileDescriptor_030765f334c86cea)
}

var fileDescriptor_030765f334c86cea = []byte{
	// 437 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x9c, 0x54, 0x5d, 0x6b, 0xe2, 0x40,
	0x14, 0x25, 0xf1, 0xfb, 0x66, 0xd7, 0xd5, 0x51, 0x77, 0xb3, 0x41, 0x58, 0x09, 0x2e, 0xeb, 0xc3,
	0xea, 0x83, 0xd2, 0xa7, 0x42, 0x1f, 0xac, 0x52, 0xfa, 0x52, 0x4a, 0x8a, 0x3f, 0x60, 0x6c, 0x46,
	0x08, 0xc4, 0x26, 0x9d, 0x19, 0x29, 0xfe, 0xb5, 0xfe, 0xb5, 0xbe, 0x94, 0xc9, 0x64, 0x34, 0x63,
	0x82, 0x48, 0xdf, 0x32, 0xf7, 0x9e, 0x73, 0xee, 0x99, 0x7b, 0x86, 0x80, 0xb5, 0x63, 0x84, 0xb2,
	0x49, 0x4c, 0x23, 0x1e, 0x21, 0x33, 0x5e, 0xbb, 0xef, 0x06, 0x94, 0x57, 0x8c, 0x50, 0xd4, 0x87,
	0xc6, 0x26, 0xa0, 0x8c, 0x3f, 0xe0, 0x2d, 0xb1, 0x8d, 0x81, 0x31, 0x6a, 0x78, 0xc7, 0x02, 0x72,
	0xa0, 0x1e, 0xe2, 0xb4, 0x69, 0x26, 0xcd, 0xc3, 0x19, 0x75, 0xa1, 0x42, 0xb6, 0x38, 0x08, 0xed,
	0x52, 0xd2, 0x90, 0x07, 0xc1, 0x88, 0x31, 0x63, 0x6f, 0x11, 0xf5, 0xed, 0xb2, 0x64, 0xa8, 0x33,
	0x6a, 0x82, 0x19, 0xf8, 0x76, 0x65, 0x60, 0x8c, 0x4a, 0x9e, 0x19, 0xf8, 0xc8, 0x86, 0x1a, 0xf6,
	0x7d, 0x4a, 0x18, 0xb3, 0xab, 0x09, 0x54, 0x1d, 0x85, 0xab, 0xe7, 0xe8, 0x65, 0x13, 0xd0, 0x2d,
	0xf1, 0xed, 0xda, 0xc0, 0x18, 0xd5, 0xbd, 0x63, 0xc1, 0x1d, 0xc3, 0xf7, 0x5b, 0x4a, 0x30, 0x27,
	0x1e, 0x79, 0xdd, 0x11, 0xc6, 0x51, 0x1f, 0xca, 0xe2, 0x82, 0x89, 0x7f, 0x6b, 0x5a, 0x9f, 0xc4,
	0xeb, 0x89, 0xb8, 0x9c, 0x97, 0x54, 0xdd, 0x3f, 0x60, 0x29, 0x78, 0x1c, 0xee, 0x51, 0x0b, 0x4a,
	0x84, 0xd2, 0xf4, 0xae, 0xe2, 0xd3, 0x1d, 0x02, 0xba, 0x23, 0x5c, 0x30, 0xe6, 0xfb, 0xfb, 0x85,
	0x12, 0x95, 0x6e, 0x0d, 0xe5, 0xd6, 0x9d, 0x43, 0x4b, 0x43, 0x09, 0xad, 0xb3, 0x83, 0xd5, 0x24,
	0xf3, 0x38, 0x69, 0x0c, 0xbd, 0x83, 0xc6, 0x52, 0xec, 0x4b, 0x0d, 0x3b, 0x2c, 0xd3, 0xc8, 0x2c,
	0xd3, 0x5d, 0x42, 0xe7, 0x14, 0xfe, 0x95, 0xa9, 0x53, 0x40, 0xab, 0xd8, 0xc7, 0x9c, 0x68, 0x23,
	0xcf, 0x2f, 0x6d, 0x08, 0x2d, 0x8d, 0x53, 0xbc, 0xb9, 0x2b, 0xe8, 0x49, 0xd4, 0x63, 0x9a, 0xf1,
	0x65, 0xe2, 0xff, 0xa0, 0x73, 0x4a, 0x2b, 0xd6, 0x9f, 0x29, 0xe0, 0x13, 0xc7, 0x7c, 0xc7, 0x2e,
	0x53, 0xff, 0x0b, 0x6d, 0x9d, 0x54, 0xa8, 0x3d, 0xfd, 0x30, 0xa1, 0x22, 0x58, 0x0c, 0xfd, 0x87,
	0xaa, 0x7c, 0x20, 0xa8, 0x2d, 0xa4, 0xb4, 0xb7, 0xe5, 0xfc, 0xc8, 0x96, 0x84, 0xd2, 0x35, 0x58,
	0x99, 0x77, 0x80, 0x7e, 0x8a, 0x7e, 0xfe, 0xf9, 0x38, 0xdd, 0x5c, 0x5d, 0x90, 0x17, 0xd0, 0xd4,
	0x13, 0x45, 0xbf, 0x35, 0x5c, 0x36, 0x21, 0xe7, 0x57, 0x51, 0x2b, 0xb5, 0x90, 0x09, 0x47, 0x5a,
	0xc8, 0x27, 0x2c, 0x2d, 0xe4, 0x52, 0x5c, 0x40, 0x53, 0x5f, 0xbe, 0xb4, 0x50, 0x98, 0xa3, 0xb4,
	0x50, 0x94, 0xd5, 0x0d, 0x7c, 0xcb, 0x2e, 0x19, 0x65, 0x80, 0x5a, 0x56, 0x4e, 0x2f, 0xdf, 0x88,
	0xc3, 0xfd, 0xba, 0x9a, 0xfc, 0x8b, 0x66, 0x9f, 0x01, 0x00, 0x00, 0xff, 0xff, 0x71, 0x2c, 0xb0,
	0xae, 0x9a, 0x04, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConnInterface

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion6

// UsersClient is the client API for Users service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type UsersClient interface {
	Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateReply, error)
	GetUserByID(ctx context.Context, in *GetUserByIDRequest, opts ...grpc.CallOption) (*GetUserByIDReply, error)
	GetUserByEmail(ctx context.Context, in *GetUserByEmailRequest, opts ...grpc.CallOption) (*GetUserByEmailReply, error)
	UpdateEmail(ctx context.Context, in *UpdateEmailRequest, opts ...grpc.CallOption) (*UpdateEmailReply, error)
	UpdatePassword(ctx context.Context, in *UpdatePasswordRequest, opts ...grpc.CallOption) (*UpdatePasswordReply, error)
	UpdateStatus(ctx context.Context, in *UpdateStatusRequest, opts ...grpc.CallOption) (*UpdateStatusReply, error)
}

type usersClient struct {
	cc grpc.ClientConnInterface
}

func NewUsersClient(cc grpc.ClientConnInterface) UsersClient {
	return &usersClient{cc}
}

func (c *usersClient) Create(ctx context.Context, in *CreateRequest, opts ...grpc.CallOption) (*CreateReply, error) {
	out := new(CreateReply)
	err := c.cc.Invoke(ctx, "/pb.Users/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersClient) GetUserByID(ctx context.Context, in *GetUserByIDRequest, opts ...grpc.CallOption) (*GetUserByIDReply, error) {
	out := new(GetUserByIDReply)
	err := c.cc.Invoke(ctx, "/pb.Users/GetUserByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersClient) GetUserByEmail(ctx context.Context, in *GetUserByEmailRequest, opts ...grpc.CallOption) (*GetUserByEmailReply, error) {
	out := new(GetUserByEmailReply)
	err := c.cc.Invoke(ctx, "/pb.Users/GetUserByEmail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersClient) UpdateEmail(ctx context.Context, in *UpdateEmailRequest, opts ...grpc.CallOption) (*UpdateEmailReply, error) {
	out := new(UpdateEmailReply)
	err := c.cc.Invoke(ctx, "/pb.Users/UpdateEmail", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersClient) UpdatePassword(ctx context.Context, in *UpdatePasswordRequest, opts ...grpc.CallOption) (*UpdatePasswordReply, error) {
	out := new(UpdatePasswordReply)
	err := c.cc.Invoke(ctx, "/pb.Users/UpdatePassword", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *usersClient) UpdateStatus(ctx context.Context, in *UpdateStatusRequest, opts ...grpc.CallOption) (*UpdateStatusReply, error) {
	out := new(UpdateStatusReply)
	err := c.cc.Invoke(ctx, "/pb.Users/UpdateStatus", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UsersServer is the server API for Users service.
type UsersServer interface {
	Create(context.Context, *CreateRequest) (*CreateReply, error)
	GetUserByID(context.Context, *GetUserByIDRequest) (*GetUserByIDReply, error)
	GetUserByEmail(context.Context, *GetUserByEmailRequest) (*GetUserByEmailReply, error)
	UpdateEmail(context.Context, *UpdateEmailRequest) (*UpdateEmailReply, error)
	UpdatePassword(context.Context, *UpdatePasswordRequest) (*UpdatePasswordReply, error)
	UpdateStatus(context.Context, *UpdateStatusRequest) (*UpdateStatusReply, error)
}

// UnimplementedUsersServer can be embedded to have forward compatible implementations.
type UnimplementedUsersServer struct {
}

func (*UnimplementedUsersServer) Create(ctx context.Context, req *CreateRequest) (*CreateReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (*UnimplementedUsersServer) GetUserByID(ctx context.Context, req *GetUserByIDRequest) (*GetUserByIDReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserByID not implemented")
}
func (*UnimplementedUsersServer) GetUserByEmail(ctx context.Context, req *GetUserByEmailRequest) (*GetUserByEmailReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserByEmail not implemented")
}
func (*UnimplementedUsersServer) UpdateEmail(ctx context.Context, req *UpdateEmailRequest) (*UpdateEmailReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateEmail not implemented")
}
func (*UnimplementedUsersServer) UpdatePassword(ctx context.Context, req *UpdatePasswordRequest) (*UpdatePasswordReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePassword not implemented")
}
func (*UnimplementedUsersServer) UpdateStatus(ctx context.Context, req *UpdateStatusRequest) (*UpdateStatusReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateStatus not implemented")
}

func RegisterUsersServer(s *grpc.Server, srv UsersServer) {
	s.RegisterService(&_Users_serviceDesc, srv)
}

func _Users_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Users/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).Create(ctx, req.(*CreateRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Users_GetUserByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserByIDRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).GetUserByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Users/GetUserByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).GetUserByID(ctx, req.(*GetUserByIDRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Users_GetUserByEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserByEmailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).GetUserByEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Users/GetUserByEmail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).GetUserByEmail(ctx, req.(*GetUserByEmailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Users_UpdateEmail_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateEmailRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).UpdateEmail(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Users/UpdateEmail",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).UpdateEmail(ctx, req.(*UpdateEmailRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Users_UpdatePassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).UpdatePassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Users/UpdatePassword",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).UpdatePassword(ctx, req.(*UpdatePasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Users_UpdateStatus_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateStatusRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UsersServer).UpdateStatus(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/pb.Users/UpdateStatus",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UsersServer).UpdateStatus(ctx, req.(*UpdateStatusRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var _Users_serviceDesc = grpc.ServiceDesc{
	ServiceName: "pb.Users",
	HandlerType: (*UsersServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _Users_Create_Handler,
		},
		{
			MethodName: "GetUserByID",
			Handler:    _Users_GetUserByID_Handler,
		},
		{
			MethodName: "GetUserByEmail",
			Handler:    _Users_GetUserByEmail_Handler,
		},
		{
			MethodName: "UpdateEmail",
			Handler:    _Users_UpdateEmail_Handler,
		},
		{
			MethodName: "UpdatePassword",
			Handler:    _Users_UpdatePassword_Handler,
		},
		{
			MethodName: "UpdateStatus",
			Handler:    _Users_UpdateStatus_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "users.proto",
}
