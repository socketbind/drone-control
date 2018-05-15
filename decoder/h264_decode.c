#include <libavutil/log.h>
#include <libavcodec/avcodec.h>

#include "h264_decode.h"

extern void handleFrame(AVFrame*);

int h264dec_new(h264dec_t *h) {
    avcodec_register_all();

    h->pkt = av_packet_alloc();
    if (!h->pkt) {
        fprintf(stderr, "unable to allocate packet\n");
        return -1;
    }

    h->inbuf = (uint8_t*) malloc(INBUF_SIZE + AV_INPUT_BUFFER_PADDING_SIZE);
    if (!h->inbuf) {
        fprintf(stderr, "unable to malloc input buffer\n");
        return -1;
    }

    memset(h->inbuf + INBUF_SIZE, 0, AV_INPUT_BUFFER_PADDING_SIZE);

    h->c = avcodec_find_decoder(AV_CODEC_ID_H264);
    if (!h->c) {
        fprintf(stderr, "unable to find the specified codec\n");
        return -1;
    }

    h->parser = av_parser_init(h->c->id);
    if (!h->parser) {
        fprintf(stderr, "unable to init the parser\n");
        return -1;
    }

    h->ctx = avcodec_alloc_context3(h->c);
    if (!h->ctx) {
        fprintf(stderr, "unable to allocate context\n");
        return -1;
    }

    h->f = av_frame_alloc();
    if (!h->f) {
        fprintf(stderr, "unable to allocate frame\n");
        return -1;
    }

    return avcodec_open2(h->ctx, h->c, 0);
}

void h264dec_free(h264dec_t *h) {
    avcodec_send_packet(h->ctx, NULL);

    free(h->inbuf);

    av_parser_close(h->parser);
    avcodec_free_context(&h->ctx);
    av_frame_free(&h->f);
    av_packet_free(&h->pkt);
}

int h264dec_decode(h264dec_t *h, uint8_t *input, int input_size) {
    int ret, data_size;
    uint8_t *data;

    while (input_size > 0) {
        if (input_size > INBUF_SIZE) {
            data_size = INBUF_SIZE;
        } else {
            data_size = input_size;
        }

        memcpy(h->inbuf, input, data_size);

        input      += data_size;
        input_size -= data_size;

        data = h->inbuf;
        while (data_size > 0) {
            ret = av_parser_parse2(h->parser, h->ctx, &h->pkt->data, &h->pkt->size,
                                   data, data_size, AV_NOPTS_VALUE, AV_NOPTS_VALUE, 0);
            if (ret < 0) {
                return -1;
            }

            data      += ret;
            data_size -= ret;

            if (h->pkt->size) {
                ret = avcodec_send_packet(h->ctx, h->pkt);
                if (ret < 0) {
                    return -1;
                }

                while (ret >= 0) {
                    ret = avcodec_receive_frame(h->ctx, h->f);
                    if (ret == AVERROR(EAGAIN) || ret == AVERROR_EOF) {
                        break;
                    } else if (ret < 0) {
                        return -1;
                    }

                    // h->f now usable
                    handleFrame(h->f);
                }
            }
        }
    }

    return 0;
}
