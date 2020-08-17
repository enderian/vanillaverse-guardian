package daemon

import (
	"net"
	"os"
)

func (d *Daemon) initializeSocket() {
	var err error
	d.recv, err = net.Listen("unix", d.Options.UnixSocket)
	if err != nil {
		d.Error("error while opening unix socket: %s", err)
		os.Exit(1)
	}

	defer d.recv.Close()
	d.Info("listening for commands on socket %s", d.Options.UnixSocket)

	for {
		conn, err := d.recv.Accept()
		if err != nil {
			d.Error("error while opening conn: %s", err)
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {

}
