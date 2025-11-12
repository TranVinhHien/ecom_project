package services

import (
	"context"
	"mime/multipart"
	"sync"
	"time"

	assets_services "github.com/TranVinhHien/ecom_product_service/services/assets"
	services "github.com/TranVinhHien/ecom_product_service/services/entity"
)

func (s *service) RenderImage(ctx context.Context, id string) string {
	filePath := s.env.ImagePath + id
	return filePath
}
func saveMediaFiles(files []*multipart.FileHeader, dstDir string, userID string) ([]services.Media, error) {
	var wg sync.WaitGroup
	resultCh := make(chan services.Media, len(files))
	errCh := make(chan error, len(files))

	for _, f := range files {
		wg.Add(1)
		go func(file *multipart.FileHeader) {
			defer wg.Done()
			if file_name, err := assets_services.SaveFile(file, dstDir); err != nil {
				errCh <- err
			} else {
				resultCh <- services.Media{
					FileName:  *file_name,
					Size:      file.Size,
					FilePath:  dstDir,
					MediaType: file.Header.Get("Content-Type"),
					CreatedAt: time.Now(),
					CreateBy:  userID,
				}
			}
		}(f)
	}

	wg.Wait()
	close(resultCh)
	close(errCh)

	var medias []services.Media
	for m := range resultCh {
		medias = append(medias, m)
	}

	if len(errCh) > 0 {
		return medias, <-errCh
	}
	return medias, nil
}

func (s *service) UploadMultiMedia(ctx context.Context, user_id string, files []*multipart.FileHeader) (result []string, err *assets_services.ServiceError) {
	// upload file to local folder
	medias, errors := saveMediaFiles(files, s.env.ImagePath, user_id)
	if errors != nil {
		return nil, &assets_services.ServiceError{Code: 400, Err: errors}
	}
	// urls:=
	urls := make([]string, 0, len(medias))
	for _, media := range medias {
		// urls = append(urls, s.env.HTTPServerAddress+"/v1/media/stream/"+media.FileName)
		urls = append(urls, media.FileName)
	}
	// result, errors = assets_services.HideFields(urls, "urls")
	// if errors != nil {
	// 	return nil, &assets_services.ServiceError{Code: 400, Err: errors}
	// }
	return urls, nil
}

func (s *service) DeleteMultiImage(ctx context.Context, user_id string, image_files []string) (err *assets_services.ServiceError) {

	for _, image := range image_files {
		errors := assets_services.DeleteFile(s.env.ImagePath, image)
		if errors != nil {
			return &assets_services.ServiceError{Code: 400, Err: errors}
		}
	}
	return nil
}
