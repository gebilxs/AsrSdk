package main

/*

typedef void (*onStartSuccess)( char* taskId);
typedef void (*onResult)(const char * msg);
typedef void (*onWarning)(const char * code,const char * msg);
typedef void (*onError)(const char * code,const char * msg);

extern void c_onResult (onResult, const char *msg);
extern void c_onWarning (onWarning, const char *code, const char *msg);
extern void c_onError (onError, const char *code, const char *msg);
extern void c_onStartSuccess(onStartSuccess, char*taskId);

*/
import "C"

func onStartSuccess(callback C.onStartSuccess, taskId string) {
	C.c_onStartSuccess(callback, C.CString(taskId))
}

func onResult(callback C.onResult, msg string) {
	C.c_onResult(callback, C.CString(msg))
}

func onWarning(callback C.onWarning, code, msg string) {
	C.c_onWarning(callback, C.CString(code), C.CString(msg))
}

func onError(callback C.onError, code, msg string) {
	C.c_onError(callback, C.CString(code), C.CString(msg))
}
