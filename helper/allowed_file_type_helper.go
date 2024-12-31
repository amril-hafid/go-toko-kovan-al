package helper

import (
	"fmt"
	"go-toko-kovan-al/config"
	"mime/multipart"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func AllowedFileType(file *multipart.FileHeader, conf config.Config) (bool, string) {
	// Mendapatkan ekstensi file
	// allowedExtensions := []string{".jpg", ".jpeg", ".png", ".PNG", ".pdf", ".PDF", ".doc", ".docx", ".xls", "xlsx", ".pptx", ".ppt"}
	allowedExtensions := strings.Split(conf.FileConf.FileType, "|")

	// allowedExtensions := strings.Split(os.Getenv("APP_FILE_TYPE"), "|")
	ext := filepath.Ext(file.Filename)

	// Memeriksa apakah ekstensi file diizinkan
	for _, allowedExt := range allowedExtensions {
		if strings.EqualFold(allowedExt, ext) {
			return true, allowedExt
		}
	}

	return false, ""
}

func StringWithoutSpaces(inputString string) string {
	var stringWithoutSpaces string
	for _, char := range inputString {
		if char != ' ' {
			stringWithoutSpaces += string(char)
		}
	}
	return stringWithoutSpaces
}

func RenameFile(file *multipart.FileHeader) string {
	now := time.Now()
	unixMilliseconds := now.UnixMilli()

	extension := filepath.Ext(file.Filename)

	nameFile := strings.TrimSpace(strings.TrimSuffix(file.Filename, extension))
	nameFile = strings.ReplaceAll(nameFile, " ", "_")

	fileString := fmt.Sprintf("%d-%s%s", unixMilliseconds, nameFile, extension)

	return fileString
}

func ChekFileSize(file *multipart.FileHeader, conf config.Config) bool {
	FileSizeFormatToMB := float64(file.Size) / 1024 / 1024
	var num float64
	num, _ = strconv.ParseFloat(conf.FileConf.FileMaxSizeTypeMB, 64)
	if num == 0 {
		return true
	}

	if FileSizeFormatToMB > num {
		return false
	}
	return true
}

// func IsAllowedFileTypeImage(file *multipart.FileHeader) (bool, string) {

// 	allowedExtensions := strings.Split(os.Getenv("APP_FILE_TYPE_IMAGE_PROFILE"), "|")
// 	ext := filepath.Ext(file.Filename)

// 	for _, allowedExt := range allowedExtensions {
// 		if strings.EqualFold(allowedExt, ext) {
// 			return true, allowedExt
// 		}
// 	}

// 	return false, ""
// }

// func DeleteFile(fileName string, conf config.Config) error {
// 	path := fmt.Sprintf("%s%s", conf.FileConf.FileDerektory, fileName)
// 	err := os.Remove(path)
// 	if err != nil {
// 		return err
// 	}
// 	return nil
// }

func DeleteFile(filePath string) error {
	// Cek apakah file ada
	path := fmt.Sprintf("assets/%s", filePath)

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return fmt.Errorf("file tidak ditemukan: %s", filePath)
		fmt.Println("File tidak ditemukan: %s", err)
	}

	// Hapus file
	err := os.Remove(filePath)
	if err != nil {
		return fmt.Errorf("gagal menghapus file: %w", err)
		fmt.Println("File tidak ditemukan 2: %s", err)

	}
	return nil
}

func FindRenameFile(conf config.Config, fileName string) (string, string, error) {
	oldFilePath := filepath.Join(conf.FileConf.FileDerektory, fileName)

	// Periksa apakah file ada
	if _, err := os.Stat(oldFilePath); os.IsNotExist(err) {
		return "", "", fmt.Errorf("file %s tidak ditemukan di direktori %s", fileName, conf.FileConf.FileDerektory)
	}

	// Buat nama baru
	newFileName := fmt.Sprintf("%d-%s", time.Now().UnixMilli(), strings.ReplaceAll(fileName, " ", "_"))
	newFilePath := filepath.Join(conf.FileConf.FileDerektory, newFileName)

	// Ganti nama file
	if err := os.Rename(oldFilePath, newFilePath); err != nil {
		return "", "", fmt.Errorf("gagal mengganti nama file: %v", err)
	}

	return newFileName, newFilePath, nil
}
