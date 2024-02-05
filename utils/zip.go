package utils

import (
	"archive/zip"
	"io"
	"os"
	"path/filepath"
)

func ExtractZip(zipReader *zip.Reader, dir string) error {
	extractFile := func(fileInfo *zip.File) error {
		fSrc, err := fileInfo.Open()
		if err != nil {
			return err
		}
		fDst, err := os.Create(filepath.Join(dir, fileInfo.Name))
		if err != nil {
			return err
		}
		io.Copy(fDst, fSrc)
		defer fDst.Close()
		return nil
	}
	for _, fileInfo := range zipReader.File {
		err := extractFile(fileInfo)
		if err != nil {
			return err
		}
	}
	return nil
}
