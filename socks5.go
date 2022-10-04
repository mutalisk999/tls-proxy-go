package tls_proxy_go

import (
	"bytes"
	"errors"
	"fmt"
)

const (
	ProtocolVersion = byte(0x05)
)

func getMethodsDescription(m byte) string {
	if m == byte(0x0) {
		return "NO AUTHENTICATION REQUIRED"
	} else if m == byte(0x01) {
		return "GSSAPI"
	} else if m == byte(0x02) {
		return "USERNAME/PASSWORD"
	} else if m == byte(0xff) {
		return "NO ACCEPTABLE METHODS"
	} else if m >= byte(0x03) && m <= byte(0x7f) {
		return "IANA ASSIGNED"
	} else if m >= byte(0x80) && m <= byte(0xfe) {
		return "RESERVED FOR PRIVATE METHODS"
	} else {
		return ""
	}
}

func parseHandshakeBody(body []byte) (bool, error) {
	if len(body) <= 2 {
		return false, errors.New("invalid handshake body length")
	}
	if len(body) != int(body[1])+2 {
		return false, errors.New("invalid handshake body length")
	}
	if body[0] != ProtocolVersion {
		return false, errors.New("only support socks5 protocol")
	}
	// check if method 'NO AUTHENTICATION REQUIRED' supported
	if bytes.Index(body[2:], []byte{0x0}) != -1 {
		return true, nil
	} else {
		return false, errors.New("method 'NO AUTHENTICATION REQUIRED' not supported")
	}
}

func parseRequestBody(body []byte) (*[]interface{}, error) {
	if len(body) <= 4 {
		return nil, errors.New("invalid handshake body length")
	}
	if body[3] != byte(0x01) && body[3] != byte(0x03) && body[3] != byte(0x04) {
		return nil, errors.New("invalid field ATYP")
	}
	if body[0] != ProtocolVersion {
		return nil, errors.New("only support socks5 protocol")
	}
	if body[1] != byte(0x01) && body[1] != byte(0x02) && body[1] != byte(0x03) {
		return nil, errors.New("invalid field CMD")
	}
	if body[2] != byte(0x0) {
		return nil, errors.New("invalid field RSV")
	}
	if body[3] == byte(0x1) {
		// ip v4
		if len(body) != 10 {
			return nil, errors.New("invalid request body length")
		} else {
			ret := make([]interface{}, 0)
			ret = append(ret, body[1], body[3])
			ret = append(ret, fmt.Sprintf("%d.%d.%d.%d", body[4], body[5], body[6], body[7]))
			ret = append(ret, (uint16(body[8])<<8)+uint16(body[9]))
			return &ret, nil
		}
	} else if body[3] == byte(0x3) {
		// domain name
		if len(body) != 7+int(body[4]) {
			return nil, errors.New("invalid request body length")
		} else {
			ret := make([]interface{}, 0)
			ret = append(ret, body[1], body[3])
			ret = append(ret, string(body[5:5+int(body[4])]))
			ret = append(ret, (uint16(body[5+int(body[4])])<<8)+uint16(body[6+int(body[4])]))
			return &ret, nil
		}
	} else if body[3] == byte(0x4) {
		// ip v6
		//if len(body) != 22 {
		//	return nil, errors.New("invalid request body length")
		//} else {
		//	ret := make([]interface{}, 0)
		//	ret = append(ret, body[1], body[3])
		//	ret = append(ret, fmt.Sprintf("%04x:%04x:%04x:%04x:%04x:%04x:%04x:%04x",
		//		body[4:6], body[6:8], body[8:10], body[10:12],body[12:14], body[14:16], body[16:18], body[18:20]))
		//	ret = append(ret, (uint16(body[20])<<8)+uint16(body[21]))
		//	return &ret, nil
		//}
		return nil, errors.New("ip v6 not supported")
	} else {
		return nil, nil
	}
}
