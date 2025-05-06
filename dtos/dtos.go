// dtos/movie_dto.go
package dtos

type CreateMovieDto struct {
	Title       string `json:"title" validate:"required,min=1"`
	Description string `json:"description" validate:"required,min=1"`
	Director    string `json:"director" validate:"required,min=1"`
	ReleaseYear int    `json:"release_year" validate:"required,gte=1888,lte=2100"`
	Genre       string `json:"genre" validate:"required,min=1"`
}

type UpdateMovieDto struct {
	Title       *string `json:"title,omitempty"`
	Description *string `json:"description,omitempty"`
	Director    *string `json:"director,omitempty"`
	ReleaseYear *int    `json:"release_year,omitempty"`
	Genre       *string `json:"genre,omitempty"`
}
