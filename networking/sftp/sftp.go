package sftp

import (
	"bytes"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"time"
)

type Files struct {
	Name    string
	Content []byte
}

func RunSFTPClient(sshClient *ssh.Client) (*sftp.Client, error) {
	// connect to sftp server
	sftpClient, err := sftp.NewClient(sshClient)
	if err != nil {
		return nil, err
	}

	return sftpClient, nil
}

// Find takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
func Find(slice []string, val string) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

// FindTime takes a slice and looks for an element in it. If found it will
// return it's key, otherwise it will return -1 and a bool of false.
func FindTime(slice []time.Time, val time.Time) (int, bool) {
	for i, item := range slice {
		if item == val {
			return i, true
		}
	}
	return -1, false
}

func FindFiles(slice []Files, name string, content []byte) (int, bool) {
	for i, item := range slice {
		if item.Name == name {
			res := bytes.Compare(item.Content, content)
			if res == 0 {
				return i, true
			}
			return i, false
		}
	}
	return -1, false
}

func FindProcessedFiles(slice []*byte, name string, content string) (int, bool) {
	//for i, item := range slice {
	//	if item.Name == name {
	//		if item.Content == content {
	//			return i, true
	//		}
	//		return i, false
	//	}
	//}
	return -1, false
}
