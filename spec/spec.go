package spec

type Topic int

const (
    _rawTopic       = "v1/devices/me/raw"
    _telemetryTopic = "v1/devices/me/telemetry"
    _attrTopic      = "v1/devices/me/attributes"
    _cmdTopic       = "v1/devices/me/commands"
    _cmdTopicResp   = "v1/devices/me/command/response"
)

const (
    RawTopic Topic = 1 << iota
    TelemetryTopic
    AttributeTopic
    CommandTopic
    CommandRespTopic

    MaskTopic = (CommandRespTopic << 1) - 1
)

func (t Topic) String() string {
    switch t {
    case RawTopic:
        return _rawTopic
    case TelemetryTopic:
        return _telemetryTopic
    case AttributeTopic:
        return _attrTopic
    case CommandTopic:
        return _cmdTopic
    case CommandRespTopic:
        return _cmdTopicResp
    }
    return ""
}

func (t Topic) Valid() bool {
    return (t & MaskTopic) > 0
}
