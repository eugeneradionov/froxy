package models

import (
	"io"

	"github.com/google/uuid"
)

type Stream struct {
	ID uuid.UUID `json:"id"`

	// map[chunk position]*Chunk
	Chunks map[uint]*Chunk `json:"chunks,omitempty"`
}

type Chunk struct {
	ID uuid.UUID `json:"id"`

	FilePath string `json:"file_path"`
}

type AddChunkReq struct {
	Storage  StorageProvider `json:"storage"`
	Compress CompressAlgo    `json:"compress"`
	FileName string          `json:"file_name"`
	FilePath string          `json:"file_path"`
	Position uint            `json:"position"`

	StreamID uuid.UUID `json:"-"`
	File     io.Reader `json:"-"`
}
