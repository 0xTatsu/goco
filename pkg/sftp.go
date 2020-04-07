package pkg

import (
	"github.com/pkg/errors"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
)

type SFTP struct {
	sftpConfig  *ssh.ClientConfig
	sftpAddress string
}

func NewSFTP(sftpUser, sftpPass, sftpAddress string) *SFTP {
	// nolint: gosec
	return &SFTP{
		sftpConfig: &ssh.ClientConfig{
			User: sftpUser,
			Auth: []ssh.AuthMethod{
				ssh.Password(sftpPass),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		},
		sftpAddress: sftpAddress,
	}
}

func (r *SFTP) getSFTPClient() (*sftp.Client, error) {
	conn, err := ssh.Dial("tcp", r.sftpAddress, r.sftpConfig)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot create connection to %s@%s", r.sftpConfig.User, r.sftpAddress)
	}

	sftpClient, err := sftp.NewClient(conn)
	if err != nil {
		return nil, errors.Wrapf(err, "cannot create SFTP client for %s@%s", r.sftpConfig.User, r.sftpAddress)
	}

	return sftpClient, nil
}
