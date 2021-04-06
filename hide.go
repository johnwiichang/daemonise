//+build !windows

package daemonise

func (d *Daemon) hide() error {
	return nil
}
