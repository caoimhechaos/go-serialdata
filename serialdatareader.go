/**
 * (c) 2014, Caoimhe Chaos <caoimhechaos@protonmail.com>,
 *	     Ancient Solutions. All rights reserved.
 *
 * Redistribution and use in source  and binary forms, with or without
 * modification, are permitted  provided that the following conditions
 * are met:
 *
 * * Redistributions of  source code  must retain the  above copyright
 *   notice, this list of conditions and the following disclaimer.
 * * Redistributions in binary form must reproduce the above copyright
 *   notice, this  list of conditions and the  following disclaimer in
 *   the  documentation  and/or  other  materials  provided  with  the
 *   distribution.
 * * Neither  the  name  of  Ancient Solutions  nor  the  name  of its
 *   contributors may  be used to endorse or  promote products derived
 *   from this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS
 * "AS IS"  AND ANY EXPRESS  OR IMPLIED WARRANTIES  OF MERCHANTABILITY
 * AND FITNESS  FOR A PARTICULAR  PURPOSE ARE DISCLAIMED. IN  NO EVENT
 * SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR ANY DIRECT,
 * INDIRECT, INCIDENTAL, SPECIAL,  EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED  TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE,  DATA, OR PROFITS; OR BUSINESS INTERRUPTION)
 * HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT,
 * STRICT  LIABILITY,  OR  TORT  (INCLUDING NEGLIGENCE  OR  OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED
 * OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package serialdata

import (
	"encoding/binary"
	"errors"
	"io"

	"code.google.com/p/goprotobuf/proto"
)

// A reader for a serial structure of consecutive data.
type SerialDataReader struct {
	io.ReadSeeker
	r io.Reader
}

// Create a new record reader around "r".
func NewSerialDataReader(r io.Reader) *SerialDataReader {
	return &SerialDataReader{
		r: r,
	}
}

// An error returned if the buffer isn't large enough to hold the next record.
var ENOBUFS error = errors.New("No buffer space available")

// Read exactly the next records worth of bytes from the reader and
// return its contents.
func (s *SerialDataReader) ReadRecord() ([]byte, error) {
	var seek io.ReadSeeker
	var buf []byte
	var pos int64 = -1
	var lendata []byte = make([]byte, 4)
	var header_len int
	var body_len uint32
	var read_len int
	var err error
	var ok bool

	seek, ok = s.r.(io.ReadSeeker)
	if ok {
		pos, err = seek.Seek(0, 1)
		if err != nil {
			return buf, err
		}
	}

	header_len, err = s.r.Read(lendata)
	if err != nil {
		if ok {
			seek.Seek(pos, 0)
		}

		return buf, err
	}

	if header_len != 4 {
		if ok {
			seek.Seek(pos, 0)
		}

		return buf, errors.New("Short read for header")
	}

	body_len = binary.BigEndian.Uint32(lendata)
	buf = make([]byte, body_len)
	read_len, err = s.r.Read(buf)
	if err == nil && uint32(read_len) < body_len {
		err = errors.New("Short read for body")
	}
	if err != nil {
		if ok {
			seek.Seek(pos, 0)
		}
	}

	return buf, err
}

// Read the next record from the wrapped reader. If the reader supports
// seek and the read fails, the position will be rolled back.
func (s *SerialDataReader) Read(data []byte) (int, error) {
	var buf []byte
	var err error

	buf, err = s.ReadRecord()
	if err != nil {
		return 0, err
	}

	if len(buf) > cap(data) {
		return 0, ENOBUFS
	}

	copy(data, buf)

	return len(buf), nil
}

// Read the next record and interpret it as a protocol buffer message.
func (s *SerialDataReader) ReadMessage(pb proto.Message) error {
	var buf []byte
	var err error

	buf, err = s.ReadRecord()
	if err != nil {
		return err
	}

	return proto.Unmarshal(buf, pb)
}

// Support seeks if the underlying object has got seek support. Seeks must
// always occur to the beginning of a record. This is mostly useful for
// reading from indexes.
func (s *SerialDataReader) Seek(offset int64, whence int) (int64, error) {
	var seeker io.Seeker
	var ok bool

	seeker, ok = s.r.(io.Seeker)
	if ok {
		return seeker.Seek(offset, whence)
	}

	return -1, errors.New("Underlying object does not support seeking")
}
