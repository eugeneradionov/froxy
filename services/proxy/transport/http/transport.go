package http

import (
	"context"
	"net/http"

	"github.com/eugeneradionov/froxy/models"
	"github.com/eugeneradionov/froxy/pkg/logger"
	"github.com/eugeneradionov/xerrors"
	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"
)

type service interface {
	GetMaxFileSize() int64
	Upload(ctx context.Context, req *models.UploadReq) (*models.UploadFileResp, xerrors.XError)
	StartStream(ctx context.Context) (*models.Stream, xerrors.XError)
	AddChunk(ctx context.Context, req models.AddChunkReq) (*models.Chunk, xerrors.XError)
	GetStream(ctx context.Context, streamID uuid.UUID) (*models.Stream, xerrors.XError)
}

type Transport struct {
	svc    service
	log    logger.Logger
	router chi.Router
}

// New returns a new instance of Transport.
func New(log logger.Logger, svc service, mdl ...func(http.Handler) http.Handler) http.Handler {
	router := chi.NewRouter()
	router.Use(mdl...)

	h := &Transport{log: log, router: router, svc: svc}

	h.attachRoutes()

	return h
}

func (h *Transport) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.router.ServeHTTP(w, r)
}

func (h *Transport) attachRoutes() {
	h.router.Route("/upload", func(r chi.Router) {
		r.Post("/{storage}", h.uploadFile)
		r.Post("/{storage}/{compress}", h.uploadFile)
	})

	h.router.Route("/streams", func(r chi.Router) {
		r.Post("/", h.startStream)
		r.Post("/{streamID}/chunks", h.addChunk)
		r.Get("/{streamID}", h.getStream)
	})
}
