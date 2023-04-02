package convert

import (
	"fmt"
	"testing"

	"kusionstack.io/kclvm-go/pkg/kcl"
)

func TestConvert(t *testing.T) {
	result, err := kcl.GetSchemaType("./testdata/main.k", "", "")
	if err != nil {
		t.Fatal(err)
	}
	goCode := GenGoCodeFromKclType(result)
	fmt.Println(goCode)
	t.Logf("go code: %s", goCode)
}
