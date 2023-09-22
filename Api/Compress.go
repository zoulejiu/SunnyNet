package Api

import "C"
import (
	"bytes"
	"compress/flate"
	"compress/gzip"
	"compress/zlib"
	"github.com/andybalholm/brotli"
	"github.com/qtgolang/SunnyNet/public"
	"github.com/qtgolang/SunnyNet/src/brotli-go/enc"
	"io"
	"io/ioutil"
)

// DeflateCompress Deflate压缩 (可能等同于zlib压缩)
func DeflateCompress(data uintptr, dataLen int) uintptr {
	input := public.CStringToBytes(data, dataLen)
	var o = &public.ZlibCompress{}
	f, _ := flate.NewWriter(o, flate.BestCompression)
	if a, b := f.Write(input); a == 0 || b != nil {
		return 0
	}
	if f.Flush() != nil {
		return 0
	}
	bx := o.Bytes()
	o.Close()
	bx = public.BytesCombine(public.IntToBytes(len(bx)), bx)
	return public.PointerPtr(string(bx))
}

// DeflateUnCompress Deflate解压缩 (可能等同于zlib解压缩)
func DeflateUnCompress(data uintptr, dataLen int) uintptr {
	bin := public.CStringToBytes(data, dataLen)
	if len(bin) < 1 {
		return 0
	}

	zr := flate.NewReader(ioutil.NopCloser(bytes.NewBuffer(bin)))
	defer func() { _ = zr.Close() }()
	bx, _ := io.ReadAll(zr)
	bx = public.BytesCombine(public.IntToBytes(len(bx)), bx)
	return public.PointerPtr(string(bx))
}

// ZlibUnCompress zlib解压缩
func ZlibUnCompress(data uintptr, dataLen int) uintptr {
	bin := public.CStringToBytes(data, dataLen)
	b := bytes.NewReader(bin)
	var out bytes.Buffer
	r, e := zlib.NewReader(b)
	if e != nil {
		return 0
	}
	_, _ = io.Copy(&out, r)
	bx := out.Bytes()
	bx = public.BytesCombine(public.IntToBytes(len(bx)), bx)
	return public.PointerPtr(string(bx))
}

// ZlibCompress zlib压缩
func ZlibCompress(data uintptr, dataLen int) uintptr {
	bin := public.CStringToBytes(data, dataLen)
	var buf bytes.Buffer
	compressor, err := zlib.NewWriterLevel(&buf, zlib.DefaultCompression)
	if err != nil {
		return 0
	}
	_, _ = compressor.Write(bin)
	_ = compressor.Close()
	out := buf.Bytes()
	out = public.BytesCombine(public.IntToBytes(len(out)), out)
	return public.PointerPtr(string(out))
}

// GzipCompress Gzip压缩
func GzipCompress(data uintptr, dataLen int) uintptr {
	bin := public.CStringToBytes(data, dataLen)
	var (
		buffer bytes.Buffer
		out    []byte
		err    error
	)
	writer := gzip.NewWriter(&buffer)
	_, err = writer.Write(bin)
	if err != nil {
		_ = writer.Close()
		return 0
	}
	err = writer.Close()
	if err != nil {
		return 0
	}
	out = buffer.Bytes()
	out = public.BytesCombine(public.IntToBytes(len(out)), out)
	return public.PointerPtr(string(out))
}

// BrCompress br压缩
func BrCompress(data uintptr, dataLen int) uintptr {
	bin := public.CStringToBytes(data, dataLen)
	compressedData, e := enc.CompressBuffer(nil, bin, make([]byte, 0))
	if e != nil {
		return 0
	}
	compressedData = public.BytesCombine(public.IntToBytes(len(compressedData)), compressedData)
	return public.PointerPtr(string(compressedData))
}

// BrUnCompress br解压缩
func BrUnCompress(data uintptr, dataLen int) uintptr {
	bin := public.CStringToBytes(data, dataLen)
	if len(bin) < 1 {
		return 0
	}
	b, _ := io.ReadAll(brotli.NewReader(ioutil.NopCloser(bytes.NewBuffer(bin))))
	b = public.BytesCombine(public.IntToBytes(len(b)), b)
	return public.PointerPtr(string(b))
}

// GzipUnCompress Gzip解压缩
func GzipUnCompress(data uintptr, dataLen int) uintptr {
	bin := public.CStringToBytes(data, dataLen)
	if len(bin) < 1 {
		return 0
	}
	gr, err := gzip.NewReader(ioutil.NopCloser(bytes.NewBuffer(bin)))
	if err != nil {
		return 0
	}
	b, _ := io.ReadAll(gr)
	b = public.BytesCombine(public.IntToBytes(len(b)), b)
	return public.PointerPtr(b)
}
