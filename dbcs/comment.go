package dbcs

import (
	"bytes"
	"strings"

	"github.com/Ptt-official-app/go-openbbsmiddleware/schema"
	"github.com/Ptt-official-app/go-openbbsmiddleware/types"
	"github.com/Ptt-official-app/go-pttbbs/bbs"
)

var (
	MATCH_RECOMMEND_BYTES = []byte{ //\n推 abcd: test
		0x0a, 0x1b, 0x5b, 0x31, 0x3b, 0x33, 0x37, 0x6d,
		0xb1, 0xc0, 0x20,
	}
	LEN_MATCH_RECOMMEND_TYPES = len(MATCH_RECOMMEND_BYTES)

	MATCH_BOO_BYTES = []byte{ //\n噓 abcde: test2
		0x0a, 0x1b, 0x5b, 0x31, 0x3b, 0x33, 0x31, 0x6d,
		0xbc, 0x4e, 0x20,
	}

	MATCH_ARROW_BYTES = []byte{ //\n→ abcde: test3
		0x0a, 0x1b, 0x5b, 0x31, 0x3b, 0x33, 0x31, 0x6d,
		0xa1, 0xf7, 0x20,
	}

	MATCH_EDIT_BYTES = []byte{ //\n※ 編輯: abcdef (1.2.3.4 臺灣), 03/20/2021 23:58:08
		0x0a, 0xa1, 0xb0, 0x20, 0xbd, 0x73, 0xbf, 0xe8, 0x3a, 0x20,
	}

	MATCH_PREFIX_BYTES = []byte{ //\n※
		0x0a, 0xa1, 0xb0, 0x20,
	}

	MATCH_FORWARD_BYTES = []byte{ //※ silent2468:轉錄至看板 BLAZERS
		0x3a, 0xc2, 0xe0, 0xbf, 0xfd, 0xa6, 0xdc, 0xac, 0xdd, 0xaa, 0x4f, 0x20,
	}
)

func matchComment(content []byte) int {
	theIdx := -1

	//recommend
	idxRecommend := bytes.Index(content, MATCH_RECOMMEND_BYTES)
	if idxRecommend != -1 {
		theIdx = matchCommentIntegratedIdx(theIdx, idxRecommend)
	}
	//boo
	idxBoo := bytes.Index(content, MATCH_BOO_BYTES)
	if idxBoo != -1 {
		theIdx = matchCommentIntegratedIdx(theIdx, idxBoo)
	}
	//arrow
	idxArrow := bytes.Index(content, MATCH_ARROW_BYTES)
	if idxArrow != -1 {
		theIdx = matchCommentIntegratedIdx(theIdx, idxArrow)
	}
	//edit
	idxEdit := bytes.Index(content, MATCH_EDIT_BYTES)
	if idxEdit != -1 {
		theIdx = matchCommentIntegratedIdx(theIdx, idxEdit)
	}

	//idx-prefix
	idxPrefix := bytes.Index(content, MATCH_PREFIX_BYTES)
	if idxPrefix != -1 {
		prefixContent := content[idxPrefix:]
		idxNewLine := bytes.Index(prefixContent, []byte{'\n'})
		if idxNewLine != -1 {
			prefixContent = prefixContent[:idxNewLine]
		}

		//forward
		idxForward := bytes.Index(prefixContent, MATCH_FORWARD_BYTES)
		if idxForward != -1 {
			theIdx = matchCommentIntegratedIdx(theIdx, idxEdit)

		}
	}

	if theIdx == -1 {
		return theIdx
	}

	return theIdx + 1 //the prefix \n belongs to signature or content
}

func matchCommentIntegratedIdx(theIdx int, idx int) int {
	if theIdx == -1 {
		return idx
	}
	if theIdx < idx {
		return theIdx
	}
	return idx
}

//ParseComments
//
//1. 有可能 reply-edit-info (※  編輯:) 不在 commentsDBCS 裡
//   但是在 allCommentsDBCS 裡.
//   我們需要 allCommentsDBCS 來拿到 reply 的時間.
//
//2. commentsDBCS: firstComments 或是 allComments
//   allCommentsDBCS: allComments.
//   commentsDBCS 和 allCommentsDBCS 的 starting offset 是一樣的.
//
//3. 在 ParseComments 結束以後. 還無法決定 create-time 和 comment-id. 需要跟 db 整合來決定.
//
//4. 目前考慮: 推/噓/→/轉錄/編輯.
//
//5. 對於每一行的 parse. 如果不是上面 5 種. 都會被當成是 reply. reply 是 multi-line block.
//
//Implementation:
//1. 根據 '\n' estimate 有多少 comments
func ParseComments(
	bboardID bbs.BBoardID,
	articleID bbs.ArticleID,
	ownerID bbs.UUserID,
	commentsDBCS []byte,
	allCommentsDBCS []byte,
) (comments []*schema.Comment) {
	if len(commentsDBCS) == 0 {
		return nil
	}

	//1. estimate nComments
	nEstimatedComments := bytes.Count(commentsDBCS, []byte{'\n'})

	comments = make([]*schema.Comment, 0, nEstimatedComments)

	p_commentsDBCS := commentsDBCS
	p_allCommentsDBCS := allCommentsDBCS

	//2. for-loop '\n'
	for idxNewLine := bytes.Index(p_commentsDBCS, []byte{'\n'}); len(p_commentsDBCS) > 0 && idxNewLine != -1; idxNewLine = bytes.Index(p_commentsDBCS, []byte{'\n'}) {
		//2.1 parse comment
		commentDBCS := p_commentsDBCS[:idxNewLine]
		comment := parseComment(bboardID, articleID, commentDBCS)
		comments = append(comments, comment)

		p_commentsDBCS = p_commentsDBCS[idxNewLine:] // with '\n'
		p_allCommentsDBCS = p_allCommentsDBCS[idxNewLine:]

		//2.2 find next comment
		nextCommentIdx := matchComment(p_commentsDBCS)

		if nextCommentIdx == -1 { // no more comments
			//2.2.1 unable to find next comment, dealing with reply
			p_commentsDBCS = p_commentsDBCS[1:] //step forward '\n'
			p_allCommentsDBCS = p_allCommentsDBCS[1:]
			if len(p_commentsDBCS) > 0 { //整個都是 reply
				replyDBCS := p_commentsDBCS
				reply := parseReply(bboardID, articleID, replyDBCS, p_allCommentsDBCS)
				if reply != nil {
					comments = append(comments, reply)
				}

				p_allCommentsDBCS = p_allCommentsDBCS[len(p_commentsDBCS):]
				p_commentsDBCS = nil
			}
			break
		}

		if nextCommentIdx > 1 { // p_commentsDBCS[0] is '\n', get reply from p_commentsDBCS[1:]
			//2.3 with reply
			replyDBCS := p_commentsDBCS[1:nextCommentIdx]

			reply := parseReply(bboardID, articleID, replyDBCS, p_allCommentsDBCS[1:])
			if reply != nil {
				comments = append(comments, reply)
			}
		}

		p_commentsDBCS = p_commentsDBCS[nextCommentIdx:]
		p_allCommentsDBCS = p_allCommentsDBCS[nextCommentIdx:]
	}

	if len(p_commentsDBCS) > 0 {
		//XXX 2.4 shouldn't be here, assuming comment without reply.
		comment := parseComment(bboardID, articleID, p_commentsDBCS)
		comments = append(comments, comment)
	}

	return comments
}

//parseComment
//
//commentDBCS: excluding '\n'
func parseComment(
	bboardID bbs.BBoardID,
	articleID bbs.ArticleID,
	commentDBCS []byte) (

	comment *schema.Comment) {

	theType := parseCommentType(commentDBCS)

	var ownerID bbs.UUserID
	var content [][]*types.Rune
	var ip string
	var theDate string
	switch theType {
	case types.COMMENT_TYPE_RECOMMEND:
		ownerID, content, ip, theDate = parseCommentRecommendBooComment(commentDBCS)
	case types.COMMENT_TYPE_BOO:
		ownerID, content, ip, theDate = parseCommentRecommendBooComment(commentDBCS)
	case types.COMMENT_TYPE_COMMENT:
		ownerID, content, ip, theDate = parseCommentRecommendBooComment(commentDBCS)
	case types.COMMENT_TYPE_EDIT:
	}

	comment = &schema.Comment{
		TheType:   theType,
		BBoardID:  bboardID,
		ArticleID: articleID,
		Big5:      commentDBCS,
		MD5:       md5sum(commentDBCS),

		Owner:   ownerID,
		Content: content,
		IP:      ip,
		TheDate: theDate,
	}
	comment.CleanComment()

	return comment
}

func parseCommentType(p_commmentDBCS []byte) (theType types.CommentType) {
	if bytes.HasPrefix(p_commmentDBCS, MATCH_RECOMMEND_BYTES[1:]) {
		return types.COMMENT_TYPE_RECOMMEND
	} else if bytes.HasPrefix(p_commmentDBCS, MATCH_BOO_BYTES[1:]) {
		return types.COMMENT_TYPE_BOO
	} else if bytes.HasPrefix(p_commmentDBCS, MATCH_ARROW_BYTES[1:]) {
		return types.COMMENT_TYPE_COMMENT
	} else if bytes.HasPrefix(p_commmentDBCS, MATCH_EDIT_BYTES[1:]) {
		return types.COMMENT_TYPE_EDIT
	} else if bytes.HasPrefix(p_commmentDBCS, MATCH_PREFIX_BYTES[1:]) {
		if bytes.Contains(p_commmentDBCS, MATCH_FORWARD_BYTES) {
			return types.COMMENT_TYPE_FORWARD
		}
	}

	return types.COMMENT_TYPE_COMMENT
}

func parseCommentRecommendBooComment(commentDBCS []byte) (ownerID bbs.UUserID, content [][]*types.Rune, ip string, theDate string) {

	p_commentDBCS := commentDBCS[LEN_MATCH_RECOMMEND_TYPES-1:]
	ownerID, p_commentDBCS = parseCommentOwnerID(p_commentDBCS)
	contentDBCS, p_commentDBCS := parseCommentContent(p_commentDBCS)
	contentBig5 := dbcsToBig5(contentDBCS) //the last 11 chars are the dates
	content = big5ToUtf8(contentBig5)
	ip, theDate = parseCommentIPTheDate(p_commentDBCS)

	return ownerID, content, ip, theDate
}

func parseCommentOwnerID(p_commmentDBCS []byte) (ownerID bbs.UUserID, nextCommentDBCS []byte) {
	if len(p_commmentDBCS) == 0 {
		return "", nil
	}
	theIdx := bytes.Index(p_commmentDBCS, []byte{'\x1b'})
	if theIdx == -1 {
		return bbs.UUserID(""), nil
	}

	ownerID = bbs.UUserID(string(p_commmentDBCS[:theIdx]))
	if len(p_commmentDBCS) <= theIdx+8 {
		return ownerID, nil
	}
	nextCommentDBCS = p_commmentDBCS[theIdx+8:]

	return ownerID, nextCommentDBCS
}

func parseCommentContent(p_commmentDBCS []byte) (contentDBCS []byte, nextCommentDBCS []byte) {
	if len(p_commmentDBCS) == 0 {
		return nil, nil
	}

	idx := bytes.Index(p_commmentDBCS, []byte{'\x1b'})
	if idx == -1 {
		return p_commmentDBCS[1:], nil
	}

	contentDBCS, nextCommentDBCS = p_commmentDBCS[1:idx], p_commmentDBCS[idx:]
	if len(contentDBCS) > 0 && contentDBCS[0] == ' ' {
		contentDBCS = contentDBCS[1:]
	}
	if len(contentDBCS) == 0 {
		contentDBCS = nil
	}
	if len(nextCommentDBCS) == 0 {
		nextCommentDBCS = nil
	}
	idx = bytes.Index(nextCommentDBCS, []byte{'m'})
	if idx == -1 {
		nextCommentDBCS = nil
	}
	nextCommentDBCS = nextCommentDBCS[idx+1:]
	if len(nextCommentDBCS) == 0 {
		nextCommentDBCS = nil
	}

	return contentDBCS, nextCommentDBCS
}

//parseCommentIPTheDate
//
//Already separate the data by color.
//There are only ip/create-time information in p_commentDBCS.
func parseCommentIPTheDate(p_commentDBCS []byte) (ip string, theDate string) {
	if len(p_commentDBCS) == 0 {
		return "", ""
	}
	theIdx := bytes.Index(p_commentDBCS, []byte("\xb1\xc0")) //推
	if theIdx != -1 {                                        //old
		postfix := strings.TrimSpace(types.Big5ToUtf8(p_commentDBCS[theIdx+2:]))
		postfixList := strings.Split(postfix, " ")
		if len(postfixList) != 2 { //unable to parse. return createTime + 10-millisecond
			return "", ""
		}
		ip = postfixList[0]
		theDate = postfixList[1]

		return ip, theDate
	}

	//new: MM/DD HH:mm
	theDate = strings.TrimSpace(string(p_commentDBCS))

	return ip, theDate
}

func parseReply(
	bboardID bbs.BBoardID,
	articleID bbs.ArticleID,
	replyDBCS []byte,
	editDBCS []byte) (

	reply *schema.Comment) {

	//clean '\r\n'
	if len(replyDBCS) == 0 {
		return nil
	}
	if replyDBCS[len(replyDBCS)-1] == '\n' {
		replyDBCS = replyDBCS[:len(replyDBCS)-1]
	}
	if len(replyDBCS) == 0 {
		return nil
	}

	origReplyDBCS := replyDBCS
	if replyDBCS[len(replyDBCS)-1] == '\r' {
		replyDBCS = replyDBCS[:len(replyDBCS)-1]
	}
	if len(replyDBCS) == 0 {
		return nil
	}

	//md5sum
	replyMD5 := md5sum(replyDBCS)

	replyID := types.ToReplyID(commentID)

	replyBig5 := dbcsToBig5(replyDBCS)
	replyUtf8 := big5ToUtf8(replyBig5)

	editOwnerID, editNanoTS, editIP, editHost := parseReplyIPHost(editDBCS)

	reply = &schema.Comment{
		BBoardID:  bboardID,
		ArticleID: articleID,
		TheType:   types.COMMENT_TYPE_REPLY,
		Owner:     editOwnerID,
		Content:   replyUtf8,
		IP:        editIP,
		Host:      editHost,
		MD5:       replyMD5,
		Big5:      origReplyDBCS,

		EditNanoTS: editNanoTS,
	}

	reply.CleanReply()
	if len(reply.Content) == 0 {
		return nil
	}

	return reply
}

var (
	EDIT_PREFIX = []byte("\xa1\xb0 \xbds\xbf\xe8: ")
)

//parseReplyIPHost
//
//※ 編輯: abcde (1.2.3.4 臺灣)
func parseReplyIPHost(editDBCS []byte) (editOwnerID string, editNanoTS types.NanoTS, editIP string, editHost string) {

	p_editDBCS := editDBCS
	theIdx := bytes.Index(p_editDBCS, EDIT_PREFIX)
	if theIdx == -1 {
		return "", 0, "", ""
	}

	p_editDBCS = p_editDBCS[theIdx+len(EDIT_PREFIX):]

	theIdx = bytes.Index(p_editDBCS, []byte{'('})
	if theIdx == -1 {
		return "", 0, "", ""
	}

	editOwnerID = string(p_editDBCS[:theIdx-1])

	p_editDBCS = p_editDBCS[theIdx+1:]

	theIdx = bytes.Index(p_editDBCS, []byte{')'})
	if theIdx == -1 {
		return "", 0, "", ""
	}
	ipHost := types.Big5ToUtf8(p_editDBCS[:theIdx])

	ipHostList := strings.Split(ipHost, " ")
	if len(ipHostList) == 1 {
		editIP = ipHostList[0]
	} else {
		editIP = ipHostList[0]
		editHost = ipHostList[1]
	}

	p_editDBCS = p_editDBCS[theIdx:]

	theIdx = bytes.Index(p_editDBCS, []byte(", "))
	p_editDBCS = p_editDBCS[theIdx+2:]

	theTime, err := types.DateYearTimeStrToTime(string(p_editDBCS[:19]))
	if err != nil {
		return "", 0, "", ""
	}

	return editOwnerID, types.TimeToNanoTS(theTime), editIP, editHost

}
