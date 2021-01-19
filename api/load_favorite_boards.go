package api

import (
	"github.com/Ptt-official-app/go-openbbsmiddleware/apitypes"
	"github.com/Ptt-official-app/go-openbbsmiddleware/schema"
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	pttbbsfav "github.com/Ptt-official-app/go-pttbbs/ptt/fav"
	"github.com/gin-gonic/gin"
)

const LOAD_FAVORITE_BOARDS_R = "/user/:user_id/favorites"

type LoadFavoriteBoardsParams struct {
	LevelIdx  schema.LevelIdx `json:"level_idx,omitempty" form:"level_idx,omitempty" url:"level_idx,omitempty"`
	StartIdx  string          `json:"start_idx,omitempty" form:"start_idx,omitempty" url:"start_idx,omitempty"`
	Ascending bool            `json:"asc,omitempty"  form:"asc,omitempty" url:"asc,omitempty"`
	Max       int             `json:"limit,omitempty" form:"limit,omitempty" url:"limit,omitempty"`
}

type LoadFavoriteBoardsPath struct {
	UserID bbs.UUserID `json:"user_id"`
}

type LoadFavoriteBoardsResult struct {
	List    []*apitypes.BoardSummary `json:"list"`
	NextIdx string                   `json:"next_idx"`
}

func NewLoadFavoriteBoardsParams() *LoadFavoriteBoardsParams {
	return &LoadFavoriteBoardsParams{
		Ascending: DEFAULT_ASCENDING,
		Max:       DEFAULT_MAX_LIST,
	}
}

func LoadFavoriteBoardsWrapper(c *gin.Context) {
	params := NewLoadFavoriteBoardsParams()
	path := &LoadFavoriteBoardsPath{}
	LoginRequiredPathQuery(LoadFavoriteBoards, params, path, c)
}

func LoadFavoriteBoards(remoteAddr string, userID bbs.UUserID, params interface{}, path interface{}, c *gin.Context) (result interface{}, statusCode int, err error) {

	theParams, ok := params.(*LoadFavoriteBoardsParams)
	if !ok {
		return nil, 400, ErrInvalidParams
	}

	thePath, ok := path.(*LoadFavoriteBoardsPath)
	if !ok {
		return nil, 400, ErrInvalidPath
	}

	if userID != thePath.UserID {
		return nil, 403, ErrInvalidUser
	}

	userFavorites_db, nextIdx, statusCode, err := tryGetUserFavorites(thePath.UserID, theParams.LevelIdx, theParams.StartIdx, theParams.Ascending, theParams.Max, c)
	if err != nil {
		return nil, statusCode, err
	}

	boardSummaryMap_db, userBoardInfoMap, statusCode, err := tryGetBoardSummaryMapFromUserFavorites(thePath.UserID, userFavorites_db, c)
	if err != nil {
		return nil, statusCode, err
	}

	r := NewLoadFavoriteBoardsResult(userFavorites_db, boardSummaryMap_db, nextIdx)

	err = checkBoardInfo(userID, userBoardInfoMap, r.List)
	if err != nil {
		return nil, 500, err
	}

	return r, 200, nil
}

func NewLoadFavoriteBoardsResult(userFavorites_db []*schema.UserFavorites, boardSummaryMap_db map[int]*schema.BoardSummary, nextIdx string) (result *LoadFavoriteBoardsResult) {

	theList := make([]*apitypes.BoardSummary, len(userFavorites_db))
	for idx, each := range userFavorites_db {
		if each.TheType != pttbbsfav.FAVT_BOARD {
			theList[idx] = apitypes.NewBoardSummaryFromUserFavorites(each, nil)
			continue
		}

		boardSummary_db := boardSummaryMap_db[each.TheID]
		theList[idx] = apitypes.NewBoardSummaryFromUserFavorites(each, boardSummary_db)
	}

	result = &LoadFavoriteBoardsResult{
		List:    theList,
		NextIdx: nextIdx,
	}

	return result
}
