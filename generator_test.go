package qbg

import (
	"io/ioutil"
	"testing"

	"github.com/favclip/genbase"
)

func TestGeneratorParsePackageDir(t *testing.T) {

	p := &genbase.Parser{SkipSemanticsCheck: true}
	pInfo, err := p.ParsePackageDir("./misc/fixture/a")
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if pInfo.Name() != "a" {
		t.Log("package name is not a, actual", pInfo.Name())
		t.Fail()
	}
	if len(pInfo.Files) != 2 {
		t.Log("package files length is not 2, actual", len(pInfo.Files))
		t.Fail()
	}
	if pInfo.Dir != "./misc/fixture/a" {
		t.Log("package dir is not ./misc/fixture/a, actual", pInfo.Dir)
		t.Fail()
	}
}

func TestGeneratorParsePackageFiles(t *testing.T) {
	p := &genbase.Parser{SkipSemanticsCheck: true}
	pInfo, err := p.ParsePackageFiles([]string{"./misc/fixture/a/model.go"})
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	if pInfo.Name() != "a" {
		t.Log("package name is not a, actual", pInfo.Name())
		t.Fail()
	}
	if len(pInfo.Files) != 1 {
		t.Log("package files length is not 1, actual", len(pInfo.Files))
		t.Fail()
	}
	if pInfo.Dir != "." {
		t.Log("package dir is not ., actual", pInfo.Dir)
		t.Fail()
	}
}

func TestGeneratorGenerate(t *testing.T) {

	testCase := []string{"a", "b"}
	for _, postFix := range testCase {
		p := &genbase.Parser{SkipSemanticsCheck: true}
		pInfo, err := p.ParsePackageDir("./misc/fixture/" + postFix)
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		args := []string{"-type", "Sample", "-output", "misc/fixture/" + postFix + "/model_query.go", "misc/fixture/" + postFix}
		typeNames := []string{"Sample"}
		if postFix == "g" {
			args = []string{"-type", "Sample,Inner", "-output", "misc/fixture/" + postFix + "/model_query.go", "misc/fixture/" + postFix}
			typeNames = []string{"Sample", "Inner"}
		}

		bu, err := Parse(pInfo, pInfo.CollectTypeInfos(typeNames))
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		src, err := bu.Emit(&args)
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		expected, err := ioutil.ReadFile("./misc/fixture/" + postFix + "/model_query.go")
		if err != nil {
			t.Log(err)
			t.Fail()
		}
		if string(src) != string(expected) {
			t.Log("not emit expected code in "+postFix+", actual\n", string(src))
			t.Fail()
		}
	}
}

func TestCollectTaggedTypes(t *testing.T) {
	p := &genbase.Parser{SkipSemanticsCheck: true}
	pInfo, err := p.ParsePackageFiles([]string{"./misc/fixture/c/model.go"})
	if err != nil {
		t.Log(err)
		t.Fail()
	}

	names := pInfo.CollectTaggedTypeInfos("+qbg")
	if len(names) != 1 {
		t.Log("names length is not 1, actual", len(names))
		t.Fail()
	}
	if names[0].Name() != "Sample" {
		t.Log("name[0] is not \"Sample\", actual", names[0].Name())
		t.Fail()
	}
}
