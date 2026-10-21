package ipar

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"strings"

	"github.com/chigopher/pathlib"
)

const (
	InternetArchiveMetadataEndpointStr = "https://archive.org/metadata"
)

var (
	InternetArchiveMetadataEndpoint *url.URL
)

func init() {
	var err error
	InternetArchiveMetadataEndpoint, err = url.Parse(InternetArchiveMetadataEndpointStr)
	if err != nil {
		panic(err)
	}
}

type InternetArchiveRepository interface {
	GetItem(identifier string) (*InternetArchiveItem, error)
	DownloadItemFiles(item *InternetArchiveItem, itemArchivalDir *pathlib.Path) error
}

type InternetArchiveItemIdentifier = string

type InternetArchiveItem struct {
	Directory string `json:"dir"`
	Files     []struct {
		CRC32 string
		MD5   string
		Name  string
		SHA1  string
		Size  string
	} `json:"files"`
	FileCnt     uint                        `json:"files_count"`
	ItemSize    uint64                      `json:"item_size"`
	LastUpdated uint64                      `json:"item_last_updated"`
	Metadata    InternetArchiveItemMetadata `json:"metadata"`
	Server      string                      `json:"server"`
}

type InternetArchiveItemMetadata struct {
	Identifier InternetArchiveItemIdentifier `json:"identifier"`
}

type InternetArchiveRepositoryHttpAPI struct{}

func (api *InternetArchiveRepositoryHttpAPI) GetItem(identifier InternetArchiveItemIdentifier) (*InternetArchiveItem, error) {
	if strings.TrimSpace(identifier) == "" {
		return nil, fmt.Errorf("item identifier is blank: '%s'", identifier)
	}
	var result InternetArchiveItem

	itemApiUrl := InternetArchiveMetadataEndpoint.JoinPath(identifier)

	res, err := http.Get(itemApiUrl.String())
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	statusCodeFamily := res.StatusCode / 100
	if statusCodeFamily != 2 {
		return nil, fmt.Errorf("GET request for %s resulted in non-200 status code %d", itemApiUrl, res.StatusCode)
	}

	itemBody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(itemBody, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (api *InternetArchiveRepositoryHttpAPI) DownloadItemFiles(item *InternetArchiveItem, itemArchivalDir *pathlib.Path) error {
	ok, err := itemArchivalDir.Exists()
	if err != nil {
		return err
	}
	if !ok {
		if err := itemArchivalDir.Mkdir(); err != nil {
			return err
		}
	}

	for _, file := range item.Files {
		filePath := itemArchivalDir.Join(file.Name)
		// Skip the download if the file already exists and its SHA1 fingerprint
		// matches
		exists, err := filePath.Exists()
		if err != nil {
			return err
		}
		if exists {
			hasher := sha1.New()
			fp, err := filePath.Open()
			if err != nil {
				return err
			}
			if _, err := io.Copy(hasher, fp); err != nil {
				return err
			}
			fileSHA1Hash := hex.EncodeToString(hasher.Sum(nil))
			if fileSHA1Hash == file.SHA1 {
				slog.Debug("skipped download, file already exists and contents match", "filename", file.Name)
				continue
			}
		}

		fileUrl, err := url.Parse(item.Server)
		if err != nil {
			return err
		}
		fileUrl = fileUrl.JoinPath(item.Directory).JoinPath(file.Name)
		err = api.downloadFile(fileUrl, filePath)
		if err != nil {
			return err
		}
	}

	slog.Debug("finished downloading", "id", item.Metadata.Identifier)

	return nil
}

func (api *InternetArchiveRepositoryHttpAPI) downloadFile(fileUrl *url.URL, targetPath *pathlib.Path) error {
	if fileUrl.Scheme == "" {
		fileUrl.Scheme = "https"
	}
	res, err := http.Get(fileUrl.String())
	if err != nil {
		return err
	}
	defer res.Body.Close()

	statusCodeFamily := res.StatusCode / 100
	if statusCodeFamily != 2 {
		return fmt.Errorf("GET request for %s resulted in non-200 status code %d", fileUrl.String(), res.StatusCode)
	}

	err = targetPath.WriteReader(res.Body)
	if err != nil {
		return err
	}

	slog.Debug("finished downloading", "url", fileUrl.String())
	return nil
}

var _ = InternetArchiveRepository((*InternetArchiveRepositoryHttpAPI)(nil))

func NewInternetArchiveRepositoryHttpAPI() *InternetArchiveRepositoryHttpAPI {
	return &InternetArchiveRepositoryHttpAPI{}
}
