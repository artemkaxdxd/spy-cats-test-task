package cat

import (
	"backend/config"
	request "backend/internal/controller/http/request/cat"
	response "backend/internal/controller/http/response/cat"
	"context"
	"fmt"
)

func (s service) GetCats(ctx context.Context, breed string) ([]response.Cat, config.ServiceCode, error) {
	cats, err := s.repo.GetCats(ctx, breed)
	return response.CatsToResponse(cats), config.DBErrToServiceCode(err), err
}

func (s service) GetCatByID(ctx context.Context, catID uint) (response.Cat, config.ServiceCode, error) {
	cat, err := s.repo.GetCatByID(ctx, catID)
	if cat.ID == 0 {
		return response.Cat{}, config.CodeNotFound, config.ErrCatNotFound
	}
	return response.CatToResponse(cat), config.DBErrToServiceCode(err), err
}

func (s service) CreateCat(ctx context.Context, body request.Cat) (response.Cat, config.ServiceCode, error) {
	valid, err := s.breedValidator.IsValid(ctx, body.Breed)
	if err != nil {
		s.l.Error("error validating breed with TheCatAPI", "breed", body.Breed, "err", err)
		return response.Cat{}, config.CodeExternalRequestFail, fmt.Errorf("could not validate breed: %v", err)
	}
	if !valid {
		return response.Cat{}, config.CodeBadRequest, fmt.Errorf("invalid breed: %s", body.Breed)
	}

	catEntity := body.ToEntity()

	createdCat, err := s.repo.CreateCat(ctx, catEntity)
	if err != nil {
		return response.Cat{}, config.DBErrToServiceCode(err), err
	}

	return response.CatToResponse(createdCat), config.CodeOK, nil
}

func (s service) UpdateCat(ctx context.Context, body request.UpdateCat, catID uint) (config.ServiceCode, error) {
	err := s.repo.UpdateCat(ctx, body.ToEntity(catID))
	return config.DBErrToServiceCode(err), err
}

func (s service) DeleteCat(ctx context.Context, catID uint) (config.ServiceCode, error) {
	err := s.repo.DeleteCat(ctx, catID)
	return config.DBErrToServiceCode(err), err
}
