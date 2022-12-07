package main

/*

   typedef void (*onStartSuccess)();
   typedef void (*onResult)(const char * msg);
   typedef void (*onWarning)(const char * code,const char * msg);
   typedef void (*onError)(const char * code,const char * msg);

   struct Params{
const char* scheme
const char* addr
const char* path
   	const char* langType
   	bool enableIntermediateResult
   	int sampleRate
   	const char* format;
   	int maxSentenceSilence
   	bool enablePunctuationPrediction
   	bool enableInverseTextNormalization
   	bool enableWords
   	const char* languageModelId
   	//unsure
   	const char* hotwordsId
   	float hotwordsWeight
   	const char* correctionWordsId
   	const char* forbiddenWordsId
   	const char* paramsJson;
   };
*/
import "C"
import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"github.com/tidwall/sjson"
	"net/url"
)

func main() {}

type AsrParams struct {
	scheme                         string
	addr                           string
	path                           string
	langType                       string
	enableIntermediateResult       bool
	sampleRate                     int
	format                         string
	maxSentenceSilence             string
	enablePunctuationPrediction    bool
	enableInverseTextNormalization bool
	enableWords                    bool
	languageModelId                string
	hotwordsId                     string
	hotwordsWeight                 float64
	correctionWordsId              string
	forbiddenWordsId               string
	paramsJson                     string
	Url                            url.URL
}

const (
	ERR_CONN_FAIL, MSG_CONN_FAIL         = "23110", "连接失败"
	ERR_PARAM_ABSENCE, MSG_PARAM_ABSENCE = "23121", "参数缺失:"
	ERR_PARAM, MSG_PARAM                 = "23122", "参数错误:"
)

// 放置websocket链接的map
var connMap = make(map[string]*websocket.Conn)

//export start
func start(filename *C.char, cParams *C.struct_Params, startSuccess C.onStartSuccess,
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
	maxSentenceSilence := C.GoString(cParams.maxSentenceSilence)
	enablePunctuationPrediction := bool(cParams.enablePunctuationPrediction)
	enableInverseTextNormalization := bool(cParams.enableInverseTextNormalization)
	enableWords := bool(cParams.enableWords)
	languageModelId := C.Gostring(cParams.languageModelId)
	hotwordsId := C.Gostring(cParams.hotwordsId)
	hotwordsWeight := C.GoFloat64(cParams.hotwordsWeight)
	correctionWordsId := C.GoString(cParams.correctionWordsId)
	forbiddenWordsId := C.GoString(cParams.forbiddenWordsId)
	//paramsJson := C.GoString(cParams.paramsJson)
	//获取url
	Url := url.URL{
		//这里同样可以先做数据处理
		Scheme: C.Gostring(cParams.scheme),
		Host:   C.Gostring(cParams.addr),
		Path:   C.Gostring(cParams.path),
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
	connMap[filename] = conn
	//params中应该是payload中的参数
	params := AsrParams{
		langType:                       langType,
		enableIntermediateResult:       enableIntermediateResult,
		sampleRate:                     sampleRate,
		format:                         format,
		maxSentenceSilence:             maxSentenceSilence,
		enablePunctuationPrediction:    enablePunctuationPrediction,
		enableInverseTextNormalization: enableInverseTextNormalization,
		enableWords:                    enableWords,
		languageModelId:                languageModelId,
		hotwordsId:                     hotwordsId,
		hotwordsWeight:                 hotwordsWeight,
		correctionWordsId:              correctionWordsId,
		forbiddenWordsId:               forbiddenWordsId,
	}
	err = sendStartJson(conn, params)
	if err != nil {
		log(err.Error())
		return
	}
	go func() {
		for {
			_, msg, err := conn.ReadMessage()
			//log(string(msg))
			if err != nil {
				break
			}
			m := make(map[string]interface{})
			json.Unmarshal(msg, &m)
			header := m["header"].(map[string]interface{})
			headerName := header["name"]
			switch headerName {
			case "EvaluationStarted":
				onStartSuccess(startSuccess)
			case "EvaluationResult":
				onResult(result, string(msg))
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
	payload["enablePunctuationPrediction"] = params.enablePunctuationPrediction
	payload["enableInverseTextNormalization"] = params.enableInverseTextNormalization
	payload["enableWords"] = params.enableWords
	payload["languageModelId"] = params.languageModelId
	payload["hotwordsId"] = params.hotwordsId
	payload["hotwordsWeight"] = params.hotwordsWeight
	payload["correctionWordsId"] = params.correctionWordsId
	payload["forbiddenWordsId"] = params.forbiddenWordsId

	p["header"] = header
	p["payload"] = payload
	data, _ := json.Marshal(p)

	data, _ = sjson.SetRawBytes(data, "payload.params", []byte(params.paramsJson))

	return data
}

//export feed
func feed() {}

//export stop
func stop(taskId *C.char) {
	connMap[C.GoString(taskId)].WriteMessage(websocket.TextMessage, getStopJson())
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
