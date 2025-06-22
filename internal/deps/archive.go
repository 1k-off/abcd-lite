package deps

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
)

func extractZip(zipFile *os.File) error {
	stat, err := zipFile.Stat()
	if err != nil {
		return err
	}

	reader, err := zip.NewReader(zipFile, stat.Size())
	if err != nil {
		return err
	}

	for _, file := range reader.File {
		if file.Name == "LICENSE" {
			continue
		}

		if file.FileInfo().IsDir() {
			continue
		}

		rc, err := file.Open()
		if err != nil {
			return err
		}
		defer rc.Close()

		outPath := filepath.Join(binDir, filepath.Base(file.Name))
		outFile, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, file.Mode())
		if err != nil {
			return err
		}
		defer outFile.Close()

		if _, err := io.Copy(outFile, rc); err != nil {
			return err
		}
	}

	return nil
}

func extractTarGz(tarGzFile *os.File) error {
	gzReader, err := gzip.NewReader(tarGzFile)
	if err != nil {
		return err
	}
	defer gzReader.Close()

	tarReader := tar.NewReader(gzReader)

	for {
		header, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		if header.Name == "LICENSE" {
			continue
		}

		if header.Typeflag == tar.TypeDir {
			continue
		}

		outPath := filepath.Join(binDir, filepath.Base(header.Name))
		outFile, err := os.OpenFile(outPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, os.FileMode(header.Mode))
		if err != nil {
			return err
		}
		defer outFile.Close()

		if _, err := io.Copy(outFile, tarReader); err != nil {
			return err
		}
	}

	return nil
}
