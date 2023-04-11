package archive

import "os"

type ByModTime []os.DirEntry

func (b ByModTime) Len() int {
	return len(b)
}

func (b ByModTime) Less(i, j int) bool {
	iInfo, _ := b[i].Info()
	jInfo, _ := b[j].Info()
	return iInfo.ModTime().Before(jInfo.ModTime())
}

func (b ByModTime) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
