package crd

import "testing"

func TestCrd(t *testing.T) {
	BuildDigest("./test_data/")
	PrintAst("./test_data/test.go", "./.crd/test.ast")
}
