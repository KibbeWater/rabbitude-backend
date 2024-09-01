// swift-tools-version: 6.0
// The swift-tools-version declares the minimum version of Swift required to build this package.

import PackageDescription

let package = Package(
    name: "apple",
    products: [
        .library(
            name: "apple",
            targets: ["apple"])
    ],
    targets: [
        .target(
            name: "apple")
    ]
)
