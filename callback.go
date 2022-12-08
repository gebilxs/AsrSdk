package main

/*
#include <stdio.h>

typedef void (*onStartSuccess)();
typedef void (*onResult)(const char *);
typedef void (*onWarning)(const char *,const char *);
typedef void (*onError)(const char *,const char *);


void c_onStartSuccess(onStartSuccess _cb,const char *taskId){
   _cb(taskId);
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
