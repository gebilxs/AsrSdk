package main

/*

typedef void (*onStartSuccess)(  char* taskId);
typedef void(*onSentenceBeginResult)(const char* msg);
typedef void(*onTranscriptionResultChangedResult)(const char*msg);
typedef void(*onSentenceEndResult)(const char *msg);
typedef void(*onTranscriptionCompletedResult)(const char *msg);
typedef void (*onError)(const char * code,const char * msg);

extern void c_onError (onError, const char *code, const char *msg);
extern void c_onStartSuccess(onStartSuccess,  char*taskId);
extern void c_onSentenceBeginResult(onSentenceBeginResult,const char *msg);
extern void c_onTranscriptionResultChangedResult(onTranscriptionResultChangedResult,const char *msg);
extern void c_onSentenceEndResult(onSentenceEndResult,const char *msg);
extern void c_onTranscriptionCompletedResult(onTranscriptionCompletedResult,const char*msg);

*/
import "C"
import "fmt"

func onStartSuccess(callback C.onStartSuccess, taskId string) {
	fmt.Println(&callback, callback)

	C.c_onStartSuccess(callback, C.CString(taskId))
}
func onSentenceBeginResult(callback C.onSentenceBeginResult, msg string) {
	C.c_onSentenceBeginResult(callback, C.CString(msg))
}
func onTranscriptionResultChangedResult(callback C.onTranscriptionResultChangedResult, msg string) {
	C.c_onTranscriptionResultChangedResult(callback, C.CString(msg))
}
func onSentenceEndResult(callback C.onSentenceEndResult, msg string) {
	C.c_onSentenceEndResult(callback, C.CString(msg))
}
func onTranscriptionCompletedResult(callback C.onTranscriptionCompletedResult, msg string) {
	C.c_onTranscriptionCompletedResult(callback, C.CString(msg))
}
func onError(callback C.onError, code, msg string) {
	C.c_onError(callback, C.CString(code), C.CString(msg))

}
