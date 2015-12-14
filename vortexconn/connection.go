package vortexconn

import (
	"encoding/binary"
	"io"
	"net"

	"github.com/pavben/Vortex/aesstream"
	"github.com/pavben/Vortex/pubkeycrypto"
)

// Connection is an AES-encrypted TCP connection.
type Connection struct {
	tcpConn        net.Conn
	aesStream      *aesstream.AesStream
	theirPublicKey *pubkeycrypto.PublicKey
}

func (c *Connection) Write(b []byte) error {
	return writeByteChunkPlain(c.aesStream, b)
}

func (c *Connection) Read() ([]byte, error) {
	return readByteChunkPlain(c.aesStream)
}

func writeByteChunkPlain(writer io.Writer, b []byte) error {
	chunkLen := uint32(len(b))
	if int(chunkLen) != len(b) {
		panic("writeByteChunkPlain chunk length overflows 32 bits")
	}
	err := binary.Write(writer, binary.BigEndian, chunkLen)
	if err != nil {
		return err
	}
	_, err = writer.Write(b)
	return err
}

func readByteChunkPlain(reader io.Reader) ([]byte, error) {
	var chunkLen uint32
	err := binary.Read(reader, binary.BigEndian, &chunkLen)
	if err != nil {
		return nil, err
	}
	b := make([]byte, chunkLen)
	err = plaintextBufferedRead(reader, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func plaintextBufferedRead(reader io.Reader, buf []byte) error {
	bytesRead := 0
	for bytesRead < len(buf) {
		// read up to the remainder of the empty space in buf
		n, err := reader.Read(buf[bytesRead:])
		if err != nil {
			return err
		}
		// update bytesRead to the new total number of bytes read
		bytesRead += n
	}
	return nil
}
