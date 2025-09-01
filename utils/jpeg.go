package utils

/*
#cgo LDFLAGS: -lturbojpeg

#include <stdlib.h>
#include "turbojpeg.h"

#define PIXEL_FORMAT TJPF_RGB

int recompressJpeg(unsigned char *jpegInput, unsigned long inputSize, unsigned char **outBuf, unsigned long *outLen, int jpegQual)
{
    int retcode = 0, width, height, jpegSubsamp, jpegColorspace;
    unsigned char *imageBuf = NULL;
    tjhandle handle = NULL;

    *outBuf = NULL;
    *outLen = 0;

    if (!(handle = tjInitTransform()))
    {
        retcode = 1;
        goto cleanup;
    }
    if (tjDecompressHeader3(handle, jpegInput, inputSize, &width, &height, &jpegSubsamp, &jpegColorspace))
    {
        retcode = 1;
        goto cleanup;
    }
    if (!(imageBuf = malloc(width * height * tjPixelSize[PIXEL_FORMAT])))
    {
        retcode = 1;
        goto cleanup;
    }
    if (tjDecompress2(handle, jpegInput, inputSize, imageBuf, width, 0, height, PIXEL_FORMAT, 0))
    {
        retcode = 1;
        goto cleanup;
    }
    if (tjCompress2(handle, imageBuf, width, 0, height, PIXEL_FORMAT, outBuf, outLen, TJSAMP_444, jpegQual, 0))
    {
        retcode = 2;
        goto cleanup;
    }

cleanup:
    if (handle)
    {
        tjDestroy(handle);
    }
    if (imageBuf)
    {
        free(imageBuf);
    }
    if (retcode && *outBuf)
    {
        tjFree(*outBuf);
        *outBuf = NULL;
    }
    return retcode;
}

void freeJpeg(unsigned char *buf)
{
    if (buf)
    {
        tjFree(buf);
    }
}
*/
import "C"

import (
	"unsafe"
)

func ReencodeJPEG(imageBytes []byte) []byte {
	cBytes := C.CBytes(imageBytes)
	defer C.free(cBytes)
	var outBuf *C.uchar
	var outLen C.ulong

	res := C.recompressJpeg((*C.uchar)(cBytes), C.ulong(len(imageBytes)), &outBuf, &outLen, C.int(85))
	if int(res) == 0 {
		ret := C.GoBytes(unsafe.Pointer(outBuf), C.int(outLen))
		C.freeJpeg(outBuf)
		if len(ret) < len(imageBytes) {
			return ret
		}
	}
	return imageBytes
}
