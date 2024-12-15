package common

import (
	"encoding/json"
	"strings"
	"sync"
)

// from songquanpeng/one-api
const (
	USD2RMB = 7.3           // 暂定 1 USD = 7.3 RMB
	USD     = 500           // $0.002 = 1 -> $1 = 500 一美元所对应的基础倍率数
	RMB2    = 1             // 锚定 人民币的定价
	RMB     = USD / USD2RMB // 1 RMB 对应基础倍率
)

// modelRatio
// https://platform.openai.com/docs/models/model-endpoint-compatibility
// https://cloud.baidu.com/doc/WENXINWORKSHOP/s/Blfmc9dlf
// https://openai.com/pricing
// TODO: when a new api is enabled, check the pricing here
// 1 === $0.002 / 1K tokens
// 1 === ￥0.014 / 1k tokens

var defaultModelRatio = map[string]float64{
	//"midjourney":                50,
	"gpt-4-gizmo-*":  15,
	"gpt-4o-gizmo-*": 2.5,
	"gpt-4-all":      15,
	"gpt-4o-all":     15,
	"gpt-4":          15,
	//"gpt-4-0314":                   15, //deprecated
	"gpt-4-0613": 15,
	"gpt-4-32k":  30,
	//"gpt-4-32k-0314":               30, //deprecated
	"gpt-4-32k-0613":                  30,
	"gpt-4-1106-preview":              5,    // $10 / 1M tokens
	"gpt-4-0125-preview":              5,    // $10 / 1M tokens
	"gpt-4-turbo-preview":             5,    // $10 / 1M tokens
	"gpt-4-vision-preview":            5,    // $10 / 1M tokens
	"gpt-4-1106-vision-preview":       5,    // $10 / 1M tokens
	"chatgpt-4o-latest":               2.5,  // $5 / 1M tokens
	"gpt-4o":                          1.25, // $2.5 / 1M tokens
	"gpt-4o-audio-preview":            1.25, // $2.5 / 1M tokens
	"gpt-4o-audio-preview-2024-10-01": 1.25, // $2.5 / 1M tokens
	"gpt-4o-2024-05-13":               2.5,  // $5 / 1M tokens
	"gpt-4o-2024-08-06":               1.25, // $2.5 / 1M tokens
	"gpt-4o-2024-11-20":               1.25, // $2.5 / 1M tokens
	"gpt-4o-realtime-preview":         2.5,
	"o1-preview":                      7.5,
	"o1-preview-2024-09-12":           7.5,
	"o1-mini":                         1.5,
	"o1-mini-2024-09-12":              1.5,
	"gpt-4o-mini":                     0.075,
	"gpt-4o-mini-2024-07-18":          0.075,
	"gpt-4-turbo":                     5, // $0.01 / 1K tokens
	"gpt-4-turbo-2024-04-09":          5, // $0.01 / 1K tokens
	//"gpt-3.5-turbo-0301":           0.75, //deprecated
	"gpt-3.5-turbo":          0.25,
	"gpt-3.5-turbo-0613":     0.75,
	"gpt-3.5-turbo-16k":      1.5, // $0.003 / 1K tokens
	"gpt-3.5-turbo-16k-0613": 1.5,
	"gpt-3.5-turbo-instruct": 0.75, // $0.0015 / 1K tokens
	"gpt-3.5-turbo-1106":     0.5,  // $0.001 / 1K tokens
	"gpt-3.5-turbo-0125":     0.25,
	"babbage-002":            0.2, // $0.0004 / 1K tokens
	"davinci-002":            1,   // $0.002 / 1K tokens
	"text-ada-001":           0.2,
	"text-babbage-001":       0.25,
	"text-curie-001":         1,
	//"text-davinci-002":               10,
	//"text-davinci-003":               10,
	"text-davinci-edit-001":        10,
	"code-davinci-edit-001":        10,
	"whisper-1":                    15,  // $0.006 / minute -> $0.006 / 150 words -> $0.006 / 200 tokens -> $0.03 / 1k tokens
	"tts-1":                        7.5, // 1k characters -> $0.015
	"tts-1-1106":                   7.5, // 1k characters -> $0.015
	"tts-1-hd":                     15,  // 1k characters -> $0.03
	"tts-1-hd-1106":                15,  // 1k characters -> $0.03
	"davinci":                      10,
	"curie":                        10,
	"babbage":                      10,
	"ada":                          10,
	"text-embedding-3-small":       0.01,
	"text-embedding-3-large":       0.065,
	"text-embedding-ada-002":       0.05,
	"text-search-ada-doc-001":      10,
	"text-moderation-stable":       0.1,
	"text-moderation-latest":       0.1,
	"claude-instant-1":             0.4,   // $0.8 / 1M tokens
	"claude-2.0":                   4,     // $8 / 1M tokens
	"claude-2.1":                   4,     // $8 / 1M tokens
	"claude-3-haiku-20240307":      0.125, // $0.25 / 1M tokens
	"claude-3-5-haiku-20241022":    0.5,   // $1 / 1M tokens
	"claude-3-sonnet-20240229":     1.5,   // $3 / 1M tokens
	"claude-3-5-sonnet-20240620":   1.5,
	"claude-3-5-sonnet-20241022":   1.5,
	"claude-3-opus-20240229":       7.5, // $15 / 1M tokens
	"BLOOMZ-7B":                    0.004 * RMB,
	"Embedding-V1":                 0.002 * RMB,
	"bge-large-zh":                 0.002 * RMB,
	"bge-large-en":                 0.002 * RMB,
	"tao-8k":                       0.002 * RMB,
	"PaLM-2":                       1,
	"gemini-pro":                   1, // $0.00025 / 1k characters -> $0.001 / 1k tokens
	"gemini-pro-vision":            1, // $0.00025 / 1k characters -> $0.001 / 1k tokens
	"gemini-1.0-pro-vision-001":    1,
	"gemini-1.0-pro-001":           1,
	"gemini-1.5-pro-latest":        1.75, // $3.5 / 1M tokens
	"gemini-1.5-pro-exp-0827":      1.75, // $3.5 / 1M tokens
	"gemini-1.5-flash-latest":      1,
	"gemini-1.5-flash-exp-0827":    1,
	"gemini-1.0-pro-latest":        1,
	"gemini-1.0-pro-vision-latest": 1,
	"gemini-ultra":                 1,
	// ￥0.002 / 1k tokens

	// Perplexity online 模型对搜索额外收费，有需要应自行调整，此处不计入搜索费用
	"llama-3-sonar-small-32k-chat":   0.2 / 1000 * USD,
	"llama-3-sonar-small-32k-online": 0.2 / 1000 * USD,
	"llama-3-sonar-large-32k-chat":   1 / 1000 * USD,
	"llama-3-sonar-large-32k-online": 1 / 1000 * USD,
}

var defaultModelPrice = map[string]float64{
	"suno_music":        0.1,
	"suno_lyrics":       0.01,
	"dall-e-3":          0.04,
	"gpt-4-gizmo-*":     0.1,
	"mj_imagine":        0.1,
	"mj_variation":      0.1,
	"mj_reroll":         0.1,
	"mj_blend":          0.1,
	"mj_modal":          0.1,
	"mj_zoom":           0.1,
	"mj_shorten":        0.1,
	"mj_high_variation": 0.1,
	"mj_low_variation":  0.1,
	"mj_pan":            0.1,
	"mj_inpaint":        0,
	"mj_custom_zoom":    0,
	"mj_describe":       0.05,
	"mj_upscale":        0.05,
	"swap_face":         0.05,
	"mj_upload":         0.05,
}

var (
	modelPriceMap      map[string]float64 = nil
	modelPriceMapMutex                    = sync.RWMutex{}
)
var (
	modelRatioMap      map[string]float64 = nil
	modelRatioMapMutex                    = sync.RWMutex{}
)

var CompletionRatio map[string]float64 = nil
var defaultCompletionRatio = map[string]float64{
	"gpt-4-gizmo-*":  2,
	"gpt-4o-gizmo-*": 3,
	"gpt-4-all":      2,
}

func GetModelPriceMap() map[string]float64 {
	modelPriceMapMutex.Lock()
	defer modelPriceMapMutex.Unlock()
	if modelPriceMap == nil {
		modelPriceMap = defaultModelPrice
	}
	return modelPriceMap
}

func ModelPrice2JSONString() string {
	GetModelPriceMap()
	jsonBytes, err := json.Marshal(modelPriceMap)
	if err != nil {
		SysError("error marshalling model price: " + err.Error())
	}
	return string(jsonBytes)
}

func UpdateModelPriceByJSONString(jsonStr string) error {
	modelPriceMapMutex.Lock()
	defer modelPriceMapMutex.Unlock()
	modelPriceMap = make(map[string]float64)
	return json.Unmarshal([]byte(jsonStr), &modelPriceMap)
}

// GetModelPrice 返回模型的价格，如果模型不存在则返回-1，false
func GetModelPrice(name string, printErr bool) (float64, bool) {
	GetModelPriceMap()
	if strings.HasPrefix(name, "gpt-4-gizmo") {
		name = "gpt-4-gizmo-*"
	}
	if strings.HasPrefix(name, "gpt-4o-gizmo") {
		name = "gpt-4o-gizmo-*"
	}
	price, ok := modelPriceMap[name]
	if !ok {
		if printErr {
			SysError("model price not found: " + name)
		}
		return -1, false
	}
	return price, true
}

func GetModelRatioMap() map[string]float64 {
	modelRatioMapMutex.Lock()
	defer modelRatioMapMutex.Unlock()
	if modelRatioMap == nil {
		modelRatioMap = defaultModelRatio
	}
	return modelRatioMap
}

func ModelRatio2JSONString() string {
	GetModelRatioMap()
	jsonBytes, err := json.Marshal(modelRatioMap)
	if err != nil {
		SysError("error marshalling model ratio: " + err.Error())
	}
	return string(jsonBytes)
}

func UpdateModelRatioByJSONString(jsonStr string) error {
	modelRatioMapMutex.Lock()
	defer modelRatioMapMutex.Unlock()
	modelRatioMap = make(map[string]float64)
	return json.Unmarshal([]byte(jsonStr), &modelRatioMap)
}

func GetModelRatio(name string) float64 {
	GetModelRatioMap()
	if strings.HasPrefix(name, "gpt-4-gizmo") {
		name = "gpt-4-gizmo-*"
	}
	ratio, ok := modelRatioMap[name]
	if !ok {
		SysError("model ratio not found: " + name)
		return 30
	}
	return ratio
}

func DefaultModelRatio2JSONString() string {
	jsonBytes, err := json.Marshal(defaultModelRatio)
	if err != nil {
		SysError("error marshalling model ratio: " + err.Error())
	}
	return string(jsonBytes)
}

func GetDefaultModelRatioMap() map[string]float64 {
	return defaultModelRatio
}

func CompletionRatio2JSONString() string {
	if CompletionRatio == nil {
		CompletionRatio = defaultCompletionRatio
	}
	jsonBytes, err := json.Marshal(CompletionRatio)
	if err != nil {
		SysError("error marshalling completion ratio: " + err.Error())
	}
	return string(jsonBytes)
}

func UpdateCompletionRatioByJSONString(jsonStr string) error {
	CompletionRatio = make(map[string]float64)
	return json.Unmarshal([]byte(jsonStr), &CompletionRatio)
}

func GetCompletionRatio(name string) float64 {
	if strings.HasPrefix(name, "gpt-4-gizmo") {
		name = "gpt-4-gizmo-*"
	}
	if strings.HasPrefix(name, "gpt-4o-gizmo") {
		name = "gpt-4o-gizmo-*"
	}
	if strings.HasPrefix(name, "gpt-4") && !strings.HasSuffix(name, "-all") && !strings.HasSuffix(name, "-gizmo-*") {
		if strings.HasPrefix(name, "gpt-4o") {
			if name == "gpt-4o-2024-05-13" {
				return 3
			}
			return 4
		}
		if strings.HasPrefix(name, "gpt-4-turbo") || strings.HasSuffix(name, "preview") {
			return 3
		}
		return 2
	}
	if strings.HasPrefix(name, "o1-") {
		return 4
	}
	if name == "chatgpt-4o-latest" {
		return 3
	}
	if strings.Contains(name, "claude-instant-1") {
		return 3
	} else if strings.Contains(name, "claude-2") {
		return 3
	} else if strings.Contains(name, "claude-3") {
		return 5
	}
	if strings.HasPrefix(name, "gpt-3.5") {
		if name == "gpt-3.5-turbo" || strings.HasSuffix(name, "0125") {
			// https://openai.com/blog/new-embedding-models-and-api-updates
			// Updated GPT-3.5 Turbo model and lower pricing
			return 3
		}
		if strings.HasSuffix(name, "1106") {
			return 2
		}
		return 4.0 / 3.0
	}
	if strings.HasPrefix(name, "mistral-") {
		return 3
	}
	if strings.HasPrefix(name, "gemini-") {
		return 4
	}
	if strings.HasPrefix(name, "command") {
		switch name {
		case "command-r":
			return 3
		case "command-r-plus":
			return 5
		case "command-r-08-2024":
			return 4
		case "command-r-plus-08-2024":
			return 4
		default:
			return 2
		}
	}
	if strings.HasPrefix(name, "deepseek") {
		return 2
	}
	if strings.HasPrefix(name, "ERNIE-Speed-") {
		return 2
	} else if strings.HasPrefix(name, "ERNIE-Lite-") {
		return 2
	} else if strings.HasPrefix(name, "ERNIE-Character") {
		return 2
	} else if strings.HasPrefix(name, "ERNIE-Functions") {
		return 2
	}
	switch name {
	case "llama2-70b-4096":
		return 0.8 / 0.64
	case "llama3-8b-8192":
		return 2
	case "llama3-70b-8192":
		return 0.79 / 0.59
	}
	if ratio, ok := CompletionRatio[name]; ok {
		return ratio
	}
	return 1
}

func GetAudioRatio(name string) float64 {
	if strings.HasPrefix(name, "gpt-4o-realtime") {
		return 20
	} else if strings.HasPrefix(name, "gpt-4o-audio") {
		return 40
	}
	return 20
}

func GetAudioCompletionRatio(name string) float64 {
	if strings.HasPrefix(name, "gpt-4o-realtime") {
		return 2
	}
	return 2
}

//func GetAudioPricePerMinute(name string) float64 {
//	if strings.HasPrefix(name, "gpt-4o-realtime") {
//		return 0.06
//	}
//	return 0.06
//}
//
//func GetAudioCompletionPricePerMinute(name string) float64 {
//	if strings.HasPrefix(name, "gpt-4o-realtime") {
//		return 0.24
//	}
//	return 0.24
//}

func GetCompletionRatioMap() map[string]float64 {
	if CompletionRatio == nil {
		CompletionRatio = defaultCompletionRatio
	}
	return CompletionRatio
}