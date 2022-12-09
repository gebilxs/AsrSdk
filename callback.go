package main

/*
#include <stdio.h>

typedef void (*onStartSuccess)(char*);
typedef void (*onResult)(const char *);
//对不同的入参进行处理 总共4个
typedef void (*onSentenceBeginResult)(const char *);
typedef void (*onTranscriptionResultChangedResult)(const char *);
typedef void (*onSentenceEndResult)(const char*);
typedef void(*onTranscriptionCompletedResult)(const char*);
//
typedef void (*onWarning)(const char *,const char *);
typedef void (*onError)(const char *,const char *);



void c_onStartSuccess(onStartSuccess _cb,char *taskId){
   _cb(taskId);
}
//下方4个函数的实现
void c_onSentenceBeginResult(onSentenceBeginResult _cb,const char *msg){
	_cb(msg);
}
void c_onTranscriptionResultChangedResult(onTranscriptionResultChangedResult _cb,const char *msg){
	_cb(msg);
}
void c_onSentenceEndResult(onSentenceEndResult _cb,const char *msg){
	_cb(msg);
}
void c_onTranscriptionCompletedResult(onTranscriptionCompletedResult _cb,const char *msg){
	_cb(msg);
}
void c_onResult(onResult _cb,const char * msg){
   _cb(msg);
}

void c_onWarning(onWarning _cb,const char * code,const char * msg){
   _cb(code,msg);
}

void c_onError(onError _cb,const char * code,const char * msg){
   _cb(code,msg);
}
*/
import "C"
