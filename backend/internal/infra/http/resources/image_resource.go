package resources

import (
	"fmt"
	"github.com/DmytroKha/nix-chat/internal/infra/database"
)

type ImageDto struct {
	Id   int64  `json:"id"`
	Link string `json:"link"`
}

type ImagesResource struct {
	images []database.Image
}

func NewImageResource(images []database.Image) ImagesResource {
	return ImagesResource{images: images}
}

func (ir ImagesResource) Serialize() []ImageDto {
	result := make([]ImageDto, 0, len(ir.images))

	for _, i := range ir.images {
		var dto ImageDto
		result = append(result, dto.DatabaseToDto(i))
	}

	return result
}

func (d *ImageDto) DatabaseToDto(i database.Image) ImageDto {
	link := fmt.Sprintf("/static/%s", i.Name)
	d.Id = i.Id
	d.Link = link
	return *d
}
