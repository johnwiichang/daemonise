package daemonise

import (
	"encoding/binary"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

//GetPidFileName Associate a pid file to Daemon instance, named after `.{program}_pid` with default
func GetPidFileName(file ...string) (name string) {
	if len(file) == 0 {
		file = []string{os.Args[0]}
	}
	dir := filepath.Dir(file[0])
	if dot := strings.LastIndex(dir, "."); dot != -1 && dot != 0 {
		dir = dir[:dot]
	}
	name = strings.Replace(file[0], dir, "", 1) + "_pid"
	if len(name) > 0 && name[0] == filepath.Separator {
		name = name[1:]
	}
	name = "." + name
	return
}

//kill Send os.Kill if force is false, otherwise will kill the process directly
func (d *Daemon) kill(force ...bool) error {
	pid := d.getPid()
	if pid < 1 {
		return errors.New("invalid file")
	}
	p, err := os.FindProcess(pid)
	if err == nil {
		if len(force) > 0 && force[0] {
			err = p.Kill()
		} else {
			err = p.Signal(os.Kill)
		}
	}
	return err
}

//getPid Get PID associates with pid file
func (d *Daemon) getPid() (pid int) {
	d.check()
	bin, err := ioutil.ReadFile(d.pidf)
	if err != nil {
		return -1
	}
	return int(binary.BigEndian.Uint32(bin))
}

//savePid Save PID info to pid file.
func (d *Daemon) savePid(pid int) (err error) {
	d.check()
	var f *os.File
	if f, err = os.Create(d.pidf); err != nil {
		return
	}
	defer f.Close()
	//the pid will trans into uint32 format and write to file in BigEndian sequence
	if err = binary.Write(f, binary.BigEndian, uint32(pid)); err == nil {
		//hide() method has 2 design:
		//- Windows will use syscall to set file attribute
		//- Others is just a empty method
		d.hide()
	}
	return
}
