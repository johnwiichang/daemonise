package daemonise

import (
	"fmt"
	"os"
	"os/exec"
)

type Daemon struct {
	pidf string
}

//New 新建 pid 文件
func New(file ...string) *Daemon {
	if len(file) == 0 {
		return &Daemon{GetPidFileName()}
	}
	return &Daemon{file[0]}
}

//ReportFunc 报告函数
func ReportFunc(pid int, ps *os.ProcessState) {
	if ps != nil {
		fmt.Fprintf(os.Stderr, "service failed with code %d\r\n", ps.ExitCode())
	} else {
		fmt.Fprintf(os.Stdout, "service started with pid %d\r\n", pid)
	}
}

func (d *Daemon) check() {
	if len(d.pidf) == 0 {
		d.pidf = GetPidFileName()
	}
}

//Run 运行
func (d *Daemon) Run(f, value string, callback ...func(int, *os.ProcessState)) {
	args := append([]string{f, value}, os.Args[1:]...)
	cmd := exec.Command(os.Args[0], args...)
	if d.Kill(); cmd.Start() == nil {
		d.savePid(cmd.Process.Pid)
	}
	if len(callback) > 0 {
		callback[0](cmd.Process.Pid, cmd.ProcessState)
	}
	os.Exit(0)
}

//Kill 杀死进程
func (d *Daemon) Kill() {
	d.kill()
	d.Done()
}

//IsRunning 判断进程是否运行
func (d *Daemon) IsRunning() bool {
	if pid := d.getPid(); pid > 0 {
		_, e := os.FindProcess(pid)
		return e == nil
	}
	return false
}

//Done 结束守护
func (d *Daemon) Done() {
	if len(d.pidf) > 0 {
		os.Remove(d.pidf)
	}
}
