package slack15

import "fmt"

// TODO: this should be really exported by log15
const errorKey = "LOG15_ERROR"

// ctxReader extracts key-value pairs from log15.Record.Ctx.
// Errors are handled transparently by setting LOG15_ERROR key
// with value describing the issue.
// TODO: on error key-value pair is lost, this could be improved.
type ctxReader struct {
	ctx []interface{}

	key   string
	value interface{}
	n     int // field counter for error reporting
}

func newCtxReader(ctx []interface{}) *ctxReader {
	return &ctxReader{ctx: ctx}
}

// Next process next key-value pair
func (r *ctxReader) Next() bool {
	if len(r.ctx) < 2 {
		// This shouldnever happen as log15 is normalising length of Record.Ctx
		return false
	}
	var ok bool
	r.key, ok = r.ctx[0].(string)
	r.value = r.ctx[1]
	if !ok {
		r.key = errorKey
		r.value = fmt.Errorf("ctx[%d] is not a string key (value=%+v)", r.n, r.value)
	}
	r.ctx = r.ctx[2:]
	r.n += 2
	return true
}

func (r *ctxReader) Key() string {
	return r.key
}

func (r *ctxReader) Value() interface{} {
	return r.value
}
