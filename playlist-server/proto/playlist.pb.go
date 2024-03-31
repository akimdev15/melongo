// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v4.25.3
// source: playlist.proto

package proto

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

type CreatePlaylistRequest struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	ApiKey       string `protobuf:"bytes,1,opt,name=apiKey,proto3" json:"apiKey,omitempty"`
	AccessToken  string `protobuf:"bytes,2,opt,name=accessToken,proto3" json:"accessToken,omitempty"`
	PlaylistName string `protobuf:"bytes,3,opt,name=playlistName,proto3" json:"playlistName,omitempty"`
	Description  string `protobuf:"bytes,4,opt,name=description,proto3" json:"description,omitempty"`
	IsPublic     bool   `protobuf:"varint,5,opt,name=isPublic,proto3" json:"isPublic,omitempty"`
	UserID       string `protobuf:"bytes,6,opt,name=userID,proto3" json:"userID,omitempty"`
}

func (x *CreatePlaylistRequest) Reset() {
	*x = CreatePlaylistRequest{}
	if protoimpl.UnsafeEnabled {
		mi := &file_playlist_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreatePlaylistRequest) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreatePlaylistRequest) ProtoMessage() {}

func (x *CreatePlaylistRequest) ProtoReflect() protoreflect.Message {
	mi := &file_playlist_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreatePlaylistRequest.ProtoReflect.Descriptor instead.
func (*CreatePlaylistRequest) Descriptor() ([]byte, []int) {
	return file_playlist_proto_rawDescGZIP(), []int{0}
}

func (x *CreatePlaylistRequest) GetApiKey() string {
	if x != nil {
		return x.ApiKey
	}
	return ""
}

func (x *CreatePlaylistRequest) GetAccessToken() string {
	if x != nil {
		return x.AccessToken
	}
	return ""
}

func (x *CreatePlaylistRequest) GetPlaylistName() string {
	if x != nil {
		return x.PlaylistName
	}
	return ""
}

func (x *CreatePlaylistRequest) GetDescription() string {
	if x != nil {
		return x.Description
	}
	return ""
}

func (x *CreatePlaylistRequest) GetIsPublic() bool {
	if x != nil {
		return x.IsPublic
	}
	return false
}

func (x *CreatePlaylistRequest) GetUserID() string {
	if x != nil {
		return x.UserID
	}
	return ""
}

type CreatePlaylistResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	SpotifyPlaylistID string `protobuf:"bytes,1,opt,name=spotifyPlaylistID,proto3" json:"spotifyPlaylistID,omitempty"`
	ExternalURL       string `protobuf:"bytes,2,opt,name=externalURL,proto3" json:"externalURL,omitempty"`
	Name              string `protobuf:"bytes,3,opt,name=name,proto3" json:"name,omitempty"`
}

func (x *CreatePlaylistResponse) Reset() {
	*x = CreatePlaylistResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_playlist_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *CreatePlaylistResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*CreatePlaylistResponse) ProtoMessage() {}

func (x *CreatePlaylistResponse) ProtoReflect() protoreflect.Message {
	mi := &file_playlist_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use CreatePlaylistResponse.ProtoReflect.Descriptor instead.
func (*CreatePlaylistResponse) Descriptor() ([]byte, []int) {
	return file_playlist_proto_rawDescGZIP(), []int{1}
}

func (x *CreatePlaylistResponse) GetSpotifyPlaylistID() string {
	if x != nil {
		return x.SpotifyPlaylistID
	}
	return ""
}

func (x *CreatePlaylistResponse) GetExternalURL() string {
	if x != nil {
		return x.ExternalURL
	}
	return ""
}

func (x *CreatePlaylistResponse) GetName() string {
	if x != nil {
		return x.Name
	}
	return ""
}

var File_playlist_proto protoreflect.FileDescriptor

var file_playlist_proto_rawDesc = []byte{
	0x0a, 0x0e, 0x70, 0x6c, 0x61, 0x79, 0x6c, 0x69, 0x73, 0x74, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xcb, 0x01, 0x0a, 0x15, 0x43, 0x72, 0x65, 0x61,
	0x74, 0x65, 0x50, 0x6c, 0x61, 0x79, 0x6c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x71, 0x75, 0x65, 0x73,
	0x74, 0x12, 0x16, 0x0a, 0x06, 0x61, 0x70, 0x69, 0x4b, 0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x06, 0x61, 0x70, 0x69, 0x4b, 0x65, 0x79, 0x12, 0x20, 0x0a, 0x0b, 0x61, 0x63, 0x63,
	0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x0b,
	0x61, 0x63, 0x63, 0x65, 0x73, 0x73, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x12, 0x22, 0x0a, 0x0c, 0x70,
	0x6c, 0x61, 0x79, 0x6c, 0x69, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28,
	0x09, 0x52, 0x0c, 0x70, 0x6c, 0x61, 0x79, 0x6c, 0x69, 0x73, 0x74, 0x4e, 0x61, 0x6d, 0x65, 0x12,
	0x20, 0x0a, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f, 0x6e, 0x18, 0x04,
	0x20, 0x01, 0x28, 0x09, 0x52, 0x0b, 0x64, 0x65, 0x73, 0x63, 0x72, 0x69, 0x70, 0x74, 0x69, 0x6f,
	0x6e, 0x12, 0x1a, 0x0a, 0x08, 0x69, 0x73, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x18, 0x05, 0x20,
	0x01, 0x28, 0x08, 0x52, 0x08, 0x69, 0x73, 0x50, 0x75, 0x62, 0x6c, 0x69, 0x63, 0x12, 0x16, 0x0a,
	0x06, 0x75, 0x73, 0x65, 0x72, 0x49, 0x44, 0x18, 0x06, 0x20, 0x01, 0x28, 0x09, 0x52, 0x06, 0x75,
	0x73, 0x65, 0x72, 0x49, 0x44, 0x22, 0x7c, 0x0a, 0x16, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x50,
	0x6c, 0x61, 0x79, 0x6c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x2c, 0x0a, 0x11, 0x73, 0x70, 0x6f, 0x74, 0x69, 0x66, 0x79, 0x50, 0x6c, 0x61, 0x79, 0x6c, 0x69,
	0x73, 0x74, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x11, 0x73, 0x70, 0x6f, 0x74,
	0x69, 0x66, 0x79, 0x50, 0x6c, 0x61, 0x79, 0x6c, 0x69, 0x73, 0x74, 0x49, 0x44, 0x12, 0x20, 0x0a,
	0x0b, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x55, 0x52, 0x4c, 0x18, 0x02, 0x20, 0x01,
	0x28, 0x09, 0x52, 0x0b, 0x65, 0x78, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x55, 0x52, 0x4c, 0x12,
	0x12, 0x0a, 0x04, 0x6e, 0x61, 0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x09, 0x52, 0x04, 0x6e,
	0x61, 0x6d, 0x65, 0x32, 0x60, 0x0a, 0x0f, 0x50, 0x6c, 0x61, 0x79, 0x6c, 0x69, 0x73, 0x74, 0x53,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x12, 0x4d, 0x0a, 0x0e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65,
	0x50, 0x6c, 0x61, 0x79, 0x6c, 0x69, 0x73, 0x74, 0x12, 0x1c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x43, 0x72, 0x65, 0x61, 0x74, 0x65, 0x50, 0x6c, 0x61, 0x79, 0x6c, 0x69, 0x73, 0x74, 0x52,
	0x65, 0x71, 0x75, 0x65, 0x73, 0x74, 0x1a, 0x1d, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x43,
	0x72, 0x65, 0x61, 0x74, 0x65, 0x50, 0x6c, 0x61, 0x79, 0x6c, 0x69, 0x73, 0x74, 0x52, 0x65, 0x73,
	0x70, 0x6f, 0x6e, 0x73, 0x65, 0x42, 0x0a, 0x5a, 0x08, 0x2e, 0x2e, 0x2f, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_playlist_proto_rawDescOnce sync.Once
	file_playlist_proto_rawDescData = file_playlist_proto_rawDesc
)

func file_playlist_proto_rawDescGZIP() []byte {
	file_playlist_proto_rawDescOnce.Do(func() {
		file_playlist_proto_rawDescData = protoimpl.X.CompressGZIP(file_playlist_proto_rawDescData)
	})
	return file_playlist_proto_rawDescData
}

var file_playlist_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_playlist_proto_goTypes = []interface{}{
	(*CreatePlaylistRequest)(nil),  // 0: proto.CreatePlaylistRequest
	(*CreatePlaylistResponse)(nil), // 1: proto.CreatePlaylistResponse
}
var file_playlist_proto_depIdxs = []int32{
	0, // 0: proto.PlaylistService.CreatePlaylist:input_type -> proto.CreatePlaylistRequest
	1, // 1: proto.PlaylistService.CreatePlaylist:output_type -> proto.CreatePlaylistResponse
	1, // [1:2] is the sub-list for method output_type
	0, // [0:1] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_playlist_proto_init() }
func file_playlist_proto_init() {
	if File_playlist_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_playlist_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreatePlaylistRequest); i {
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
		file_playlist_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*CreatePlaylistResponse); i {
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
			RawDescriptor: file_playlist_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   1,
		},
		GoTypes:           file_playlist_proto_goTypes,
		DependencyIndexes: file_playlist_proto_depIdxs,
		MessageInfos:      file_playlist_proto_msgTypes,
	}.Build()
	File_playlist_proto = out.File
	file_playlist_proto_rawDesc = nil
	file_playlist_proto_goTypes = nil
	file_playlist_proto_depIdxs = nil
}
