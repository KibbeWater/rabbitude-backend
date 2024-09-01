package providers

import (
	"fmt"
	"main/external"
	"main/structures"
	"time"
)

var appleProviderName = "apple"
var appleProvider = structures.Provider{
	ProviderName: appleProviderName,
}

var (
	apple_setup bool
	apple_ready bool
)

func RegisterApple(baseServices *[]structures.BaseService, services *[]structures.Service) {
	if !external.Apple_IsLoaded() {
		fmt.Println("[Warn] Apple library not loaded, skipping registration")
		return
	}

	*baseServices = append(*baseServices, structures.BaseService{
		Provider:    appleProvider,
		ServiceType: structures.SPEECH_SERVICE,
		Run:         appleSpeech,
		Setup:       appleSetup,
	})
}

func appleSetup() {
	if !external.Apple_IsLoaded() {
		fmt.Errorf("Fatal: Apple library not loaded but setup was ran")
	}

	var status int = -1
	external.Apple_RequestSpeechPermissions(&status)

	for status == -1 {
		time.Sleep(1)
	}

	switch status {
	case 0:
		fmt.Println("[Apple] Speech Permissions Authorized")
		apple_ready = true
	case 1:
		fmt.Println("[Fatal] [Apple] Speech Permissions Denied")
	case 2:
		fmt.Println("[Fatal] [Apple] Speech Permissions Restricted")
	case 3:
		fmt.Println("[Fatal] [Apple] Speech Permissions Not Determined")
	default:
		fmt.Println("[Fatal] [Apple] Speech Permissions Unknown")
	}
}

func appleSpeech(client *structures.Client, data []byte, preventDef *bool) ([]byte, error) {
	fmt.Println("Apple Speech")
	if !apple_ready {
		return nil, fmt.Errorf("Apple not ready")
	}

	// filename := utils.GenerateUniqueID() + ".wav"
	// exeDir, err := utils.GetExecutableDir()
	// if err != nil {
	// 	return nil, fmt.Errorf(fmt.Sprintf("Error getting executable directory: %s", err))
	// }
	//
	// filePath := fmt.Sprintf("%s/%s", exeDir, filename)
	//
	// // Save the audio to a file
	// if err := os.WriteFile(filePath, data, 0644); err != nil {
	// 	return nil, fmt.Errorf(fmt.Sprintf("Error writing audio to file: %s", err))
	// }
	// defer os.Remove(filePath)

	speech, err := external.Apple_SpeechRecognition(data)
	if err != nil {
		return nil, fmt.Errorf(fmt.Sprintf("Error recognizing speech: %s", err))
	}

	fmt.Println("Apple Speech:", speech)

	return []byte(speech), nil
}
