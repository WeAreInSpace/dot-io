package main

import (
	"bytes"
	"encoding/base64"
	"io"
	"os"

	dotio "github.com/WeAreInSpace/dot-io"
)

func main() {
	exampleFile, err := os.OpenFile("demo/local/example.exe", os.O_CREATE|os.O_RDONLY, 0777)
	if err != nil {
		panic(err)
	}

	outFile, err := os.OpenFile("demo/local/example.dio", os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		panic(err)
	}
	l := dotio.NewLocal(outFile)

	fileInBase64 := new(bytes.Buffer)
	fileEncoder := base64.NewEncoder(base64.RawStdEncoding, fileInBase64)
	io.Copy(fileEncoder, exampleFile)

	l.Ob.WriteStreamString(int64(fileInBase64.Len()), fileInBase64)
}
