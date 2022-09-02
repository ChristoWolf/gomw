package logging

import "time"

const (
	// requestTemplate is the template for logging requests.
	requestTemplate = `
# Request
{{if .WithMethod}}- Method: {{.Method}}{{end}}
{{if .WithUrl}}- URL: {{.Url}}{{end}}
{{if .WithContentLength}}- Content length: {{.ContentLength}}{{end}}
{{if .WithBody}}- Body: {{.Body}}{{end}}
`

	// responseTemplate is the template for logging responses.
	responseTemplate = `
# Response
{{if .WithDuration}}- Duration: {{.Duration.Microseconds}}µs{{end}}
{{if .WithStatus}}- Status: {{.Status}}{{end}}
{{if .WithContentLength}}- Content length: {{.ContentLength}}{{end}}
{{if .WithBody}}- Body: {{.Body}}{{end}}
`
)

// requestData caches any log relevant request data
// which can be in injected into the related template.
type requestData struct {
	WithBody          bool
	WithMethod        bool
	WithUrl           bool
	WithContentLength bool
	Body              string
	Method            string
	Url               string
	ContentLength     int64
}

// responseData caches any log relevant response data
// which can be in injected into the related template.
type responseData struct {
	WithBody          bool
	WithStatus        bool
	WithContentLength bool
	WithDuration      bool
	Body              string
	Status            int
	ContentLength     int64
	Duration          time.Duration
}
