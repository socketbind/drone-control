#ifndef __H264_DECODE_H
#define __H264_DECODE_H

#include <libavcodec/avcodec.h>

#define INBUF_SIZE 4096

typedef struct {
    AVPacket *pkt;
    uint8_t *inbuf;
    AVCodec *c;
    AVCodecParserContext *parser;
    AVCodecContext *ctx;
    AVFrame *f;
} h264dec_t;

int h264dec_new(h264dec_t *h);
void h264dec_free(h264dec_t *h);
int h264dec_decode(h264dec_t *h, uint8_t *input, int input_size);

#endif