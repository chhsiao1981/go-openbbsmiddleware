package api

import (
	"strconv"

	"github.com/Ptt-official-app/go-openbbsmiddleware/apitypes"
	"github.com/Ptt-official-app/go-openbbsmiddleware/schema"
	"github.com/Ptt-official-app/go-pttbbs/bbs"
	pttbbsfav "github.com/Ptt-official-app/go-pttbbs/ptt/fav"
	"github.com/Ptt-official-app/go-pttbbs/ptttype"
	"github.com/gin-gonic/gin"
)

const ADD_FAVORITE_BOARD_R = "/user/:user_id/favorites/add"

type AddFavoriteBoardParams struct {
	LevelIdx schema.LevelIdx   `json:"level_idx,omitempty" form:"level_idx,omitempty" url:"level_idx,omitempty"`
	FBoardID apitypes.FBoardID `json:"bid" form:"bid" url:"bid"`
}

type AddFavoriteBoardPath struct {
	UserID bbs.UUserID `uri:"user_id"`
}

type AddFavoriteBoardResult struct {
	Summary *apitypes.BoardSummary `json:"brd"`
}

func AddFavoriteBoardWrapper(c *gin.Context) {
	params := &AddFavoriteBoardParams{}
	path := &LoadFavoriteBoardsPath{}
	LoginRequiredPathQuery(AddFavoriteBoard, params, path, c)
}

func AddFavoriteBoard(remoteAddr string, userID bbs.UUserID, params interface{}, path interface{}, c *gin.Context) (result interface{}, statusCode int, err error) {
	theParams, ok := params.(*AddFavoriteBoardParams)
	if !ok {
		return nil, 400, ErrInvalidParams
	}

	thePath, ok := path.(*AddFavoriteBoardPath)
	if !ok {
		return nil, 400, ErrInvalidPath
	}

	if userID != thePath.UserID {
		return nil, 403, ErrInvalidUser
	}

	boardID, err := toBoardID(theParams.FBoardID, remoteAddr, userID, c)
	if err != nil {
		return nil, 500, err
	}

	_, statusCode, err = isBoardValidUser(boardID, c)
	if err != nil {
		return nil, statusCode, err
	}

	bid, _, err := boardID.ToRaw()
	if err != nil {
		return nil, 500, err
	}

	userFavoritesMeta, userFavorites, err := getAllUserFavoritesFromDB(userID)
	if err != nil {
		return nil, 500, err
	}

	newUserFavorite := &schema.UserFavorites{
		UserID:   userID,
		FavIdx:   0,
		LevelIdx: theParams.LevelIdx,
		Idx:      0,
		TheType:  pttbbsfav.FAVT_BOARD,
		Attr:     pttbbsfav.FAVH_FAV,
		TheID:    int(bid),
	}

	userFavorites = insertNewUserFavorite(userFavorites, newUserFavorite)

	theFav, _ := schema.UserFavoritesToFav(&userFavoritesMeta.FolderMeta, userFavorites, 0)

	err := tryWriteFav(boardID, theFav, remoteAddr, userID, c)
	if err != nil {
		return nil, 500, err
	}

	startIdxStr := strconv.Itoa(newUserFavorite.Idx)

	newUserFavorites, _, statusCode, err := tryGetUserFavorites(userID, theParams.LevelIdx, startIdxStr, true, 1, c)
	if err != nil {
		return nil, statusCode, err
	}

	if len(newUserFavorites) != 1 {
		return nil, 500, ErrInvalidFav
	}

	boardSummaryMap_db, userBoardInfoMap, statusCode, err := tryGetBoardSummaryMapFromUserFavorites(thePath.UserID, newUserFavorites, c)
	if err != nil {
		return nil, statusCode, err
	}

	boardSummaries_db := make([]*schema.BoardSummary, 0, len(boardSummaryMap_db))
	for _, each := range boardSummaryMap_db {
		boardSummaries_db = append(boardSummaries_db, each)
	}

	userBoardInfoMap, err = checkUserReadBoard(userID, userBoardInfoMap, boardSummaries_db)
	if err != nil {
		return nil, 500, err
	}

	newUserFavorite = newUserFavorites[0]
	boardSummary_db, ok := boardSummaryMap_db[ptttype.Bid(newUserFavorite.TheID)]
	if !ok {
		return nil, 500, ErrInvalidFav
	}
	userBoardInfo, ok := userBoardInfoMap[boardSummary_db.BBoardID]
	if !ok {
		return nil, 500, ErrInvalidFav
	}

	summary := apitypes.NewBoardSummaryFromUserFavorites(newUserFavorites[0], boardSummary_db, userBoardInfo)

	ret := &AddFavoriteBoardResult{
		Summary: summary,
	}

	return ret, 200, nil
}
