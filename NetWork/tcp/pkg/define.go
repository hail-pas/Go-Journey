package pkg

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"unsafe"
)

type Header struct {
	packageLength uint64
	infoLength    uint64
	info          string
	length        uint64
	bodyLength    uint64
}

type Body struct {
	Version string `json:"version"`
	Topic   string `json:"topic"`
	Author  string `json:"author"`
	Content string `json:"content"`
}

type TcpPackage struct {
	RawData []byte
	Header  *Header
	Body    *Body
}

func (p TcpPackage) packageLength() uint64 {
	return ConvertByte2Int(p.RawData[:8])
}
func (p TcpPackage) headerInfoLength() uint64 {
	return ConvertByte2Int(p.RawData[8:16])
}

func (p TcpPackage) headerInfo() string {
	bs := p.RawData[16:p.headerInfoLength()]
	return *(*string)(unsafe.Pointer(&bs))
}

func (p TcpPackage) headerBodyLength() uint64 {
	return ConvertByte2Int(p.RawData[8+8+p.headerInfoLength()+8 : 8+8+p.headerInfoLength()+8+8])
}

func (p TcpPackage) bodyContent() *Body {
	bs := p.RawData[p.headerLength() : p.headerLength()+p.headerBodyLength()]
	var body = new(Body)
	err := json.Unmarshal(bs, &body)
	if err != nil {
		return body
	}
	return body
}

func (p TcpPackage) headerLength() uint64 {
	return ConvertByte2Int(p.RawData[8+8+p.headerInfoLength() : 8+8+p.headerInfoLength()+8])
}

func ConvertByte2Int(bs []byte) uint64 {
	var num uint64
	switch len(bs) {
	case 1:
		num = uint64(bs[0])
	case 2:
		num = uint64(binary.BigEndian.Uint16(bs))
	case 4:
		num = uint64(binary.BigEndian.Uint32(bs))
	case 8:
		num = binary.BigEndian.Uint64(bs)

	default:
		return 0
	}
	return num
}

func NewTcpPackage(rawData []byte) *TcpPackage {
	tcpPackage := &TcpPackage{
		RawData: rawData,
	}

	tcpPackage.Header = &Header{
		infoLength: tcpPackage.headerInfoLength(),
		info:       tcpPackage.headerInfo(),
		length:     tcpPackage.headerLength(),
		bodyLength: tcpPackage.headerBodyLength()}
	tcpPackage.Body = tcpPackage.bodyContent()
	return tcpPackage
}

func GenerateRawData(headerInfo string, body Body) (data []byte, err error) {
	headerInfoByte := []byte(headerInfo)
	headerInfoLength := make([]byte, 8)
	headerLength := make([]byte, 8)
	bodyLength := make([]byte, 8)
	pkgLength := make([]byte, 8)
	bodyByte, err := json.Marshal(body)
	if err != nil {
		fmt.Println("json marshal failed, err: ", err)
		return
	}
	binary.BigEndian.PutUint64(headerInfoLength, uint64(len(headerInfoByte)))
	binary.BigEndian.PutUint64(headerLength, uint64(8+8+len(headerInfoByte)+8+8))
	binary.BigEndian.PutUint64(bodyLength, uint64(len(bodyByte)))
	data = append(headerInfoLength, headerInfoByte...)
	data = append(data, headerLength...)
	data = append(data, bodyLength...)
	data = append(data, bodyByte...)
	binary.BigEndian.PutUint64(pkgLength, uint64(len(data)+8))
	data = append(pkgLength, data...)
	return
}
