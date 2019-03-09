//go:generate gojay -s=$GOFILE -t=Ack -o=generated_$GOFILE
package responses

type Ack struct {
	Ack bool `gojay:"ack"`
}
