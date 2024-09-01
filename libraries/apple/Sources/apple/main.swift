// The Swift Programming Language
// https://docs.swift.org/swift-book
import Foundation

@_cdecl("swift_greet")
public func swift_greet(name: UnsafePointer<CChar>) -> UnsafeMutablePointer<CChar>? {
    let swiftName = String(cString: name)
    let greeting = "Hello, \(swiftName)!"
    return strdup(greeting)  // strdup returns a C string allocated on the heap
}
