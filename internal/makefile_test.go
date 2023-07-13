package makefile

import (
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestReadMakefile(t*testing.T){

	t.Parallel()
	t.Run("Invalid file path",func(t *testing.T) {
		_,_,err:=ReadMakefile("invalid path")

		assert.NotNil(t,err,"should return no such file or directory but got a file")
	})
	
}