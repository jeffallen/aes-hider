package main

import (
	"crypto/aes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
)

func Test(t *testing.T) {
	ready := make(chan struct{})
	done := make(chan struct{})
	go func() {
		k := make([]byte, 16)

		// Worst key evar.
		for i := range k {
			k[i] = byte(i)
		}

		c, err := aes.NewCipher(k)
		if err != nil {
			t.Fatal(err)
		}

		dst := make([]byte, 256)
		src := make([]byte, 256)
		c.Encrypt(dst, src)
		close(ready)

		// Sleeps here, holding key material in harm's way.
		_ = <-done
	}()
	_ = <-ready

	// Hi, harm.
	pid := os.Getpid()
	cmd := exec.Command("sudo", "../aes-finder/bin/aes-finder", "-1", "-p", fmt.Sprintf("%v", pid))
	out, err := cmd.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}
	if strings.Contains(string(out), "000102030405060708090a0b0c0d0e0f") {
		t.Fatal("key was found")
	}
	t.Log("output: ", string(out))
	close(done)
}
