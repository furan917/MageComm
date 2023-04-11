package archive

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"magecomm/logger"
	"os"
	"os/exec"
	"strings"
)

const (
	ZipArchive      = "zip"
	TarArchive      = "tar"
	GzipArchive     = "gzip"
	XzArchive       = "xz"
	RarArchive      = "rar"
	Bzip2Archive    = "bzip2"
	SevenZipArchive = "7z"
)

var SupportedCatArchives = []string{".zip", ".tar", ".gz", ".bz2", ".rar", ".7z", ".xz"}

func CatFileFromDeploy(filePath string) error {
	deployPath, err := GetLatestDeploy()
	if err != nil {
		logger.Fatalf("Failed to get latest deploy, unable to cat file: %s", err)
	}

	err = CatFileFromArchive(deployPath, filePath)
	if err != nil {
		logger.Fatalf("Failed to cat file from deploy: %s", err)
	}

	return nil
}

func CatFileFromArchive(archivePath string, filePath string) error {
	if _, err := os.Stat(archivePath); os.IsNotExist(err) {
		return errors.New(fmt.Sprintf("archive '%s' does not exist or cannot be read", archivePath))
	}

	archiveType, err := archiveType(archivePath)
	if err != nil {
		return err
	}

	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	err = catFromArchive(file, archiveType, archivePath, filePath)
	if err != nil {
		return err
	}

	return nil
}

func catFromArchive(reader io.Reader, archiveType, archivePath, filePath string) error {
	switch archiveType {
	case ZipArchive:
		return catFromZipReader(archivePath, filePath)
	case TarArchive:
		return catFromTarReader(tar.NewReader(reader), filePath)
	case GzipArchive:
		gzipReader, err := gzip.NewReader(reader)
		if err != nil {
			return err
		}
		defer func(gzipReader *gzip.Reader) {
			err := gzipReader.Close()
			if err != nil {
				panic(err)
			}
		}(gzipReader)
		return catFromTarReader(tar.NewReader(gzipReader), filePath)
	case Bzip2Archive:
		return catFromTarReader(tar.NewReader(bzip2.NewReader(reader)), filePath)
	case SevenZipArchive:
		return catFromSevenZip(archivePath, filePath)
	case XzArchive:
		return catFromXz(archivePath, filePath)
	case RarArchive:
		return catFromRar(archivePath, filePath)
	default:
		return errors.New(fmt.Sprintf("archive type '%s' is not yet supported", archiveType))
	}
}

func archiveType(path string) (string, error) {
	file, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	header := make([]byte, 264)
	_, err = file.Read(header)
	if err != nil {
		return "", err
	}

	//Check the header as a file extension can lie, and opening the file will be too expensive
	switch {
	case bytes.Equal(header, []byte{0x50, 0x4B, 0x03, 0x04}): // zip
		return ZipArchive, nil
	case bytes.Equal(header, []byte{0x50, 0x4B, 0x05, 0x06}): // zip64
		return ZipArchive, nil
	case bytes.Equal(header, []byte{0x50, 0x4B, 0x07, 0x08}): // spanned zip
		return ZipArchive, nil
	case bytes.Equal(header, []byte{0x50, 0x4B, 0x53, 0x70}): // split zip
		return ZipArchive, nil
	case bytes.Equal(header, []byte{0x50, 0x4B, 0x01, 0x02}): // empty zip
		return ZipArchive, nil
	case len(header) >= 6 && bytes.Equal(header[:6], []byte{0x37, 0x7A, 0xBC, 0xAF, 0x27, 0x1C}): // 7zip
		return SevenZipArchive, nil
	case len(header) >= 6 && bytes.Equal(header[:6], []byte{0x52, 0x61, 0x72, 0x21, 0x1A, 0x07}): // 7zip split archive
		return SevenZipArchive, nil
	case len(header) >= 3 && bytes.Equal(header[:3], []byte{0x1F, 0x00, 0x00}): // zlib
		return GzipArchive, nil
	case len(header) >= 2 && bytes.Equal(header[:2], []byte{0xFD, 0x37}): // xz
		return XzArchive, nil
	case len(header) >= 6 && bytes.Equal(header[:6], []byte{0x42, 0x5A, 0x68, 0x39, 0x31, 0x41}): // bzip2
		return Bzip2Archive, nil
	case len(header) >= 4 && bytes.Equal(header[:4], []byte{0x50, 0x4B, 0x03, 0x04}): // docx, xlsx, pptx, jar
		return ZipArchive, nil
	case bytes.Equal(header, []byte{0x75, 0x73, 0x74, 0x61, 0x72, 0x00}): // tar
		return TarArchive, nil
	case len(header) >= 8 && (bytes.Equal(header[:8], []byte{0x75, 0x73, 0x74, 0x61, 0x72, 0x00, 0x00, 0x00}) || bytes.Equal(header[:8], []byte{0x75, 0x73, 0x74, 0x61, 0x72, 0x20, 0x00, 0x00})): // tar
		return TarArchive, nil
	case len(header) >= 4 && bytes.Equal(header[:4], []byte{0x75, 0x73, 0x74, 0x61}): // tar
		if len(header) >= 262 && (bytes.Equal(header[256:262], []byte(" ustar")) ||
			bytes.Equal(header[257:262], []byte("ustar")) ||
			bytes.Equal(header[257:262], []byte(" ustar")) ||
			bytes.Equal(header[257:262], []byte("ustar "))) {
			return TarArchive, nil
		}
		if len(header) >= 262 && (bytes.Equal(header[257:262], []byte("gnu  ")) ||
			bytes.Equal(header[257:262], []byte("GNUtar")) ||
			bytes.Equal(header[257:262], []byte("GNU tar"))) {
			return TarArchive, nil
		}
	case len(header) >= 4 && bytes.Equal(header[:4], []byte{0x75, 0x73, 0x74, 0x61}): // tar
		if len(header) >= 263 && bytes.Equal(header[257:263], []byte("GNU tar")) {
			return TarArchive, nil
		}
		if len(header) >= 263 && bytes.Equal(header[256:262], []byte(" ustar")) &&
			header[263] == byte(0) {
			return TarArchive, nil
		}
	case len(header) >= 262 && bytes.Equal(header[257:262], []byte("ustar")): //GNU TAR
		return TarArchive, nil
	case len(header) >= 262 && bytes.Equal(header[257:262], []byte(" \x00")): //GNU magic number Tar
		return TarArchive, nil
	case len(header) >= 3 && bytes.Equal(header[:3], []byte{0x1F, 0x9E, 0x01}): // compressed tar
		return TarArchive, nil
	case bytes.Equal(header, []byte{0x1F, 0x8B, 0x08}): // gzip
		return GzipArchive, nil
	case len(header) >= 2 && bytes.Equal(header[:2], []byte{0x1F, 0x8B}): // gzip
		return GzipArchive, nil
	case len(header) >= 2 && bytes.Equal(header[:2], []byte{0x1F, 0x9D}): // old-style gzip
		return GzipArchive, nil
	case len(header) >= 2 && bytes.Equal(header[:2], []byte{0x1F, 0xA0}): // old-style gzip
		return GzipArchive, nil
	case len(header) >= 8 && bytes.Equal(header[:8], []byte{0x52, 0x61, 0x72, 0x21, 0x1A, 0x07, 0x00, 0x20}): // rar
		return RarArchive, nil
	default:
		return "", fmt.Errorf("unknown archive type for file '%s'", path)
	}

	return "", nil
}

func catFromZipReader(archivePath string, filePath string) error {
	zipReader, err := zip.OpenReader(archivePath)
	if err != nil {
		return err
	}
	defer func(zipReader *zip.ReadCloser) {
		err := zipReader.Close()
		if err != nil {
			fmt.Printf("error closing zip reader: %s", err)
		}
	}(zipReader)

	var file *zip.File
	for _, f := range zipReader.File {
		if strings.HasSuffix(f.Name, filePath) {
			file = f
			break
		}
	}

	if file == nil {
		return errors.New(fmt.Sprintf("file '%s' not found in archive", filePath))
	}

	fileReader, err := file.Open()
	if err != nil {
		return err
	}
	defer func(fileReader io.ReadCloser) {
		err := fileReader.Close()
		if err != nil {
			fmt.Printf("error closing file: %s", err)
		}
	}(fileReader)

	_, err = io.Copy(os.Stdout, fileReader)
	if err != nil {
		return err
	}

	return nil
}

func catFromTarReader(tarReader *tar.Reader, filePath string) error {
	for {
		fileHeader, err := tarReader.Next()
		if err == io.EOF {
			return errors.New(fmt.Sprintf("file '%s' not found in archive", filePath))
		} else if err != nil {
			return err
		}
		if strings.HasSuffix(fileHeader.Name, filePath) {
			break
		}
	}

	_, err := io.Copy(os.Stdout, tarReader)
	if err != nil {
		return err
	}

	return nil
}

// The below functions are based on exec commands as go cant easily handle them natively, ig you do know a better way please let me know

func catFromSevenZip(archivePath string, filePath string) error {
	sevenZip, err := exec.LookPath("7z")
	if err != nil {
		return err
	}

	cmd := exec.Command(sevenZip, "e", "-so", archivePath, filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func catFromRar(archivePath string, filePath string) error {
	unrar, err := exec.LookPath("unrar")
	if err != nil {
		return err
	}

	cmd := exec.Command(unrar, "p", "-inul", archivePath, filePath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func catFromXz(archivePath string, filePath string) error {
	xz, err := exec.LookPath("xz")
	if err != nil {
		return err
	}

	cmd := exec.Command(xz, "--decompress", "--stdout", archivePath)
	cmd.Stderr = os.Stderr
	out, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	err = cmd.Start()
	if err != nil {
		return err
	}

	err = catFromTarReader(tar.NewReader(out), filePath)
	if err != nil {
		return err
	}

	return cmd.Wait()
}
