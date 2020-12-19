package mysql

import (
	"fmt"
	"os/exec"
)

// Backup can be used to generate database backups
type Backup struct {
	Exec string
}

// NewBackup instanciate a backup interface
func NewBackup() *Backup {
	return &Backup{
		Exec: "mysqldump",
	}
}

// Run runs a backup and store it as the filename
func (m *Backup) Run(filename string) error {
	cmd := exec.Command(
		m.Exec,
		"--all-databases",
		"--lock-all-tables",
		"--host=127.0.0.1",
		fmt.Sprintf(`--result-file=%s`, filename),
	)
	return cmd.Run()
}
