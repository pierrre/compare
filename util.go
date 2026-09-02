package compare

import (
	"github.com/pierrre/go-libs/bytesutil"
)

type appender interface {
	AppendTo(b []byte) []byte
}

func stringFromAppender[A appender](a A) string {
	bw := bytesWriterPool.Get()
	defer bytesWriterPool.Put(bw)
	*bw = a.AppendTo(*bw)
	return bw.String()
}

var bytesWriterPool = &bytesutil.WriterPool{
	MaxCap: -1,
}
