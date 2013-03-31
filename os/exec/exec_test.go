package exec

import (
	"os/exec"
	"testing"

	"github.com/vsekhar/govtil/bytes"
	vtesting "github.com/vsekhar/govtil/testing"
)

func TestFileFromConn(t *testing.T) {
	in, out := vtesting.SelfConnection()
	defer in.Close()
	defer out.Close()

	fin, err := FileFromConn(in)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer fin.Close()
	fout, err := FileFromConn(out)
	if err != nil {
		t.Fatal(err.Error())
	}
	defer fout.Close()
	
	senddata := []byte("hello\n")

	go func() {
		in.Write(senddata)
	}()
	
	cmd := exec.Command("head", "-n1")
	cmd.Stdin = fout
	recvdata, err := cmd.Output()
	if !bytes.Equals(senddata, recvdata) {
		t.Error("failed, received data doesn't match sent data")
	}
}