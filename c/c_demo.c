#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <stdbool.h>
#include "libasr.h"


const char *scheme = "ws";
const char *addr ="localhost:7100";
const char *path ="C:/Users/Administrator/Desktop/chengdu.wav";
const char *langType = "zh-cmn-Hans-CN";
const char *format="wav";
const char *hotwordsId="default";
const char *hotwordsWeight="0.7";
const char *correctionWordsId="";
const char *forbiddenWordsId="";
const int sampleRate=16000;
const int thread = 1;
const int maxSentenceSilence=450;
const int serverType=1;
bool enableIntermediateResult =true;
bool enableInverseTextNormalization=true;
bool enableWords = false;
bool  saveOutput = false;
bool sleep =false;
bool punctuationPrediction=true;
void onStartSuccessCallback(const char *taskId) {
    printf("on start success ...\n");
    const int size = 6400;
    FILE *fp = fopen(path, "rb");
    if (fp == NULL) {
        printf("read wav file error.");
        return;
    }

    if (strcmp(format, "wav") == 0 || strcmp(format, "pcm") == 0) {
        //wav格式需要分包
        printf("read wav/pcm.\n");
        char buffer[size];
        int len;
        do {
            len = fread(buffer, sizeof(char), size, fp);
//            printf("len:%d,size:%d\n",len,size);
            if (len == size) {
                feed(taskId, buffer, size);
            } else {
                char newBuffer[len];
                memcpy(newBuffer, buffer, sizeof(char) * len);
                feed(taskId, newBuffer, len);
                break;
            }

        } while (len == size);
    }
    fclose(fp);
    stop(taskId);
    printf("stop ...\n");
}


void onResultCallback(const char *msg) {
    printf("demo,onResult:\n");
    printf("%s\n", msg);
}

void onWarningCallback(const char *code, const char *msg) {
    printf("demo,onWarning:%s %s\n", code, msg);
}

void onErrorCallback(const char *code, const char *msg) {
    printf("demo,onError:%s %s\n", code, msg);
    exit(0);
}

int main() {


#if defined(_WIN32) || defined(_WIN64)
    //防止返回错误信息中的中文乱码
    system("chcp 65001 > NUL");
#endif


    struct Params p;
    //必须设置以下参数
    p.scheme=scheme;
    p.addr=addr;
    p.path=path;
    p.enableIntermediateResult = enableIntermediateResult;
    p.sampleRate = sampleRate;
    p.langType = langType;
    p.format=format;
    p.maxSentenceSilence=maxSentenceSilence;
    p.enableInverseTextNormalization=enableInverseTextNormalization;
    p.enableWords=enableWords;
    p.hotwordsId=hotwordsId;
    p.hotwordsWeight=0;
    p.forbiddenWordsId=forbiddenWordsId;
    p.thread=thread;
    p.serverType=serverType;
    p.punctuationPrediction=punctuationPrediction;
    p.saveOutput=saveOutput;
    p.sleep=sleep;
    p.correctionWordsId=correctionWordsId;

    static char json[1024] = {0};

    p.paramsJson = json;

    start(&p, onStartSuccessCallback, onResultCallback, onWarningCallback,
          onErrorCallback);
    getchar();
 printf("end ...");
}