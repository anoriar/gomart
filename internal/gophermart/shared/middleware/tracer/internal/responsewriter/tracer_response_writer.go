package responsewriter

import "net/http"

type ResponseData struct {
	status  int
	payload []byte
}

func (r ResponseData) Status() int {
	return r.status
}

func (r ResponseData) Payload() []byte {
	return r.payload
}

func NewResponseData() *ResponseData {
	return &ResponseData{
		0, nil,
	}
}

type TracerResponseWriter struct {
	http.ResponseWriter
	responseData *ResponseData
}

func (l *TracerResponseWriter) ResponseData() *ResponseData {
	return l.responseData
}

func NewTracerResponseWriter(w http.ResponseWriter) *TracerResponseWriter {
	return &TracerResponseWriter{
		w,
		NewResponseData(),
	}
}

func (l *TracerResponseWriter) Write(bytes []byte) (int, error) {
	size, err := l.ResponseWriter.Write(bytes)
	l.responseData.payload = bytes
	return size, err
}

func (l *TracerResponseWriter) WriteHeader(statusCode int) {
	l.ResponseWriter.WriteHeader(statusCode)
	l.responseData.status = statusCode
}
