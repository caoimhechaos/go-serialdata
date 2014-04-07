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

// A writer for a serial structure of consecutive data.
type SerialDataWriter struct {
	io.Writer
	w io.Writer
}

// Create a new record writer around "w".
func NewSerialDataWriter(w io.Writer) *SerialDataWriter {
	return &SerialDataWriter{
		w: w,
	}
}

// Write the data designated as "data" to the wrapped writer.
// It will be written as a separate new record.
func (s *SerialDataWriter) Write(data []byte) (int, error) {
	var lendata []byte = make([]byte, 4)
	var header_len int
	var body_len int
	var err error

	binary.BigEndian.PutUint32(lendata, uint32(len(data)))

	header_len, err = s.w.Write(lendata)
	if err != nil {
		return header_len, err
	}

	body_len, err = s.w.Write(data)
	if err != nil {
		return header_len + body_len, err
	}

	if body_len < len(data) {
		return header_len + body_len, errors.New("Short write")
	}

	return header_len + body_len, nil
}

// Write the designated protobuf record to the underlying writer.
func (s *SerialDataWriter) WriteMessage(pb proto.Message) error {
	var b []byte
	var err error

	b, err = proto.Marshal(pb)
	if err != nil {
		return err
	}

	_, err = s.Write(b)
	return err
}
