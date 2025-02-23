// Code generated by protoc-gen-gogo. DO NOT EDIT.
// source: cysicmint/types/v1/indexer.proto

package types

import (
	fmt "fmt"
	_ "github.com/cosmos/gogoproto/gogoproto"
	proto "github.com/gogo/protobuf/proto"
	io "io"
	math "math"
	math_bits "math/bits"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.GoGoProtoPackageIsVersion3 // please upgrade the proto package

// TxResult is the value stored in eth tx indexer
type TxResult struct {
	// height of the blockchain
	Height int64 `protobuf:"varint,1,opt,name=height,proto3" json:"height,omitempty"`
	// tx_index of the cosmos transaction
	TxIndex uint32 `protobuf:"varint,2,opt,name=tx_index,json=txIndex,proto3" json:"tx_index,omitempty"`
	// msg_index in a batch transaction
	MsgIndex uint32 `protobuf:"varint,3,opt,name=msg_index,json=msgIndex,proto3" json:"msg_index,omitempty"`
	// eth_tx_index is the index in the list of valid eth tx in the block,
	// aka. the transaction list returned by eth_getBlock api.
	EthTxIndex int32 `protobuf:"varint,4,opt,name=eth_tx_index,json=ethTxIndex,proto3" json:"eth_tx_index,omitempty"`
	// failed is true if the eth transaction did not go succeed
	Failed bool `protobuf:"varint,5,opt,name=failed,proto3" json:"failed,omitempty"`
	// gas_used by the transaction. If it exceeds the block gas limit,
	// it's set to gas limit, which is what's actually deducted by ante handler.
	GasUsed uint64 `protobuf:"varint,6,opt,name=gas_used,json=gasUsed,proto3" json:"gas_used,omitempty"`
	// cumulative_gas_used specifies the cumulated amount of gas used for all
	// processed messages within the current batch transaction.
	CumulativeGasUsed uint64 `protobuf:"varint,7,opt,name=cumulative_gas_used,json=cumulativeGasUsed,proto3" json:"cumulative_gas_used,omitempty"`
}

func (m *TxResult) Reset()         { *m = TxResult{} }
func (m *TxResult) String() string { return proto.CompactTextString(m) }
func (*TxResult) ProtoMessage()    {}
func (*TxResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_3771ff10ca9acdb1, []int{0}
}
func (m *TxResult) XXX_Unmarshal(b []byte) error {
	return m.Unmarshal(b)
}
func (m *TxResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	if deterministic {
		return xxx_messageInfo_TxResult.Marshal(b, m, deterministic)
	} else {
		b = b[:cap(b)]
		n, err := m.MarshalToSizedBuffer(b)
		if err != nil {
			return nil, err
		}
		return b[:n], nil
	}
}
func (m *TxResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_TxResult.Merge(m, src)
}
func (m *TxResult) XXX_Size() int {
	return m.Size()
}
func (m *TxResult) XXX_DiscardUnknown() {
	xxx_messageInfo_TxResult.DiscardUnknown(m)
}

var xxx_messageInfo_TxResult proto.InternalMessageInfo

func init() {
	proto.RegisterType((*TxResult)(nil), "cysicmint.types.v1.TxResult")
}

func init() { proto.RegisterFile("cysicmint/types/v1/indexer.proto", fileDescriptor_3771ff10ca9acdb1) }

var fileDescriptor_3771ff10ca9acdb1 = []byte{
	// 304 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x44, 0x90, 0xb1, 0x4e, 0xc3, 0x30,
	0x14, 0x45, 0x63, 0xda, 0xa6, 0xc1, 0x82, 0x81, 0x80, 0xaa, 0x00, 0x52, 0xb0, 0x98, 0xc2, 0x40,
	0xa2, 0x8a, 0x8d, 0x11, 0x06, 0xc4, 0x6a, 0x95, 0x85, 0x25, 0x4a, 0x93, 0x87, 0x63, 0x91, 0xd4,
	0x55, 0xfd, 0x52, 0xd2, 0x3f, 0x60, 0xe4, 0x13, 0xf8, 0x1c, 0xc6, 0x8e, 0x8c, 0xa8, 0x15, 0xff,
	0x81, 0xea, 0x44, 0x61, 0x7b, 0x57, 0xe7, 0x3c, 0x5d, 0xe9, 0x52, 0x96, 0xae, 0xb4, 0x4c, 0x4b,
	0x39, 0xc3, 0x08, 0x57, 0x73, 0xd0, 0xd1, 0x72, 0x1c, 0xc9, 0x59, 0x06, 0x35, 0x2c, 0xc2, 0xf9,
	0x42, 0xa1, 0x72, 0xdd, 0xce, 0x08, 0x8d, 0x11, 0x2e, 0xc7, 0x67, 0x27, 0x42, 0x09, 0x65, 0x70,
	0xb4, 0xbb, 0x1a, 0xf3, 0xf2, 0x97, 0x50, 0x67, 0x52, 0x73, 0xd0, 0x55, 0x81, 0xee, 0x88, 0xda,
	0x39, 0x48, 0x91, 0xa3, 0x47, 0x18, 0x09, 0x7a, 0xbc, 0x4d, 0xee, 0x29, 0x75, 0xb0, 0x8e, 0x4d,
	0x85, 0xb7, 0xc7, 0x48, 0x70, 0xc8, 0x87, 0x58, 0x3f, 0xee, 0xa2, 0x7b, 0x4e, 0xf7, 0x4b, 0x2d,
	0x5a, 0xd6, 0x33, 0xcc, 0x29, 0xb5, 0x68, 0x20, 0xa3, 0x07, 0x80, 0x79, 0xdc, 0xfd, 0xf6, 0x19,
	0x09, 0x06, 0x9c, 0x02, 0xe6, 0x93, 0xf6, 0x7d, 0x44, 0xed, 0x97, 0x44, 0x16, 0x90, 0x79, 0x03,
	0x46, 0x02, 0x87, 0xb7, 0x69, 0xd7, 0x28, 0x12, 0x1d, 0x57, 0x1a, 0x32, 0xcf, 0x66, 0x24, 0xe8,
	0xf3, 0xa1, 0x48, 0xf4, 0x93, 0x86, 0xcc, 0x0d, 0xe9, 0x71, 0x5a, 0x95, 0x55, 0x91, 0xa0, 0x5c,
	0x42, 0xdc, 0x59, 0x43, 0x63, 0x1d, 0xfd, 0xa3, 0x87, 0xc6, 0xbf, 0xed, 0xbf, 0x7f, 0x5e, 0x58,
	0x77, 0xf7, 0x5f, 0x1b, 0x9f, 0xac, 0x37, 0x3e, 0xf9, 0xd9, 0xf8, 0xe4, 0x63, 0xeb, 0x5b, 0xeb,
	0xad, 0x6f, 0x7d, 0x6f, 0x7d, 0xeb, 0xf9, 0x4a, 0x48, 0xcc, 0xab, 0x69, 0x98, 0xaa, 0x32, 0x32,
	0xb3, 0x5d, 0x17, 0xc9, 0x54, 0xb7, 0xe7, 0x0c, 0xf0, 0x4d, 0x2d, 0x5e, 0x9b, 0x9d, 0xa7, 0xb6,
	0xd9, 0xec, 0xe6, 0x2f, 0x00, 0x00, 0xff, 0xff, 0xfe, 0xed, 0x7b, 0x07, 0x81, 0x01, 0x00, 0x00,
}

func (m *TxResult) Marshal() (dAtA []byte, err error) {
	size := m.Size()
	dAtA = make([]byte, size)
	n, err := m.MarshalToSizedBuffer(dAtA[:size])
	if err != nil {
		return nil, err
	}
	return dAtA[:n], nil
}

func (m *TxResult) MarshalTo(dAtA []byte) (int, error) {
	size := m.Size()
	return m.MarshalToSizedBuffer(dAtA[:size])
}

func (m *TxResult) MarshalToSizedBuffer(dAtA []byte) (int, error) {
	i := len(dAtA)
	_ = i
	var l int
	_ = l
	if m.CumulativeGasUsed != 0 {
		i = encodeVarintIndexer(dAtA, i, uint64(m.CumulativeGasUsed))
		i--
		dAtA[i] = 0x38
	}
	if m.GasUsed != 0 {
		i = encodeVarintIndexer(dAtA, i, uint64(m.GasUsed))
		i--
		dAtA[i] = 0x30
	}
	if m.Failed {
		i--
		if m.Failed {
			dAtA[i] = 1
		} else {
			dAtA[i] = 0
		}
		i--
		dAtA[i] = 0x28
	}
	if m.EthTxIndex != 0 {
		i = encodeVarintIndexer(dAtA, i, uint64(m.EthTxIndex))
		i--
		dAtA[i] = 0x20
	}
	if m.MsgIndex != 0 {
		i = encodeVarintIndexer(dAtA, i, uint64(m.MsgIndex))
		i--
		dAtA[i] = 0x18
	}
	if m.TxIndex != 0 {
		i = encodeVarintIndexer(dAtA, i, uint64(m.TxIndex))
		i--
		dAtA[i] = 0x10
	}
	if m.Height != 0 {
		i = encodeVarintIndexer(dAtA, i, uint64(m.Height))
		i--
		dAtA[i] = 0x8
	}
	return len(dAtA) - i, nil
}

func encodeVarintIndexer(dAtA []byte, offset int, v uint64) int {
	offset -= sovIndexer(v)
	base := offset
	for v >= 1<<7 {
		dAtA[offset] = uint8(v&0x7f | 0x80)
		v >>= 7
		offset++
	}
	dAtA[offset] = uint8(v)
	return base
}
func (m *TxResult) Size() (n int) {
	if m == nil {
		return 0
	}
	var l int
	_ = l
	if m.Height != 0 {
		n += 1 + sovIndexer(uint64(m.Height))
	}
	if m.TxIndex != 0 {
		n += 1 + sovIndexer(uint64(m.TxIndex))
	}
	if m.MsgIndex != 0 {
		n += 1 + sovIndexer(uint64(m.MsgIndex))
	}
	if m.EthTxIndex != 0 {
		n += 1 + sovIndexer(uint64(m.EthTxIndex))
	}
	if m.Failed {
		n += 2
	}
	if m.GasUsed != 0 {
		n += 1 + sovIndexer(uint64(m.GasUsed))
	}
	if m.CumulativeGasUsed != 0 {
		n += 1 + sovIndexer(uint64(m.CumulativeGasUsed))
	}
	return n
}

func sovIndexer(x uint64) (n int) {
	return (math_bits.Len64(x|1) + 6) / 7
}
func sozIndexer(x uint64) (n int) {
	return sovIndexer(uint64((x << 1) ^ uint64((int64(x) >> 63))))
}
func (m *TxResult) Unmarshal(dAtA []byte) error {
	l := len(dAtA)
	iNdEx := 0
	for iNdEx < l {
		preIndex := iNdEx
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return ErrIntOverflowIndexer
			}
			if iNdEx >= l {
				return io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= uint64(b&0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		fieldNum := int32(wire >> 3)
		wireType := int(wire & 0x7)
		if wireType == 4 {
			return fmt.Errorf("proto: TxResult: wiretype end group for non-group")
		}
		if fieldNum <= 0 {
			return fmt.Errorf("proto: TxResult: illegal tag %d (wire type %d)", fieldNum, wire)
		}
		switch fieldNum {
		case 1:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Height", wireType)
			}
			m.Height = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIndexer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.Height |= int64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 2:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field TxIndex", wireType)
			}
			m.TxIndex = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIndexer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.TxIndex |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 3:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field MsgIndex", wireType)
			}
			m.MsgIndex = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIndexer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.MsgIndex |= uint32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 4:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field EthTxIndex", wireType)
			}
			m.EthTxIndex = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIndexer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.EthTxIndex |= int32(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 5:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field Failed", wireType)
			}
			var v int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIndexer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				v |= int(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			m.Failed = bool(v != 0)
		case 6:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field GasUsed", wireType)
			}
			m.GasUsed = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIndexer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.GasUsed |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		case 7:
			if wireType != 0 {
				return fmt.Errorf("proto: wrong wireType = %d for field CumulativeGasUsed", wireType)
			}
			m.CumulativeGasUsed = 0
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return ErrIntOverflowIndexer
				}
				if iNdEx >= l {
					return io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				m.CumulativeGasUsed |= uint64(b&0x7F) << shift
				if b < 0x80 {
					break
				}
			}
		default:
			iNdEx = preIndex
			skippy, err := skipIndexer(dAtA[iNdEx:])
			if err != nil {
				return err
			}
			if (skippy < 0) || (iNdEx+skippy) < 0 {
				return ErrInvalidLengthIndexer
			}
			if (iNdEx + skippy) > l {
				return io.ErrUnexpectedEOF
			}
			iNdEx += skippy
		}
	}

	if iNdEx > l {
		return io.ErrUnexpectedEOF
	}
	return nil
}
func skipIndexer(dAtA []byte) (n int, err error) {
	l := len(dAtA)
	iNdEx := 0
	depth := 0
	for iNdEx < l {
		var wire uint64
		for shift := uint(0); ; shift += 7 {
			if shift >= 64 {
				return 0, ErrIntOverflowIndexer
			}
			if iNdEx >= l {
				return 0, io.ErrUnexpectedEOF
			}
			b := dAtA[iNdEx]
			iNdEx++
			wire |= (uint64(b) & 0x7F) << shift
			if b < 0x80 {
				break
			}
		}
		wireType := int(wire & 0x7)
		switch wireType {
		case 0:
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowIndexer
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				iNdEx++
				if dAtA[iNdEx-1] < 0x80 {
					break
				}
			}
		case 1:
			iNdEx += 8
		case 2:
			var length int
			for shift := uint(0); ; shift += 7 {
				if shift >= 64 {
					return 0, ErrIntOverflowIndexer
				}
				if iNdEx >= l {
					return 0, io.ErrUnexpectedEOF
				}
				b := dAtA[iNdEx]
				iNdEx++
				length |= (int(b) & 0x7F) << shift
				if b < 0x80 {
					break
				}
			}
			if length < 0 {
				return 0, ErrInvalidLengthIndexer
			}
			iNdEx += length
		case 3:
			depth++
		case 4:
			if depth == 0 {
				return 0, ErrUnexpectedEndOfGroupIndexer
			}
			depth--
		case 5:
			iNdEx += 4
		default:
			return 0, fmt.Errorf("proto: illegal wireType %d", wireType)
		}
		if iNdEx < 0 {
			return 0, ErrInvalidLengthIndexer
		}
		if depth == 0 {
			return iNdEx, nil
		}
	}
	return 0, io.ErrUnexpectedEOF
}

var (
	ErrInvalidLengthIndexer        = fmt.Errorf("proto: negative length found during unmarshaling")
	ErrIntOverflowIndexer          = fmt.Errorf("proto: integer overflow")
	ErrUnexpectedEndOfGroupIndexer = fmt.Errorf("proto: unexpected end of group")
)
