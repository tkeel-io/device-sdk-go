package device_sdk_go

// Adaptor is the interface that describes an adaptor
type Adaptor interface {
    // Name returns the label for the Adaptor
    Name() string
    // SetName sets the label for the Adaptor
    SetName(n string)
    // Connect initiates the Adaptor
    Connect() error
    // Finalize terminates the Adaptor
    Finalize() error
}

