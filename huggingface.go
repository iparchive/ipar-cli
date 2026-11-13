package ipar

import (
	"github.com/chigopher/pathlib"
	"github.com/gomlx/go-huggingface/hub"
)

func init() {
}

type HuggingfaceRepository interface {
	DownloadModel(modelID HuggingfaceModelIdentifier, modelArchivalDir *pathlib.Path) (*HuggingfaceModel, error)
}

type HuggingfaceModelIdentifier = string

type HuggingfaceModel struct {
	ModelID          HuggingfaceModelIdentifier
	ModelArchivalDir *pathlib.Path
}

type HuggingfaceRepositoryHttpAPI struct {
	HuggingfaceToken string
}

func (api *HuggingfaceRepositoryHttpAPI) DownloadModel(modelID HuggingfaceModelIdentifier, modelArchivalDir *pathlib.Path) (*HuggingfaceModel, error) {
	modelRepo := hub.New(modelID).WithAuth(api.HuggingfaceToken)
	if err := modelRepo.Save(modelArchivalDir.String(), false); err != nil {
		return nil, err
	}
	if err := modelRepo.DeleteCache(); err != nil {
		return nil, err
	}
	return &HuggingfaceModel{
		ModelID:          modelID,
		ModelArchivalDir: modelArchivalDir,
	}, nil
}

var _ = HuggingfaceRepository((*HuggingfaceRepositoryHttpAPI)(nil))

func NewHuggingfaceRepositoryHttpAPI(hfToken string) *HuggingfaceRepositoryHttpAPI {
	api := &HuggingfaceRepositoryHttpAPI{
		HuggingfaceToken: hfToken,
	}
	return api
}
