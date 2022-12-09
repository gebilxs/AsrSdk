package main

/*
#include <stdio.h>

typedef void (*onStartSuccess)(char*);
typedef void (*onSentenceBeginResult)(const char *);
typedef void (*onTranscriptionResultChangedResult)(const char *);
typedef void (*onSentenceEndResult)(const char*);
typedef void(*onTranscriptionCompletedResult)(const char*);
typedef void (*onError)(const char *,const char *);

void c_onStartSuccess(onStartSuccess _cb,char *taskId){
   _cb(taskId);
}
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
void c_onError(onError _cb,const char * code,const char * msg){
   _cb(code,msg);
}
*/
import "C"
