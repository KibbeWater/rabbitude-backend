// The Swift Programming Language
// https://docs.swift.org/swift-book
import AVFoundation
import Foundation
import Speech

enum SpeechPermissions: Int {
    case authorized = 0
    case denied = 1
    case restricted = 2
    case notDetermined = 3
}

@_cdecl("swift_greet")
public func swift_greet(name: UnsafePointer<CChar>) -> UnsafeMutablePointer<CChar>? {
    let swiftName = String(cString: name)
    let greeting = "Hello, \(swiftName)!"
    return strdup(greeting)  // strdup returns a C string allocated on the heap
}

@_cdecl("swift_requestSpeechPermissions")
public func swift_requestSpeechPermissions(statusPtr: UnsafeMutablePointer<Int>) {
    print("Requesting speech recognition permissions")
    SFSpeechRecognizer.requestAuthorization { authStatus in
        print("[Swift] Result returned")
        print("[Swift] Speech recognition permissions status: \(authStatus.rawValue)")
        switch authStatus {
        case .authorized:
            statusPtr.pointee = SpeechPermissions.authorized.rawValue
        case .denied:
            statusPtr.pointee = SpeechPermissions.denied.rawValue
        case .restricted:
            statusPtr.pointee = SpeechPermissions.restricted.rawValue
        case .notDetermined:
            statusPtr.pointee = SpeechPermissions.notDetermined.rawValue
        @unknown default:
            statusPtr.pointee = SpeechPermissions.notDetermined.rawValue
        }
    }
}

@_cdecl("swift_speechRecognition")
public func swift_speechRecognition(buffer: UnsafePointer<UInt8>, length: Int)
    -> UnsafeMutablePointer<CChar>?
{
    print("[Swift] Starting speech recognition")
    let data = Data(bytes: buffer, count: length)
    let recognizer = SFSpeechRecognizer(locale: Locale(identifier: "en-US"))

    print("[Swift] Creating recognition request")
    let request = SFSpeechAudioBufferRecognitionRequest()

    var recognizedText = ""

    do {
        print("[Swift] Creating temporary file")
        let tempDirectoryURL = FileManager.default.temporaryDirectory
        let tempFileURL = tempDirectoryURL.appendingPathComponent(UUID().uuidString + ".wav")
        try data.write(to: tempFileURL)
        print("[Swift] Temporary file created at: \(tempFileURL.path)")

        print("[Swift] Reading audio file")
        let audioFile = try AVAudioFile(forReading: tempFileURL)
        let audioFormat = audioFile.processingFormat
        let audioFrameCount = UInt32(audioFile.length)
        print(
            "[Swift] Audio file details - Format: \(audioFormat), Frame count: \(audioFrameCount)")

        print("[Swift] Creating audio buffer")
        let audioFileBuffer = AVAudioPCMBuffer(
            pcmFormat: audioFormat, frameCapacity: audioFrameCount)

        print("[Swift] Appending audio buffer to request")
        request.append(audioFileBuffer!)
        request.shouldReportPartialResults = true
        request.requiresOnDeviceRecognition = true
        request.endAudio()

        print("[Swift] Starting recognition task")
        var shouldRun = true
        if recognizer != nil {
            print("[Swift] Recognizer is available")
        } else {
            print("[Swift] Recognizer is not available")
        }
        if recognizer!.isAvailable {
            print("[Swift] Recognizer2 is available")
        } else {
            print("[Swift] Recognizer2 is not available")
        }
        recognizer?.recognitionTask(with: request) { (result, error) in
            print("[Swift] Recognition task callback")

            guard let result = result else {
                print(
                    "[Swift] Recognition failed with error: \(error?.localizedDescription ?? "Unknown error")"
                )
                return
            }

            if result.isFinal {
                recognizedText = result.bestTranscription.formattedString
                print("[Swift] Final recognized text: \(recognizedText)")
                shouldRun = false
            } else {
                print("[Swift] Interim result: \(result.bestTranscription.formattedString)")
            }
        }

        while shouldRun {
            RunLoop.current.run(until: Date(timeIntervalSinceNow: 0.1))
        }
        print("[Swift] Recognition completed or timed out")
    } catch {
        print("[Swift] Error processing audio file: \(error.localizedDescription)")
        return strdup(
            "Unfortunately there is a problem with my speech recognition, please try again later.")
    }

    print("[Swift] Returning recognized text")
    return strdup(recognizedText)
}
