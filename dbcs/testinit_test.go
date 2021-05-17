package dbcs

import (
	"github.com/Ptt-official-app/go-openbbsmiddleware/schema"
	"github.com/Ptt-official-app/go-openbbsmiddleware/types"
)

func setupTest() {
	initTest()
	types.SetIsTest("dbcs")
	schema.SetIsTest()
}

func teardownTest() {
	schema.UnsetIsTest()
	types.UnsetIsTest("dbcs")
}
