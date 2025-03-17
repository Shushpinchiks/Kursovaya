package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

var token = os.Getenv("TOKEN")

func uploadToYandexDisk(filePath, token string) error {
	// Получаем URL для загрузки
	uploadURL, err := getUploadURL(filePath, token)
	if err != nil {
		return err
	}

	// Читаем файл
	fileData, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("ошибка чтения файла: %v", err)
	}

	// Загружаем файл
	req, err := http.NewRequest("PUT", uploadURL, bytes.NewReader(fileData))
	if err != nil {
		return fmt.Errorf("ошибка создания запроса: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		body, _ := ioutil.ReadAll(resp.Body)
		return fmt.Errorf("ошибка загрузки файла: %s", string(body))
	}

	fmt.Println("Файл успешно загружен на Яндекс.Диск")
	return nil
}

func getUploadURL(filePath, token string) (string, error) {
	url := fmt.Sprintf("https://cloud-api.yandex.net/v1/disk/resources/upload?path=%s&overwrite=true", filePath)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("ошибка создания запроса: %v", err)
	}

	req.Header.Set("Authorization", "OAuth "+token)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ошибка выполнения запроса: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		return "", fmt.Errorf("ошибка получения URL для загрузки: %s", string(body))
	}

	var result map[string]interface{}
	body, _ := ioutil.ReadAll(resp.Body)
	if err := json.Unmarshal(body, &result); err != nil {
		return "", fmt.Errorf("ошибка разбора ответа: %v", err)
	}

	href, ok := result["href"].(string)
	if !ok {
		return "", fmt.Errorf("не удалось получить ссылку для загрузки")
	}

	return href, nil
}

func TimeBackup() {
	backupTime := "22:19"

	for {
		now := time.Now()
		currentTime := now.Format("15:04")

		if currentTime == backupTime {
			file, err := backupDatabase()
			if err != nil {
				fmt.Println("Ошибка при создании бэкапа:", err)
			} else {
				fmt.Println("Бэкап успешно создан")

				err := uploadToYandexDisk(file, token)
				if err != nil {
					fmt.Println("Ошибка загрузки на Яндекс.Диск:", err)
				} else {
					fmt.Println("Бэкап успешно загружен на Яндекс.Диск")
				}
			}
			time.Sleep(23 * time.Hour)
		} else {
			time.Sleep(1 * time.Second)
		}
	}
}
