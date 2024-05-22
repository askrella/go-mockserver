package compression

import (
	"bytes"
	gzip2 "compress/gzip"
	"github.com/andybalholm/brotli"
	"io"
	"log"
)

func DecompressGzip(data []byte) (decompressed []byte) {
	reader, err := gzip2.NewReader(bytes.NewBuffer(data))
	if err != nil {
		log.Panicln("Cannot create new gzip reader:", err)
	}
	defer func(reader *gzip2.Reader) {
		err := reader.Close()
		if err != nil {
			log.Panicln("Error closing gzip reader: ", err)
		}
	}(reader)

	decompressed, err = io.ReadAll(reader)
	if err != nil {
		log.Panicln("Reading decompressed data from gzip reader caused an error: ", err)
	}

	return decompressed
}

func CompressBrotli(data []byte) (compressed []byte, contentEncoding string) {
	out := bytes.Buffer{}
	writer := brotli.NewWriterV2(&out, brotli.BestSpeed)
	_, err := writer.Write(data)
	if err != nil {
		log.Panicln("Error writing to brotli writer: ", err)
	}

	err = writer.Close()
	if err != nil {
		log.Panicln("Error closing brotli writer: ", err)
	}

	return out.Bytes(), "br"
}

func CompressGzip(data []byte) (compressed []byte, contentEncoding string) {
	var b bytes.Buffer
	gzw, err := gzip2.NewWriterLevel(&b, 9)
	if err != nil {
		log.Panicln("Error creating gzip writer:", err)
	}
	_, err = gzw.Write(data)
	if err != nil {
		log.Panicln("gzip write error:", err)
	}
	err = gzw.Flush()
	if err != nil {
		log.Panicln("gzip flush error:", err)
	}
	err = gzw.Close()
	if err != nil {
		log.Panicln("gzip close error:", err)
	}

	return b.Bytes(), "gzip"
}
