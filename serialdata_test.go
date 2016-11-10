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
	"bytes"
	"testing"

	"github.com/golang/protobuf/proto"
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

// Serialize a protobuf message and try to read it again.
func TestSerializeAndReadMessage(t *testing.T) {
	var buf *bytes.Buffer = new(bytes.Buffer)
	var writer *SerialDataWriter = NewSerialDataWriter(buf)
	var reader *SerialDataReader
	var err error

	var data MessageForTest

	data.TestMessage = proto.String("Test data")
	err = writer.WriteMessage(&data)
	if err != nil {
		t.Error("Cannot serialize message: ", err)
	}

	data.TestMessage = proto.String("Toast Data")
	err = writer.WriteMessage(&data)
	if err != nil {
		t.Error("Cannot serialize message: ", err)
	}

	reader = NewSerialDataReader(bytes.NewReader(buf.Bytes()))
	data.Reset()

	err = reader.ReadMessage(&data)
	if err != nil {
		t.Error("Unable to re-read the message: ", err)
	}

	if data.GetTestMessage() != "Test data" {
		t.Errorf("Expected: Test data, got: %s", data.GetTestMessage())
	}
	data.Reset()

	err = reader.ReadMessage(&data)
	if err != nil {
		t.Error("Unable to re-read the message: ", err)
	}

	if data.GetTestMessage() != "Toast Data" {
		t.Errorf("Expected: Toast Data, got: %s", data.GetTestMessage())
	}
}

// Write a bunch of records to a memory buffer and read them back.
// Essentially, desperately shouting "Hello" into the void 10'000 times.
func BenchmarkRecordWriterAndReader(b *testing.B) {
	var buf *bytes.Buffer = new(bytes.Buffer)
	var writer *SerialDataWriter = NewSerialDataWriter(buf)
	var reader *SerialDataReader
	var rbuf []byte
	var err error
	var i, l int

	b.StartTimer()

	for i = 0; i < b.N; i++ {
		l, err = writer.Write([]byte("Hello"))
		if err != nil {
			b.Error("Error writing record: ", err)
		}

		if l != 9 {
			b.Error("Write length mismatched (expected 9, got ", l, ")")
		}
	}

	reader = NewSerialDataReader(bytes.NewReader(buf.Bytes()))
	for i = 0; i < b.N; i++ {
		rbuf, err = reader.ReadRecord()
		if err != nil {
			b.Error("Error reading record: ", err)
		}

		if len(rbuf) != 5 {
			b.Error("Read length mismatched (expected 5, got ", len(rbuf), ")")
		}

		if string(rbuf) != "Hello" {
			b.Error("Unexpected data: got ", string(rbuf), " (", rbuf,
				"), expected Hello")
		}
	}

	b.StopTimer()
	b.ReportAllocs()
}
