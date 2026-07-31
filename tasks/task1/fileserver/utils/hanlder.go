package utils

import (
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"sort"

	"go.uber.org/zap"
)

type (
	FileHandler struct{}
)

var (
	Files        = regexp.MustCompile(`^/files$`)
	FileWithName = regexp.MustCompile(`^/files/[^/]+$`)
)

const temp = "temp"

func (h *FileHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logger, _ := zap.NewProduction()

	query := r.URL.RequestURI()

	switch {
	case Files.MatchString(query):
		switch r.Method {
		case http.MethodPost:
			if err := r.ParseMultipartForm(10 << 20); err != nil {
				logger.Error("parse multipart form", zap.Error(err))
				w.WriteHeader(http.StatusConflict)
				return
			}

			file, header, err := r.FormFile("file")

			if err != nil {
				logger.Error("form file", zap.Error(err))
				w.WriteHeader(http.StatusConflict)
				return
			}

			defer func() {
				if err := file.Close(); err != nil {
					logger.Error("file close", zap.Error(err))
				}
			}()

			path := "./" + temp + "/" + header.Filename

			logger.Info("Path info", zap.String("path", path))

			_, err = os.Open(path)

			if !os.IsNotExist(err) {
				logger.Error("open file", zap.Error(err))
				w.WriteHeader(http.StatusConflict)
				return
			}

			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				logger.Error("make dir", zap.Error(err))
				w.WriteHeader(http.StatusConflict)
				return
			}

			dst, err := os.Create(path)
			if err != nil {
				logger.Error("create file", zap.Error(err))
				w.WriteHeader(http.StatusConflict)
				return
			}

			defer func() {
				if err := dst.Close(); err != nil {
					logger.Error("close file", zap.Error(err))
				}
			}()

			if _, err := io.Copy(dst, file); err != nil {
				logger.Error("write file", zap.Error(err))
				w.WriteHeader(http.StatusConflict)
				return
			}

			w.WriteHeader(http.StatusCreated)
			if _, err := w.Write([]byte(header.Filename)); err != nil {
				logger.Error("write response", zap.Error(err))
				return
			}

			return
		case http.MethodGet:
			files, _ := os.ReadDir("./" + temp)

			names := make([]string, 0)

			for _, file := range files {
				if !file.IsDir() {
					names = append(names, file.Name())
				}
			}

			sort.Strings(names)

			for _, name := range names {
				if _, err := w.Write([]byte(name + "\n")); err != nil {
					return
				}
			}

			return
		default:
			logger.Error("bad files request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	case FileWithName.MatchString(query):
		switch r.Method {
		case http.MethodPut:
			name := r.URL.Path[len("/files/"):]

			path := "./" + temp + "/" + name

			dst, err := os.OpenFile(path, os.O_WRONLY|os.O_TRUNC, 0644)

			if os.IsNotExist(err) {
				logger.Error("update file", zap.Error(err))
				w.WriteHeader(http.StatusNotFound)
				return
			}

			defer func() {
				if err := dst.Close(); err != nil {
					logger.Error("close file", zap.Error(err))
				}
			}()

			if err := r.ParseMultipartForm(10 << 20); err != nil {
				w.WriteHeader(http.StatusConflict)
				return
			}

			src, _, err := r.FormFile("file")

			if err != nil {
				w.WriteHeader(http.StatusConflict)
				return
			}

			defer func() {
				if err := src.Close(); err != nil {
					logger.Error("close file", zap.Error(err))
				}
			}()

			if _, err := io.Copy(dst, src); err != nil {
				logger.Error("write file", zap.Error(err))
				w.WriteHeader(http.StatusConflict)
				return
			}

			w.WriteHeader(http.StatusOK)
			return
		case http.MethodGet:
			name := r.URL.Path[len("/files/"):]

			path := "./" + temp + "/" + name

			data, err := os.ReadFile(path)
			if err != nil {
				logger.Error("read file", zap.Error(err))
				if os.IsNotExist(err) {
					w.WriteHeader(http.StatusNotFound)
				} else {
					w.WriteHeader(http.StatusConflict)
				}
				return
			}
			if _, err := w.Write(data); err != nil {
				w.WriteHeader(http.StatusConflict)
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		case http.MethodDelete:
			name := r.URL.Path[len("/files/"):]

			path := "./" + temp + "/" + name

			if err := os.Remove(path); err != nil {
				logger.Error("file already deleted", zap.Error(err))
				return
			}
			w.WriteHeader(http.StatusOK)
			return
		default:
			logger.Error("bad file with name request")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	default:
		logger.Error("bad request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
}
