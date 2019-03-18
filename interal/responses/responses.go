//go:generate gojay -s=$GOFILE -t=Ack -o=generated_$GOFILE
package responses

type Ack struct {
	// Ack Defines if the server could acknowledge the request.
	Ack bool `gojay:"ack"json:"ack"`
}
