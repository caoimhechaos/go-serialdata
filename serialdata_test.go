/**
 * (c) 2014, Tonnerre Lombard <tonnerre@ancient-solutions.com>,
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
	"bytes"
	"testing"
)

// Write a bit of test data and test the io.Reader interface.
func TestSerializeAndReadBytes(t *testing.T) {
	var buf *bytes.Buffer = new(bytes.Buffer)
	var writer *SerialDataWriter = NewSerialDataWriter(buf)
	var reader *SerialDataReader
	var rbuf []byte
	var err error
	var l int

	l, err = writer.Write([]byte("Hello"))
	if err != nil {
		t.Error("Error writing record: ", err)
	}

	if l != 9 {
		t.Error("Write length mismatched (expected 9, got ", l, ")")
	}

	if buf.Len() != 9 {
		t.Error("Expected length to be 9, was ", buf.Len())
	}

	l, err = writer.Write([]byte("World"))
	if err != nil {
		t.Error("Error writing record: ", err)
	}

	if l != 9 {
		t.Error("Write length mismatched (expected 9, got ", l, ")")
	}

	if buf.Len() != 18 {
		t.Error("Expected length to be 18, was ", buf.Len())
	}

	reader = NewSerialDataReader(bytes.NewReader(buf.Bytes()))
	rbuf = make([]byte, 20)

	l, err = reader.Read(rbuf)
	if err != nil {
		t.Error("Error reading record: ", err)
	}

	if l != 5 {
		t.Error("Read length mismatched (expected 5, got ", l, ")")
	}

	if string(rbuf[0:l]) != "Hello" {
		t.Error("Unexpected data: got ", string(rbuf), " (", rbuf,
			"), expected Hello")
	}

	rbuf = make([]byte, 20)

	l, err = reader.Read(rbuf)
	if err != nil {
		t.Error("Error reading record: ", err)
	}

	if l != 5 {
		t.Error("Read length mismatched (expected 5, got ", l, ")")
	}

	if string(rbuf[0:l]) != "World" {
		t.Error("Unexpected data: got ", string(rbuf), " (", rbuf,
			"), expected World")
	}
}

// Write a bit of test data and read individual records off it.
func TestSerializeAndReadRecord(t *testing.T) {
	var buf *bytes.Buffer = new(bytes.Buffer)
	var writer *SerialDataWriter = NewSerialDataWriter(buf)
	var reader *SerialDataReader
	var rbuf []byte
	var err error
	var l int

	l, err = writer.Write([]byte("Hello"))
	if err != nil {
		t.Error("Error writing record: ", err)
	}

	if l != 9 {
		t.Error("Write length mismatched (expected 9, got ", l, ")")
	}

	if buf.Len() != 9 {
		t.Error("Expected length to be 9, was ", buf.Len())
	}

	l, err = writer.Write([]byte("World"))
	if err != nil {
		t.Error("Error writing record: ", err)
	}

	if l != 9 {
		t.Error("Write length mismatched (expected 9, got ", l, ")")
	}

	if buf.Len() != 18 {
		t.Error("Expected length to be 18, was ", buf.Len())
	}

	reader = NewSerialDataReader(bytes.NewReader(buf.Bytes()))

	rbuf, err = reader.ReadRecord()
	if err != nil {
		t.Error("Error reading record: ", err)
	}

	if len(rbuf) != 5 {
		t.Error("Read length mismatched (expected 5, got ", len(rbuf), ")")
	}

	if string(rbuf) != "Hello" {
		t.Error("Unexpected data: got ", string(rbuf), " (", rbuf,
			"), expected Hello")
	}

	rbuf, err = reader.ReadRecord()
	if err != nil {
		t.Error("Error reading record: ", err)
	}

	if len(rbuf) != 5 {
		t.Error("Read length mismatched (expected 5, got ", len(rbuf), ")")
	}

	if string(rbuf) != "World" {
		t.Error("Unexpected data: got ", string(rbuf), " (", rbuf,
			"), expected World")
	}
}
