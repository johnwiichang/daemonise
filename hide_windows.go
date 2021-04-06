//+build windows

package daemonise

import (
	"syscall"
)

/*
	适用于 Windows 的文件隐藏方案：需要创建使用 HIDDEN 标识位的文件。
*/
func (d *Daemon) hide() error {
	filenameW, err := syscall.UTF16PtrFromString(d.pidf)
	if err == nil {
		err = syscall.SetFileAttributes(filenameW, syscall.FILE_ATTRIBUTE_HIDDEN)
	}
	if err != nil {
		return err
	}
	return nil
}
