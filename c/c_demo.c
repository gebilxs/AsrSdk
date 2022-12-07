#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <stdbool.h>
#include "libsoe.h"

char *appId = "39868f66-256c-4ff5-af51-8f2a887e85b4";
char *appSecret = "087ddf41b868880b";

const char *txtPath = "pronunciation.txt"; //文本路径
const char *audioPath = "pronunciation.wav"; //音频文件路径

const char *format = "wav"; //音频文件格式： wav/pcm/mp3
const char *mode = "word"; //评测模式：   低阶： word/sentence/chapter
//                                       高阶： qa/retell/topic

void onStartSuccessCallback() {

    const int size = 6400;
    FILE *fp = fopen(audioPath, "rb");
    if (fp == NULL) {
        printf("read wav file error.");
        return;
    }

    if (strcmp(format, "wav") == 0 || strcmp(format, "pcm") == 0) {
        //wav格式需要分包
        printf("readWav/pcm.\n");
        char buffer[size];
        int len;
        do {
            len = fread(buffer, sizeof(char), size, fp);
            if (len == size) {
                feed(appId, buffer, size);
            } else {
                char newBuffer[len];
                memcpy(newBuffer, buffer, sizeof(char) * len);
                feed(appId, newBuffer, len);
                break;
            }
        } while (len == size);
    } else if (strcmp(format, "mp3") == 0) {
        //mp3格式需要读取整个文件
        printf("readMp3.\n");
        fseek(fp, 0, SEEK_END);
        int file_size = ftell(fp);
        char buffer[file_size];
        fseek(fp, 0, SEEK_SET);
        int len = fread(buffer, sizeof(char), file_size, fp);
        feed(appId, buffer, len);
    }
    fclose(fp);
    stop(appId);
}

int readTextFile(const char *path, char *buff) {
    FILE *fp = fopen(path, "r");
    if (fp == NULL) {
        printf("read txt file error.\n");
        return -1;
    }
    fread(buff, sizeof(char), 1024, fp);
    fclose(fp);
    return 0;
}

void onResultCallback(const char *msg) {
    printf("demo,onResult:\n");
    printf("%s\n", msg);
    exit(0);
}

void onWarningCallback(const char *code, const char *msg) {
    printf("demo,onWarning:%s %s\n", code, msg);
}

void onErrorCallback(const char *code, const char *msg) {
    printf("demo,onError:%s %s\n", code, msg);
    exit(0);
}

bool isBasic(const char *mode) {
    return strcmp(mode, "word") == 0 || strcmp(mode, "sentence") == 0 || strcmp(mode, "chapter") == 0;
}

bool isAdvance(const char *mode) {
    return strcmp(mode, "qa") == 0 || strcmp(mode, "retell") == 0 || strcmp(mode, "topic") == 0;
}

int main() {


#if defined(_WIN32) || defined(_WIN64)
    //防止返回错误信息中的中文乱码
    system("chcp 65001 > NUL");
#endif


    struct Params p;
    //必须设置以下参数
    p.looseness = 4;
    p.connectTimeout = 15;
    p.responseTimeout = 15;
    p.scale = 100;
    p.sampleRate = 16000;

    p.langType = "en-US";
    p.userId = "";
    p.format = format;

    static char json[1024] = {0};
    if (isBasic(mode)) {
        char text[1024] = {0};
        int ret = readTextFile(txtPath, text);
        if (ret != 0) {
            return ret;
        }
        sprintf(json, "{\"mode\":\"%s\",\"refText\":\"%s\"}", mode, text);
    } else if (isAdvance(mode)) {
        char jsonPath[50];
        sprintf(jsonPath, "./json/%s.json", mode);
        int ret = readTextFile(jsonPath, json);
        if (ret != 0) {
            return ret;
        }
    } else {
        printf("error mode:%s", mode);
        exit(0);
    }

    p.paramsJson = json;

    start(appId, appSecret, &p, onStartSuccessCallback, onResultCallback, onWarningCallback,
          onErrorCallback);
    getchar();
}