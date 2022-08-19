package util

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path"
)

type DownloadInfo struct {
	Url  string
	Path string
}

func ParallelDownload(tasks []DownloadInfo, client *http.Client) []error {
	ch := make(chan error)

	for _, t := range tasks {
		go func(task DownloadInfo, client *http.Client, ch chan error) {
			ch <- DownloadFile(task, client)
		}(t, client, ch)
	}

	errs := make([]error, 0)

	for i := 0; i < len(tasks); i++ {
		err := <-ch
		if err != nil {
			errs = append(errs, err)
		}
	}

	return errs
}

func DownloadFile(task DownloadInfo, client *http.Client) error {
	err := os.MkdirAll(path.Dir(task.Path), os.ModePerm)
	if err != nil {
		return err
	}
	file, err := os.Create(task.Path)
	if err != nil {
		return err
	}
	defer file.Close()

	req, err := http.NewRequest("GET", task.Url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", "TakanashiDownloader/1.0.0")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.New(resp.Status)
	}

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}
