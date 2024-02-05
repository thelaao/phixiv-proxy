package handlers

import (
	"archive/zip"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/thelaao/phixiv-proxy/pixiv"
	"github.com/thelaao/phixiv-proxy/utils"
)

func (h *Handler) UgoiraHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		id := chi.URLParam(r, "id")
		format := chi.URLParam(r, "format")
		output, contentType, err := h.convertUgoira(r.Context(), id, format)
		if err != nil {
			h.reportError(w, err)
			return
		}
		w.Header().Add("Content-Type", contentType)
		w.Write(output)
	}
}

func (h *Handler) convertUgoira(ctx context.Context, id string, format string) (output []byte, contentType string, err error) {
	url := fmt.Sprintf("https://www.pixiv.net/ajax/illust/%s/ugoira_meta", id)
	metadata := &pixiv.UgoiraApiResponse{}
	err = h.Client.CallApi(ctx, url, metadata)
	if err != nil {
		return
	}

	_, body, err := h.Client.Download(ctx, metadata.Body.Src)
	if err != nil {
		return
	}
	defer body.Close()
	zipRaw, err := io.ReadAll(body)
	if err != nil {
		return
	}
	zipReader, err := zip.NewReader(bytes.NewReader(zipRaw), int64(len(zipRaw)))
	if err != nil {
		return
	}

	dir, err := os.MkdirTemp("", "ugoira")
	if err != nil {
		return
	}
	defer os.RemoveAll(dir)
	err = utils.ExtractZip(zipReader, dir)
	if err != nil {
		return
	}

	var builder strings.Builder
	var lastFile string
	builder.WriteString("ffconcat version 1.0\n")
	totalMilli := 0
	cropped := false
	for frameNum, frame := range metadata.Body.Frames {
		lastFile = filepath.Join(dir, frame.File)
		builder.WriteString("file 'file:")
		builder.WriteString(lastFile)
		builder.WriteString("'\n")
		builder.WriteString("duration ")
		builder.WriteString(strconv.FormatFloat(float64(frame.Delay)/1000, 'f', 3, 64))
		builder.WriteRune('\n')
		totalMilli += frame.Delay
		if (h.UgoiraMaxFrames == frameNum+1) || (h.UgoiraMaxDuration > 0 && h.UgoiraMaxDuration <= totalMilli) {
			if len(metadata.Body.Frames) != frameNum+1 {
				cropped = true
			}
			break
		}
	}
	builder.WriteString("file 'file:")
	builder.WriteString(lastFile)
	builder.WriteString("'\n")

	ffmpegInput := strings.NewReader(builder.String())
	filename := filepath.Join(dir, "output")
	start := time.Now()
	if format == "png" {
		contentType = "image/png"
		err = utils.CallFFmpeg(ffmpegInput, "-safe", "0", "-i", "-", "-f", "apng", "-plays", "0", "-y", filename)
	} else if format == "gif" {
		contentType = "image/gif"
		err = utils.CallFFmpeg(ffmpegInput, "-safe", "0", "-i", "-", "-f", "gif", "-loop", "0", "-vf", "split[s0][s1];[s0]palettegen=stats_mode=full[p];[s1][p]paletteuse", "-y", filename)
	} else if format == "mp4" {
		contentType = "video/mp4"
		err = utils.CallFFmpeg(ffmpegInput, "-safe", "0", "-i", "-", "-f", "mp4", "-c:v", "libx264", "-profile:v", "baseline", "-pix_fmt", "yuv420p", "-an", "-vf", "pad=ceil(iw/2)*2:ceil(ih/2)*2", "-y", filename)
	} else {
		err = errors.New("invalid output format")
	}
	if err != nil {
		return
	}
	if !cropped && format == "mp4" && totalMilli < h.UgoiraMinDuration {
		loopTimes := max(1, h.UgoiraMinDuration/totalMilli)
		newFilename := filepath.Join(dir, "output_loop")
		err = utils.CallFFmpeg(nil, "-stream_loop", strconv.Itoa(loopTimes), "-i", filename, "-f", "mp4", "-c", "copy", newFilename)
		if err != nil {
			return
		}
		filename = newFilename
	}
	elapsed := time.Now().Sub(start).Milliseconds()
	log.Printf("ffmpeg time %d ms\n", elapsed)

	output, err = os.ReadFile(filename)
	return
}
