# Parse DBCS

在 c-pttbbs 裡. 使用的 encoding 是 Big5 + [ANSI escape code](https://en.wikipedia.org/wiki/ANSI_escape_code). 另外利用 Big5 除了 0~127 以外, 都是 double-byte 的特性, 可以達到一字雙色的效果.

由於
