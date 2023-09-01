package asticaapi

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

const ModelVersion2_1_full = "2.1_full"

type Client struct {
	apiKey string
}

// Returns a caption which describes the image.
const VisionParameterDescribe = "describe"

// Returns multiple auxilliary captions that describe the image.
const VisionParameterDescribeAll = "describe_all"

// Returns the results of OCR with positional coordinates.
const VisionParameterTextRead = "text_read"

// Uses the result of asticaVision to create a GPT description. Using this VisionParameter increases the processing time of your API request. Be Patient.
const VisionParameterGPT = "gpt"

// Uses the result of asticaVision to create a GPT-4 description. Using this VisionParameter greatly increases the processing time of your API request. Please be patient.
const VisionParameterGPTDetailed = "gpt_detailed"

// Returns the age and gender of all faces detected in the image.
const VisionParameterFaces = "faces"

// Returns a calculated value for different types of sensitive materials found in the image.
const VisionParameterModerate = "moderate"

// Returns a list of descriptive terms which describe the image.
const VisionParameterTags = "tags"

// Returns a list of brands that have been identified. For example, a logos on a cup, or a t-shirt.
const VisionParameterBrands = "brands"

// Returns a list of celebrities and other known persons that have been detected in the photo.
const VisionParameterCelebrities = "celebrities"

// Returns a list of known locations and landmarks found in the photo. For example, the Eiffel Tower
const VisionParameterLandmarks = "landmarks"

const VisionParameterObjects = "objects"

func NewClient(apiKey string) *Client {
	return &Client{apiKey: apiKey}
}

type Rectangle struct {
	X int `json:"x"`
	Y int `json:"y"`
	W int `json:"w"`
	H int `json:"h"`
}

type tokenRequest struct {
	Token string `json:"tkn"`
}

type VisionDescribeRequest struct {
	tokenRequest
	ModelVersion string `json:"modelVersion"`
	Input        string `json:"input"`
	VisionParams string `json:"visionParams"`
	GptPrompt    string `json:"gpt_prompt"`
	GtpLength    string `json:"gtp_length"`
}

type VisionDescribeResponse struct {
	ModelVersion string `json:"modelVersion"`
	Astica       struct {
		Request      string  `json:"request"`
		RequestType  string  `json:"requestType"`
		ModelVersion string  `json:"modelVersion"`
		ApiQty       float64 `json:"api_qty"`
	} `json:"astica"`
	Status      string `json:"status"`
	CaptionGPTS string `json:"caption_GPTS"`
	GPTLevel    int    `json:"GPT_level"`
	Caption     struct {
		Text       string  `json:"text"`
		Confidence float64 `json:"confidence"`
	} `json:"caption"`
	CaptionList []struct {
		Text       string    `json:"text"`
		Confidence float64   `json:"confidence"`
		Rectangle  Rectangle `json:"rectangle"`
	} `json:"caption_list"`
	Objects []struct {
		Name       string    `json:"name"`
		Confidence float64   `json:"confidence"`
		Rectangle  Rectangle `json:"rectangle"`
	} `json:"objects"`
	Tags []struct {
		Name       string  `json:"name"`
		Confidence float64 `json:"confidence"`
	} `json:"tags"`
	Metadata struct {
		Width  int `json:"width"`
		Height int `json:"height"`
	} `json:"metadata"`
}

type VisionDescribeOptions struct {
	GptPrompt string
	GptLength string
}

func (self *Client) VisionDescribe(image string, options *VisionDescribeOptions, visionParameters ...string) (*VisionDescribeResponse, error) {
	vdr := VisionDescribeRequest{
		tokenRequest: tokenRequest{
			Token: self.apiKey,
		},
		ModelVersion: ModelVersion2_1_full,
		Input:        image,
		VisionParams: strings.Join(visionParameters, ", "),
	}

	if options != nil {
		vdr.GptPrompt = options.GptPrompt
		vdr.GtpLength = options.GptLength
	}

	var reqbody bytes.Buffer
	if err := json.NewEncoder(&reqbody).Encode(vdr); err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, "https://vision.astica.ai/describe", &reqbody)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if res.StatusCode != http.StatusOK {
		return nil, errors.New("Invalid status code: " + res.Status)
	}

	defer res.Body.Close()
	resbody, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	var result VisionDescribeResponse
	if err := json.Unmarshal(resbody, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
