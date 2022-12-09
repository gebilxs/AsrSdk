package main

/*
#include <stdio.h>
#include <errno.h>
#include <stdlib.h>
#include <string.h>
#include <stdbool.h>
   typedef void (*onStartSuccess)( char * taskId);
   typedef void (*onError)(const char * code,const char * msg);
	typedef void (*onSentenceBeginResult)(const char * msg);
	typedef void(*onTranscriptionResultChangedResult)(const char * msg);
	typedef void (*onSentenceEndResult)(const char * msg);
	typedef void (*onTranscriptionCompletedResult)(const char * msg);

   struct Params{
	const char* scheme;
	const char* addr;
	const char* path;
   	const char* langType;
   	bool enableIntermediateResult;
   	int sampleRate;
   	const char* format;
   	int maxSentenceSilence;
   	bool enableInverseTextNormalization;
   	bool enableWords;
   	const char* hotwordsId;
   	float hotwordsWeight;
   	const char* correctionWordsId;
   	const char* forbiddenWordsId;
   	const char* paramsJson;
	int serverType;
	int thread;
	bool punctuationPrediction;
	bool saveOutput;
	bool sleep;
   };

*/
import "C"
import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/tidwall/gjson"
	"net/url"
	"sync"
	"unsafe"
)

func main() {}

type AsrParams struct {
	taskId                         string
	filename                       string
	scheme                         string
	addr                           string
	path                           string
	langType                       string
	enableIntermediateResult       bool
	sampleRate                     int
	serverType                     int
	format                         string
	maxSentenceSilence             int
	enableInverseTextNormalization bool
	enableWords                    bool
	hotwordsId                     string
	hotwordsWeight                 float64
	correctionWordsId              string
	forbiddenWordsId               string
	paramsJson                     string
	thread                         int
	punctuationPrediction          bool
	saveOutput                     bool
	sleep                          bool
	Url                            url.URL
	Conn                           *websocket.Conn
}

const (
	ERR_CONN_FAIL, MSG_CONN_FAIL = "23110", "连接失败"
)

// 放置websocket链接的map

var connMap = sync.Map{}
var onErrorMap = sync.Map{}

//export start
func start(cParams *C.struct_Params, startSuccess C.onStartSuccess,
	SentenceBeginResult C.onSentenceBeginResult,
	TranscriptionResultChangedResult C.onTranscriptionResultChangedResult,
	SentenceEndResult C.onSentenceEndResult, TranscriptionCompletedResult C.onTranscriptionCompletedResult,
	error C.onError) {

	langType := C.GoString(cParams.langType)
	enableIntermediateResult := bool(cParams.enableIntermediateResult)
	sampleRate := int(cParams.sampleRate)
	format := C.GoString(cParams.format)
	enableInverseTextNormalization := bool(cParams.enableInverseTextNormalization)
	enableWords := bool(cParams.enableWords)
	hotwordsId := C.GoString(cParams.hotwordsId)
	hotwordsWeight := float64(cParams.hotwordsWeight)
	correctionWordsId := C.GoString(cParams.correctionWordsId)
	forbiddenWordsId := C.GoString(cParams.forbiddenWordsId)
	thread := int(cParams.thread)
	maxSentenceSilence := int(cParams.maxSentenceSilence)
	serverType := int(cParams.serverType)
	punctuationPrediction := bool(cParams.punctuationPrediction)
	saveOutput := bool(cParams.saveOutput)
	sleep := bool(cParams.sleep)
	path := C.GoString(cParams.path)
	//paramsJson := C.GoString(cParams.paramsJson)
	//获取url
	Url := url.URL{
		//这里同样可以先做数据处理
		Scheme: C.GoString(cParams.scheme),
		Host:   C.GoString(cParams.addr),
		Path:   "/ws/v1",
	}
	fmt.Println("connecting to", Url.String())

	//创建链接 先不处理resp
	conn, _, err := websocket.DefaultDialer.Dial(Url.String(), nil)
	if err != nil {
		log(err.Error())
		onError(error, ERR_CONN_FAIL, MSG_CONN_FAIL)
		return

	}
	//将链接存入 map 中 先将filename作为key 因为taskId是之后才返回的参数
	//connMap[filename] = conn
	//params中应该是payload中的参数

	params := AsrParams{
		langType:                       langType,
		enableIntermediateResult:       enableIntermediateResult,
		sampleRate:                     sampleRate,
		format:                         format,
		maxSentenceSilence:             maxSentenceSilence,
		enableInverseTextNormalization: enableInverseTextNormalization,
		enableWords:                    enableWords,
		hotwordsId:                     hotwordsId,
		hotwordsWeight:                 hotwordsWeight,
		correctionWordsId:              correctionWordsId,
		forbiddenWordsId:               forbiddenWordsId,
		thread:                         thread,
		serverType:                     serverType,
		punctuationPrediction:          punctuationPrediction,
		saveOutput:                     saveOutput,
		sleep:                          sleep,
		path:                           path,
	}
	err = sendStartJson(conn, params)
	if err != nil {
		log(err.Error())
		return
	}

	//对headername 进行相关的判断
	go func() {
		for {
			_, message, err := conn.ReadMessage()
			if err != nil {
				onError(error, "20199", err.Error())
				return
			}
			fmt.Println("result:", string(message))
			switch gjson.GetBytes(message, "header.name").String() {
			case "TranscriptionStarted":
				taskId := gjson.GetBytes(message, "header.task_id").String()
				//存储
				connMap.Store(taskId, conn)
				onErrorMap.Store(taskId, error)
				//新建key taskId value onerror回调

				onStartSuccess(startSuccess, taskId)
			//case "SentenceBegin", "TranscriptionResultChanged", "SentenceEnd", "TranscriptionCompleted":
			//	onResult(result, string(message))
			case "SentenceBegin":
				onSentenceBeginResult(SentenceBeginResult, string(message))
			case "TranscriptionResultChanged":
				onTranscriptionResultChangedResult(TranscriptionResultChangedResult, string(message))
			case "SentenceEnd":
				onSentenceEndResult(SentenceEndResult, string(message))
			case "TranscriptionCompleted":
				onTranscriptionCompletedResult(TranscriptionCompletedResult, string(message))
			case "TaskFailed":
				onError(error, gjson.GetBytes(message, "header.status").String(), gjson.GetBytes(message, "header.status_ext").String())

			}
		}
	}()
}

func sendStartJson(conn *websocket.Conn, params AsrParams) error {
	return conn.WriteMessage(websocket.TextMessage, getStartJson(params))
}
func getStartJson(params AsrParams) []byte {
	p := make(map[string]interface{})
	header := map[string]interface{}{
		"namespace": "SpeechTranscriber",
		"name":      "StartTranscription",
	}
	payload := make(map[string]interface{})

	payload["lang_type"] = params.langType
	payload["enable_intermediate_result"] = params.enableIntermediateResult
	payload["enable_intermediate_result"] = true
	payload["sample_rate"] = params.sampleRate
	payload["format"] = params.format
	payload["max_sentence_silence"] = params.maxSentenceSilence
	payload["enable_inverse_text_normalization"] = params.enableInverseTextNormalization
	payload["enable_words"] = params.enableWords
	payload["hotwords_id"] = params.hotwordsId
	payload["hotwords_weight"] = params.hotwordsWeight
	payload["correction_words_id"] = params.correctionWordsId
	payload["forbidden_words_id"] = params.forbiddenWordsId
	payload["enable_punctuation_prediction"] = params.punctuationPrediction

	p["header"] = header
	p["payload"] = payload
	fmt.Printf("%+v\n", p)
	data, _ := json.Marshal(p)

	return data
}

//export feed
func feed(taskId *C.char, data *C.char, length C.int) {
	//test
	//v, _ := onErrorMap.Load(C.GoString(taskId))
	//onErr := v.(C.onError)
	//onError(onErr, "20191", "testing")
	//return
	var buf []byte
	buf = C.GoBytes(unsafe.Pointer(data), length)

	v, ok := connMap.Load(C.GoString(taskId))
	if !ok {
		log("Wrong taskId")
		v, _ := onErrorMap.Load(C.GoString(taskId))
		onErr := v.(C.onError)
		onError(onErr, "20191", "Wrong taskId ...")
		return
	}

	conn, ok := v.(*websocket.Conn)
	if !ok {
		log("Wrong webSocket connection")
		v, _ := onErrorMap.Load(C.GoString(taskId))
		onErr := v.(C.onError)
		onError(onErr, "20191", "Wrong webSocket connection")
		return
	}

	err := conn.WriteMessage(websocket.BinaryMessage, buf)
	if err != nil {
		log(err.Error())
		v, _ := onErrorMap.Load(C.GoString(taskId))
		onErr := v.(C.onError)
		onError(onErr, "20191", "Failed write binary message ...")
		return
	}
}

//export stop
func stop(taskId *C.char) {
	var v interface{}
	v, ok := connMap.Load(C.GoString(taskId))
	if !ok {
		log("Failed write stop message")
		v, _ := onErrorMap.Load(C.GoString(taskId))
		onErr := v.(C.onError)
		onError(onErr, "20191", "Failed write stop message ...")
	}
	conn, ok := v.(*websocket.Conn)
	if !ok {
		log("Stop..Wrong webSocket connection")
		v, _ := onErrorMap.Load(C.GoString(taskId))
		onErr := v.(C.onError)
		onError(onErr, "20191", "Stop..Wrong websocket conn...")
	}
	conn.WriteMessage(websocket.TextMessage, getStopJson())
}
func getStopJson() []byte {
	p := make(map[string]interface{})
	header := map[string]interface{}{
		"namespace": "SpeechTranscriber",
		"name":      "StopTranscription",
	}
	p["header"] = header
	data, _ := json.Marshal(p)
	return data
}

func log(msg ...any) {
	fmt.Print("[== lib_asr ==] ")
	for _, m := range msg {
		fmt.Print(m, " ")
	}
	fmt.Println()
}
