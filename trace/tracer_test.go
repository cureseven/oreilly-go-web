package trace

import (
	"bytes"
	"testing"
)

func TestMew(t *testing.T) {
	var buf bytes.Buffer
	tracer := New(&buf)
	if tracer == nil {
		t.Error("Newからの戻り値がnilです")
	} else {
		tracer.Trace("こんにちは、traceパッケージ")
		if buf.String() != "こんにちは、traceパッケージ\n" {
			t.Errorf("'%s'という誤った文字列が出力されました", buf.String())
		}
	}
}

func testOff(t *testing.T) {
	var silentTracer Tracer = OFF()
	silentTracer.Trace("データ")
}
