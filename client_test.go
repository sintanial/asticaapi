package asticaapi

import (
	"fmt"
	"os"
	"testing"
)

var apiKey = os.Getenv("API_KEY")

func TestClient_VisionDescribe(t *testing.T) {
	client := NewClient(apiKey)

	result, err := client.VisionDescribe("https://www.astica.org/inputs/analyze_3.jpg", nil,
		VisionParameterGPT,
		//VisionParameterDescribe,
		//VisionParameterDescribeAll,
		//VisionParameterTags,
		//VisionParameterObjects,
	)
	if err != nil {
		t.Fatal(err)
	}

	fmt.Println(result)
}
