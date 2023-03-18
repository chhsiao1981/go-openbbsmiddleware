package api

import (
	"github.com/Ptt-official-app/go-openbbsmiddleware/schema"
	"github.com/Ptt-official-app/go-openbbsmiddleware/types"
)

var (
	testFilename3            = "M.1608386280.A.BC9"
	testContentAll3          []byte
	testContent3             []byte
	testSignature3           []byte
	testComment3             []byte
	testFirstCommentsDBCS3   []byte
	testTheRestCommentsDBCS3 []byte
	testContent3Big5         [][]*types.Rune
	testContent3Utf8         [][]*types.Rune
	testSignature3Utf8       [][]*types.Rune

	testFirstComments3 []*schema.Comment

	testFullFirstComments3 []*schema.Comment
)

func initTest3() {
	testContentAll3, testContent3, testSignature3, testComment3, testFirstCommentsDBCS3, testTheRestCommentsDBCS3 = loadTest(testFilename3)

	testContent3Utf8 = [][]*types.Rune{
		{
			{
				Utf8:   "作者: SYSOP () 看板: WhoAmI",
				Big5:   []byte("\xa7@\xaa\xcc: SYSOP () \xac\xdd\xaaO: WhoAmI"),
				DBCS:   []byte("\xa7@\xaa\xcc: SYSOP () \xac\xdd\xaaO: WhoAmI"),
				Color0: types.DefaultColor,
				Color1: types.DefaultColor,
			},
		},
		{
			{
				Utf8:   "標題: [心得] 測試一下特殊字～",
				Big5:   []byte("\xbc\xd0\xc3D: [\xa4\xdf\xb1o] \xb4\xfa\xb8\xd5\xa4@\xa4U\xafS\xae\xed\xa6r\xa1\xe3"),
				DBCS:   []byte("\xbc\xd0\xc3D: [\xa4\xdf\xb1o] \xb4\xfa\xb8\xd5\xa4@\xa4U\xafS\xae\xed\xa6r\xa1\xe3"),
				Color0: types.DefaultColor,
				Color1: types.DefaultColor,
			},
		},
		{
			{
				Utf8:   "時間: Sat Dec 19 21:57:58 2020",
				Big5:   []byte("\xae\xc9\xb6\xa1: Sat Dec 19 21:57:58 2020"),
				DBCS:   []byte("\xae\xc9\xb6\xa1: Sat Dec 19 21:57:58 2020"),
				Color0: types.DefaultColor,
				Color1: types.DefaultColor,
			},
		},
		{},
		{
			{
				Utf8:   "※這樣子有綠色嗎？～",
				Color0: types.DefaultColor,
				Color1: types.DefaultColor,
				Big5:   []byte("\xa1\xb0\xb3o\xbc\xcb\xa4l\xa6\xb3\xba\xf1\xa6\xe2\xb6\xdc\xa1H\xa1\xe3"),
				DBCS:   []byte("\xa1\xb0\xb3o\xbc\xcb\xa4l\xa6\xb3\xba\xf1\xa6\xe2\xb6\xdc\xa1H\xa1\xe3"),
			},
		},
		{
			{
				Utf8:   "※ 發信站",
				Color0: types.DefaultColor,
				Color1: types.DefaultColor,
				Big5:   []byte("\xa1\xb0 \xb5o\xabH\xaf\xb8"),
				DBCS:   []byte("\xa1\xb0 \xb5o\xabH\xaf\xb8"),
			},
		},
	}
	testSignature3Utf8 = [][]*types.Rune{
		{},
		{
			{
				Utf8:   "--",
				Color0: types.DefaultColor,
				Color1: types.DefaultColor,
				Big5:   []byte("--"),
				DBCS:   []byte("--"),
			},
		},
		{
			{
				Utf8:   "※ 發信站: 批踢踢 docker(pttdocker.test), 來自: 172.22.0.1",
				Big5:   []byte("\xa1\xb0 \xb5o\xabH\xaf\xb8: \xa7\xe5\xbd\xf0\xbd\xf0 docker(pttdocker.test), \xa8\xd3\xa6\xdb: 172.22.0.1"),
				DBCS:   []byte("\xa1\xb0 \xb5o\xabH\xaf\xb8: \xa7\xe5\xbd\xf0\xbd\xf0 docker(pttdocker.test), \xa8\xd3\xa6\xdb: 172.22.0.1"),
				Color0: types.DefaultColor,
				Color1: types.DefaultColor,
			},
		},
	}
}
