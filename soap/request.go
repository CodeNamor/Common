package soap

type Request interface {
	GetSoapEnvelope() Envelope
}
