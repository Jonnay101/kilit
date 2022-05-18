package command

import "os/exec"

type dataCreatorListener interface {
	OnDataCreation(data []byte)
}

func CreateData(l dataCreatorListener, name string, params ...string) {
	cmd := exec.Command(name, params...)
	data, _ := cmd.CombinedOutput()
	if data == nil {
		data = make([]byte, 0)
	}
	if l == nil {
		return
	}
	l.OnDataCreation(data)
}

type ProcessKiller struct {
	pids []string
}

func (k *ProcessKiller) OnPIDFound(pid string) {
	k.pids = append(k.pids, pid)
}

func (k *ProcessKiller) Kill() {
	for _, pid := range k.pids {
		killCommand := exec.Command("kill", "-9", pid)
		killCommand.Run()
	}
}
