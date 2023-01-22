package trace // パッケージ名と同じなのはたまたま

import (
	"fmt"
	"io"
)

// Tracerはコード内での出来事尾を記録できるオブジェクトを表すインスタンス
type Tracer interface { // 大文字なので公開
	Trace(...interface{}) // メソッドの指定
}

func New(w io.Writer) Tracer { // 大文字なので公開
	return &tracer{out: w} // Newするだけで隠された下の実装実行する
}

// ここからinterfaceの実装
type tracer struct { // 小文字なので非公開
	out io.Writer
}

func (t *tracer) Trace(a ...interface{}) {
	t.out.Write([]byte(fmt.Sprint(a...)))
	t.out.Write([]byte("\n"))
}
