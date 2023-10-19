package fill

import (
	"encoding/json"
	"io"
	"time"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/dynamicpb"
)

// SilentFilter is a Filler implementation that doesn't behave interactive actions.
type SilentFiller struct {
	dec *protojson.UnmarshalOptions
	in  *json.Decoder
}

// NewSilentFiller receives input as io.Reader and returns an instance of SilentFiller.
func NewSilentFiller(in io.Reader) *SilentFiller {
	return &SilentFiller{
		dec: &protojson.UnmarshalOptions{
			Resolver: nil, // TODO
		},
		in: json.NewDecoder(in),
	}
}

type EvansSleep struct {
	Duration time.Duration `json:"evans_wait_ms"`
}

// Fill fills values of each field from a JSON string. If the JSON string is invalid JSON format or v is a nil pointer,
// Fill returns ErrCodecMismatch.
func (f *SilentFiller) Fill(v *dynamicpb.Message) error {
	var in interface{}
	if err := f.in.Decode(&in); err != nil {
		return err
	}

	b, err := json.Marshal(in)
	if err != nil {
		return err
	}

	var cmd EvansSleep
	err = json.Unmarshal(b, &cmd)
	if err == nil && cmd.Duration > 0 {
		time.Sleep(cmd.Duration * time.Millisecond)
		return nil
	} else {
		return f.dec.Unmarshal(b, v)
	}
}
