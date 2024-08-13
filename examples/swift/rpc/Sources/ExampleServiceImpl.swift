import Foundation

class ExampleService: ExampleServiceDelegate {
    func add(_ a: Int32, _ b: Int32) -> Int32 {
        a + b
    }
    
    func subtract(_ a: Int32, _ b: Int32) -> Int32 {
        a - b
    }
}
