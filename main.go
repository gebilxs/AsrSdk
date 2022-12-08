package main

/*

   typedef void (*onStartSuccess)();
   typedef void (*onResult)(const char * msg);
   typedef void (*onWarning)(const char * code,const char * msg);
   typedef void (*onError)(const char * code,const char * msg);

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
	thread                         string
	punctuationPrediction          bool
	saveOutput                     bool
	sleep                          bool
	Url                            url.URL
	Conn                           *websocket.Conn
}

const (
	ERR_CONN_FAIL, MSG_CONN_FAIL         = "23110", "连接失败"
	ERR_PARAM_ABSENCE, MSG_PARAM_ABSENCE = "23121", "参数缺失:"
	ERR_PARAM, MSG_PARAM                 = "23122", "参数错误:"
)

// 放置websocket链接的map

var connMap = sync.Map{}

//export start
func start(cParams *C.struct_Params, startSuccess C.onStartSuccess,
	result C.onResult, warning C.onWarning, error C.onError) {

	/*处理c99版本逻辑的代码
	err := C.bool(cParams.enableIntermediateResult)
	if err != true {
		if C.boolean = 1 {
	}
	*/
	//进行统一的数据转换 c -> go
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
	thread := C.GoString(cParams.thread)
	maxSentenceSilence := int(cParams.maxSentenceSilence)
	serverType := int(cParams.serverType)
	punctuationPrediction := bool(cParams.punctuationPrediction)
	saveOutput := bool(cParams.saveOutput)
	sleep := bool(cParams.sleep)
	//paramsJson := C.GoString(cParams.paramsJson)
	//获取url
	Url := url.URL{
		//这里同样可以先做数据处理
		Scheme: C.GoString(cParams.scheme),
		Host:   C.GoString(cParams.addr),
		Path:   C.GoString(cParams.path),
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
	}
	err = sendStartJson(conn, params)
	if err != nil {
		log(err.Error())
		return
	}
	_, message, err := conn.ReadMessage()
	if err != nil {
		log(err.Error())
		return
	}
	name := gjson.GetBytes(message, "header.name").String()
	if name != "TranscriptionStarted" && name != "RecognitionStarted" {
		log(err.Error())
		return
	}
	taskId := gjson.GetBytes(message, "header.task_id").String()
	//存储
	connMap.Store(taskId, conn)

	fmt.Println(taskId)
	m := make(map[string]interface{})
	json.Unmarshal(message, &m)
	header := m["header"].(map[string]interface{})
	//对headername 进行相关的判断
	go func() {
		for {
			_, message, err = conn.ReadMessage()
			if err != nil {
				onError(error, "20199", err.Error())
				return
			}

			switch gjson.GetBytes(message, "header.name").String() {
			case "EvaluationStarted":
				onStartSuccess(startSuccess)
			case "EvaluationResult":
				onResult(result, string(message))
			case "EvaluationError":
				onError(error, header["status"].(string), header["statusText"].(string))
			case "EvaluationWarning":
				onWarning(warning, header["status"].(string), header["statusText"].(string))

			}
		}
	}()
}

func sendStartJson(conn *websocket.Conn, params AsrParams) error {
	return conn.WriteMessage(websocket.TextMessage, getStartJson(params))
}
func getStartJson(params AsrParams) []byte {
	p := make(map[string]interface{})
	header := make(map[string]interface{})
	header["namespace"] = "SpeechEvaluator"
	header["name"] = "StartEvaluation"
	payload := make(map[string]interface{})

	payload["langType"] = params.langType
	payload["enableIntermediateResult"] = params.enableIntermediateResult
	payload["sampleRate"] = params.sampleRate
	payload["format"] = params.format
	payload["maxSentenceSilence"] = params
	payload["enableInverseTextNormalization"] = params.enableInverseTextNormalization
	payload["enableWords"] = params.enableWords
	payload["hotwordsId"] = params.hotwordsId
	payload["hotwordsWeight"] = params.hotwordsWeight
	payload["correctionWordsId"] = params.correctionWordsId
	payload["forbiddenWordsId"] = params.forbiddenWordsId
	payload["thread"] = params.thread
	payload["serverType"] = params.serverType
	payload["punctuationPrediction"] = params.punctuationPrediction
	payload["saveOutput"] = params.saveOutput
	payload["sleep"] = params.sleep

	p["header"] = header
	p["payload"] = payload
	data, _ := json.Marshal(p)

	return data
}

//export feed
func feed(taskId string, data *C.char, length C.int) {
	var buf []byte
	buf = C.GoBytes(unsafe.Pointer(data), length)
	//取taskId 然后写进去
	var v interface{}
	v, _ = connMap.Load(taskId)
	conn, _ := v.(*websocket.Conn)
	conn.WriteMessage(websocket.BinaryMessage, buf)
	//用断言出来的conn,来writeMessage
	//两个返回一个value 一个bool
	//直接readMessage
}

//export stop
func stop(taskId *C.char) {
	var v interface{}
	v, _ = connMap.Load(taskId)
	conn, _ := v.(*websocket.Conn)
	conn.WriteMessage(websocket.TextMessage, getStopJson())
}
func getStopJson() []byte {
	p := make(map[string]interface{})
	header := make(map[string]interface{})
	header["namespace"] = "SpeechEvaluator"
	header["name"] = "StopEvaluation"
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
