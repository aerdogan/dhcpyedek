package main

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
)

// Tarih karşılaştırma fonksiyonu
func isOneDayBefore(fileTime time.Time) bool {
	// Bugünün tarihi UTC kullanarak al
	today := time.Now().UTC()
	// Bugünden bir gün önceki tarihi al
	oneDayBefore := today.Add(-24 * time.Hour)
	// Dosya tarihinin bir gün önceyle aynı olup olmadığını kontrol et
	return fileTime.Year() == oneDayBefore.Year() && fileTime.YearDay() == oneDayBefore.YearDay()
}

// Tarih karşılaştırma fonksiyonu
func isToday(fileTime time.Time) bool {
	today := time.Now().UTC()
	// Dosya tarihinin aynı olup olmadığını kontrol et
	return fileTime.Year() == today.Year() && fileTime.YearDay() == today.YearDay()
}

func zipFiles(sourceDir, targetZip string) error {
	// ZIP dosyasını oluştur
	zipFile, err := os.Create(targetZip)
	if err != nil {
		return err
	}
	defer zipFile.Close()

	// ZIP yazıcısını başlat
	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Klasördeki tüm dosyaları gez
	err = filepath.Walk(sourceDir, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Dosya ise ve bir gün önce değiştirilmişse
		if !info.IsDir() && isOneDayBefore(info.ModTime()) {

			// ZIP dosyasına eklemek için bir dosya aç
			relPath, err := filepath.Rel(sourceDir, filePath)
			if err != nil {
				return err
			}

			// ZIP dosyasına bir yazma başlat
			fileInZip, err := zipWriter.Create(relPath)
			if err != nil {
				return err
			}

			// Kaynak dosyayı aç
			file, err := os.Open(filePath)
			if err != nil {
				return err
			}
			defer file.Close()

			// Dosyayı ZIP dosyasına yaz
			_, err = io.Copy(fileInZip, file)
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func main() {
	sourceDir := `C:\Windows\System32\Dhcp` // Sıkıştırılacak dosyaların bulunduğu klasör
	targetDir := `\\Watchguardlog\tmp`      // Hedef ZIP dosyasının yolu
	currentDate := time.Now().Format("2006-01-02")

	// Hedef ZIP dosyasının ismini oluştur
	targetZip := filepath.Join(targetDir, "dhcp_log_"+currentDate+".zip")

	// Dosyaları sıkıştır
	err := zipFiles(sourceDir, targetZip)
	if err != nil {
		fmt.Println("Hata:", err)
	} else {
		fmt.Println("ZIP dosyası başarıyla oluşturuldu:", targetZip)
	}
}
