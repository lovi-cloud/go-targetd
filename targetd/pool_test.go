package targetd

import (
	"context"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestClient_GetPoolList(t *testing.T) {
	client, err := setup()
	if err != nil {
		t.Fatalf("failed to setup client: %+v", err)
	}

	tests := []struct {
		input interface{}
		want  []Pool
		err   bool
	}{
		{
			input: nil,
			want: []Pool{
				{
					Name:     "",
					Size:     0,
					FreeSize: 0,
					Type:     "",
					UUID:     0,
				},
			},
			err: false,
		},
	}

	for _, test := range tests {
		got, err := client.GetPoolList(context.Background())
		if !test.err && err != nil {
			t.Fatalf("should not be error for %+v but: %+v", test.input, err)
		}
		if test.err && err == nil {
			t.Fatalf("should be error for %+v but not:", test.input)
		}

		if diff := cmp.Diff(test.want, got); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	}
}
