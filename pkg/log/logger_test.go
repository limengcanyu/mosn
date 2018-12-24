/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package log

import (
	"github.com/alipay/sofa-mosn/pkg/buffer"
	"io"
	"os"
	"runtime"
	"testing"
	"time"
)

func TestLogPrintDiscard(t *testing.T) {
	l, _ := NewLogger("/tmp/mosn_bench/benchmark.log", DEBUG)
	buf := buffer.GetIoBuffer(100)
	buf.WriteString("BenchmarkLog BenchmarkLog BenchmarkLog BenchmarkLog BenchmarkLog")
	l.Close()
	runtime.Gosched()
	// writeBufferChan is 1000
	// l.Debugf is discard, non block
	for i := 0; i < 1001; i++ {
		l.Debugf("BenchmarkLog BenchmarkLog BenchmarkLog BenchmarkLog BenchmarkLog %v", l)
	}
	lchan := make(chan struct{})
	go func() {
		// block
		l.Print(buf, false)
		lchan <- struct{}{}
	}()

	select {
	case <-lchan:
		t.Errorf("test Print diacard failed, should be block")
	case <-time.After(time.Second * 3):
	}
}

func TestLogPrintnull(t *testing.T) {
	logName := "/tmp/mosn_bench/printnull.log"
	os.Remove(logName)
	l, _ := NewLogger(logName, DEBUG)
	buf := buffer.GetIoBuffer(0)
	buf.WriteString("testlog")
	l.Print(buf, false)
	buf = buffer.GetIoBuffer(0)
	buf.WriteString("")
	l.Print(buf, false)
	l.Close()
	f, _ := os.Open(logName)
	b := make([]byte, 1024)
	n, _ := f.Read(b)
	f.Close()
	if n != len("testlog") {
		t.Errorf("Printnull error")
	}
	if string(b[:n]) != "testlog" {
		t.Errorf("Printnull error")
	}
}

func TestLogDefaultRollerTime(t *testing.T) {
	logName := "/tmp/mosn_bench/printdefaultroller.log"
	rollerName := logName + "." + time.Now().Format("2006-01-02")
	os.Remove(logName)
	os.Remove(rollerName)
	// 5s
	defaultRollerTime = 2
	logger, err := NewLogger(logName, RAW)
	if err != nil {
		t.Errorf("TestLogDefaultRoller failed %v", err)
	}
	logger.Print(buffer.NewIoBufferString("1111111"), false)
	time.Sleep(3 * time.Second)
	logger.Print(buffer.NewIoBufferString("2222222"), false)
	time.Sleep(1 * time.Second)
	logger.Close()

	b := make([]byte, 100)
	f, err := os.Open(logName)
	if err != nil {
		t.Errorf("TestLogDefaultRoller failed %v", err)
	}
	n, err := f.Read(b)
	f.Close()
	if err != io.EOF || n != 0 {
		t.Errorf("TestLogDefaultRoller failed %v", err)
	}

	f, err = os.Open(rollerName)
	if err != nil {
		t.Errorf("TestLogDefaultRoller failed %v", err)
	}
	n, err = f.Read(b)
	f.Close()
	if n == 0 || string(b[0:n]) != "11111112222222" {
		t.Errorf("TestLogDefaultRoller failed %v", string(b[0:n]))
	}

}

func BenchmarkLog(b *testing.B) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//InitDefaultLogger("", INFO)
	l, _ := NewLogger("/tmp/mosn_bench/benchmark.log", DEBUG)

	for n := 0; n < b.N; n++ {
		l.Debugf("BenchmarkLog BenchmarkLog BenchmarkLog BenchmarkLog BenchmarkLog %v", l)
	}
}

func BenchmarkLogParallel(b *testing.B) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	//InitDefaultLogger("", INFO)
	l, _ := NewLogger("/tmp/mosn_bench/benchmark.log", DEBUG)

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			l.Debugf("BenchmarkLog BenchmarkLog BenchmarkLog BenchmarkLog BenchmarkLog %v", l)
		}
	})
}

func BenchmarkLogTimeNow(b *testing.B) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	for n := 0; n < b.N; n++ {
		time.Now()
	}
}

func BenchmarkLogTimeFormat(b *testing.B) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	for n := 0; n < b.N; n++ {
		time.Now().Format("2006/01/02 15:04:05")
	}
}

func BenchmarkLogTime(b *testing.B) {
	runtime.GOMAXPROCS(runtime.NumCPU())
	for n := 0; n < b.N; n++ {
		logTime()
	}
}
