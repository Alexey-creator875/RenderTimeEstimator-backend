package repository

import (
	"RenderTimeEstimator/internal/app/ds"
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strconv"

	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

func (r *Repository) GetPublishedRenderServerUnits() ([]ds.RenderServerUnit, error) {
	var renderServerUnits []ds.RenderServerUnit
	err := r.db.Where("status = ?", "published").Find(&renderServerUnits).Error

	if err != nil {
		return nil, err
	}

	return renderServerUnits, nil
}

func (r *Repository) GetPublishedRenderServerUnit(id int) (ds.RenderServerUnit, error) {
	renderServerUnit := ds.RenderServerUnit{}
	err := r.db.Where("id = ? AND status = ?", id, "published").First(&renderServerUnit).Error

	if err != nil {
		return ds.RenderServerUnit{}, err
	}

	return renderServerUnit, nil
}

func (r *Repository) GetPublishedRenderServerUnitsByRAM(min_ram int, max_ram int) ([]ds.RenderServerUnit, error) {
	var renderServerUnits []ds.RenderServerUnit
	err := r.db.Where("status = ? AND RAM BETWEEN ? AND ?", "published", min_ram, max_ram).Find(&renderServerUnits).Error

	if err != nil {
		return nil, err
	}

	return renderServerUnits, nil
}

func (r *Repository) GetDraftRenderServerUnit(creatorID uint) (*ds.RenderServerUnit, error) {
	draftRenderServerUnit := ds.RenderServerUnit{}
	err := r.db.Where("status = ? AND creator_id = ?", "draft", creatorID).Take(&draftRenderServerUnit).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return &ds.RenderServerUnit{}, err
	}

	return &draftRenderServerUnit, nil
}

func (r *Repository) GetNextPublishedRenderServerUnitTo(id int) (ds.RenderServerUnit, error) {
	var next ds.RenderServerUnit
	err := r.db.Where("id > ?", id).Order("id ASC").First(&next).Error

	if err != nil {
		return ds.RenderServerUnit{}, err
	}

	return next, nil
}

func (r *Repository) AddDraftRenderServerUnit(renderServerUnit *ds.RenderServerUnit) error {
	err := r.db.Create(renderServerUnit).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateRenderServerUnit(renderServerUnit ds.RenderServerUnit) error {
	err := r.db.Save(&renderServerUnit).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteRenderServerUnit(id int) error {
	query := "UPDATE render_server_units SET status = $1 WHERE id = $2"

	err := r.db.Exec(query, "deleted", id).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetLikesNumber(id int) (int, error) {
	query := "SELECT COUNT(*) FROM likes WHERE render_server_unit_id = $1"

	row := r.db.Raw(query, id).Row()

	var count int
	err := row.Scan(&count)

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *Repository) AddOrReplaceRenderServerUnitImage(RenderServerUnitID uint, header *multipart.FileHeader, ctx context.Context) error {
	filename := strconv.FormatUint(uint64(RenderServerUnitID), 10)

	file, err := header.Open()
	if err != nil {
		return fmt.Errorf("ошибка открытия файла: %w", err)
	}

	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)

	if err != nil {
		return fmt.Errorf("ошибка чтения файла: %w", err)
	}

	contentType := http.DetectContentType(buffer)

	_, err = file.Seek(0, 0)
	if err != nil {
		return fmt.Errorf("ошибка перемещения по файловому потоку: %w", err)
	}

	_, err = r.minio.PutObject(
		ctx,
		r.minio_bucket_name,
		filename,
		file,
		header.Size,
		minio.PutObjectOptions{
			ContentType: contentType,
		})

	if err != nil {
		return fmt.Errorf("не удалось добавить объект в хранилище minio: %w", err)
	}

	err = r.db.Model(&ds.RenderServerUnit{}).Where("id = ?", RenderServerUnitID).UpdateColumn("Image", "http://127.0.0.1:9000/"+r.minio_bucket_name+"/"+filename).Error
	if err != nil {
		// Если не удалось сохранить в БД, удаляем из MinIO
		r.minio.RemoveObject(
			ctx,
			r.minio_bucket_name,
			filename,
			minio.RemoveObjectOptions{})

		return fmt.Errorf("ошибка сохранения пути к изображению: %w", err)
	}

	return nil
}
