package dbcs

import (
	"bytes"
)

//splitArticleSignatureCommentsDBCS
//
//given content, try to split as article / bbs-signature / comments
//
//在 edit 之後, 整個 content 有可能會被擾亂.
//我們只能盡可能的根據得到的資訊來決定是否是屬於 signature 或是 comments.
//
//rules:
//1. 整個 content 只會被分成 3 種可能: article / bbs-signature / comments
//2. bbs-signature 辨識方式: 最後 1 個 ※  發信站:
//3. comments 辨識方式: 符合顏色的推/噓/→
//4. 如果有 bbs-signature, comments 一定是在 bbs-signature 之後.
//5. 一旦出現 comments, 之後的 content 都會是 comments 或 reply.
//
//implementation:
//1. 嘗試拿到 signature 的 idx.
//2. 如果有 signature idx: 嘗試將 signature 之後的分出 comments, return article / signature, comments
//3. 否則: 嘗試將 content 分出 comments, return article / nil, comments
func splitArticleSignatureCommentsDBCS(content []byte) (articleDBCS []byte, signatureDBCS []byte, comments []byte) {

	//get content with signature
	idxes := tryGetSimpleSignatureIdxes(content)

	//start with largest idx
	reverseList(idxes)

	for _, idx := range idxes {
		isValid := isValidSimpleSignaure(content, idx)
		if !isValid {
			continue
		}

		//comments are after signature
		splitCommentIdx := matchComment(content[idx:])
		if splitCommentIdx == -1 {
			return content[:idx], content[idx:], nil
		}

		commentIdx := idx + splitCommentIdx
		return content[:idx], content[idx:commentIdx], content[commentIdx:]
	}

	// get content with comments, but no signatures.
	idx := matchComment(content)
	if idx == -1 {
		return content, nil, nil
	}

	return content[:idx], nil, content[idx:]
}

func tryGetSimpleSignatureIdxes(content []byte) []int {
	idxes := make([]int, 0, 5) //pre-allocated possible 5 candidate idxes.

	matchInit := []byte{ //\n※  發信站: (no \n-- in forward (轉))
		0x0a, 0xa1, 0xb0, 0x20, 0xb5,
		0x6f, 0xab, 0x48, 0xaf, 0xb8, 0x3a, 0x20,
	}

	p_content := content
	nowIdx := 0
	for idx := bytes.Index(p_content, matchInit); len(p_content) > 0 && (idx != -1); idx = bytes.Index(p_content, matchInit) {
		p_content = p_content[idx+1:]
		nowIdx += idx + 1
		idxes = append(idxes, nowIdx) //the prefix \n belongs to content
	}

	return idxes
}

func isValidSimpleSignaure(content []byte, idx int) (isValid bool) {
	p_content := content[idx:]

	return isSimpleSignatureWithFrom(p_content)
}

func isSimpleSignatureWithFrom(p_content []byte) bool {
	//check 來自 in the 1st line.
	idxLineBreak := bytes.Index(p_content, []byte{'\n'})
	line := p_content
	if idxLineBreak != -1 {
		line = p_content[:idxLineBreak]
	}

	matchFrom := []byte{ //), 來自:
		0x29, 0x2c, 0x20, 0xa8, 0xd3, 0xa6, 0xdb, 0x3a, 0x20,
	}
	if bytes.Index(line, matchFrom) != -1 {
		return true
	}

	//check From in the 2nd line.
	p_content = p_content[(idxLineBreak + 1):]
	line = p_content
	idxLineBreak = bytes.Index(p_content, []byte{'\n'})
	if idxLineBreak != -1 {
		line = p_content[:idxLineBreak]
	}

	matchFrom = []byte{ //◆  From:
		0xa1, 0xbb, 0x20, 0x46, 0x72, 0x6f, 0x6d, 0x3a, 0x20,
	}
	if bytes.HasPrefix(line, matchFrom) {
		return true
	}

	//check 轉錄者 in the 2nd line.
	matchForward := []byte{ //※ 轉錄者:
		0xa1, 0xb0, 0x20, 0xc2, 0xe0, 0xbf, 0xfd, 0xaa, 0xcc, 0x3a, 0x20,
	}
	if bytes.HasPrefix(line, matchForward) {
		return true
	}

	return false
}
