package main

import (
	"bytes"
	"fmt"
	"io"
	"strings"
	"testing"
)

func Test_scanStream(t *testing.T) {
	type args struct {
		r io.Reader
	}
	tests := []struct {
		name  string
		args  args
		wantW string
	}{
		// TODO: Add test cases.
		{
			name: "test1",
			args: args{
				r: strings.NewReader(`data: {"choices":[],"created":0,"id":"","prompt_filter_results":[{"content_filter_results":{"hate":{"filtered":false,"severity":"safe"},"self_harm":{"filtered":false,"severity":"safe"},"sexual":{"filtered":false,"severity":"safe"},"violence":{"filtered":false,"severity":"safe"}},"prompt_index":0}]}

data: {"choices":[{"index":0,"content_filter_offsets":{"check_offset":95,"start_offset":95,"end_offset":129},"content_filter_results":{"hate":{"filtered":false,"severity":"safe"},"self_harm":{"filtered":false,"severity":"safe"},"sexual":{"filtered":false,"severity":"safe"},"violence":{"filtered":false,"severity":"safe"}},"delta":{"content":null,"role":"assistant"}}],"created":1706846880,"id":"chatcmpl-8nf8KfSBXm4NzcgQOvtjoSyPZOSun"}

data: {"choices":[{"index":0,"content_filter_offsets":{"check_offset":95,"start_offset":95,"end_offset":129},"content_filter_results":{"hate":{"filtered":false,"severity":"safe"},"self_harm":{"filtered":false,"severity":"safe"},"sexual":{"filtered":false,"severity":"safe"},"violence":{"filtered":false,"severity":"safe"}},"delta":{"content":"Hello","role":null}}],"created":1706846880,"id":"chatcmpl-8nf8KfSBXm4NzcgQOvtjoSyPZOSun"}

data: {"choices":[{"index":0,"content_filter_offsets":{"check_offset":95,"start_offset":95,"end_offset":129},"content_filter_results":{"hate":{"filtered":false,"severity":"safe"},"self_harm":{"filtered":false,"severity":"safe"},"sexual":{"filtered":false,"severity":"safe"},"violence":{"filtered":false,"severity":"safe"}},"delta":{"content":"!","role":null}}],"created":1706846880,"id":"chatcmpl-8nf8KfSBXm4NzcgQOvtjoSyPZOSun"}

data: {"choices":[{"index":0,"content_filter_offsets":{"check_offset":95,"start_offset":95,"end_offset":129},"content_filter_results":{"hate":{"filtered":false,"severity":"safe"},"self_harm":{"filtered":false,"severity":"safe"},"sexual":{"filtered":false,"severity":"safe"},"violence":{"filtered":false,"severity":"safe"}},"delta":{"content":" How","role":null}}],"created":1706846880,"id":"chatcmpl-8nf8KfSBXm4NzcgQOvtjoSyPZOSun"}

data: {"choices":[{"index":0,"content_filter_offsets":{"check_offset":95,"start_offset":95,"end_offset":129},"content_filter_results":{"hate":{"filtered":false,"severity":"safe"},"self_harm":{"filtered":false,"severity":"safe"},"sexual":{"filtered":false,"severity":"safe"},"violence":{"filtered":false,"severity":"safe"}},"delta":{"content":" can","role":null}}],"created":1706846880,"id":"chatcmpl-8nf8KfSBXm4NzcgQOvtjoSyPZOSun"}

data: {"choices":[{"index":0,"content_filter_offsets":{"check_offset":95,"start_offset":95,"end_offset":129},"content_filter_results":{"hate":{"filtered":false,"severity":"safe"},"self_harm":{"filtered":false,"severity":"safe"},"sexual":{"filtered":false,"severity":"safe"},"violence":{"filtered":false,"severity":"safe"}},"delta":{"content":" I","role":null}}],"created":1706846880,"id":"chatcmpl-8nf8KfSBXm4NzcgQOvtjoSyPZOSun"}

data: {"choices":[{"index":0,"content_filter_offsets":{"check_offset":95,"start_offset":95,"end_offset":129},"content_filter_results":{"hate":{"filtered":false,"severity":"safe"},"self_harm":{"filtered":false,"severity":"safe"},"sexual":{"filtered":false,"severity":"safe"},"violence":{"filtered":false,"severity":"safe"}},"delta":{"content":" assist","role":null}}],"created":1706846880,"id":"chatcmpl-8nf8KfSBXm4NzcgQOvtjoSyPZOSun"}

data: {"choices":[{"index":0,"content_filter_offsets":{"check_offset":95,"start_offset":95,"end_offset":129},"content_filter_results":{"hate":{"filtered":false,"severity":"safe"},"self_harm":{"filtered":false,"severity":"safe"},"sexual":{"filtered":false,"severity":"safe"},"violence":{"filtered":false,"severity":"safe"}},"delta":{"content":" you","role":null}}],"created":1706846880,"id":"chatcmpl-8nf8KfSBXm4NzcgQOvtjoSyPZOSun"}

data: {"choices":[{"index":0,"content_filter_offsets":{"check_offset":95,"start_offset":95,"end_offset":129},"content_filter_results":{"hate":{"filtered":false,"severity":"safe"},"self_harm":{"filtered":false,"severity":"safe"},"sexual":{"filtered":false,"severity":"safe"},"violence":{"filtered":false,"severity":"safe"}},"delta":{"content":" today","role":null}}],"created":1706846880,"id":"chatcmpl-8nf8KfSBXm4NzcgQOvtjoSyPZOSun"}

data: {"choices":[{"index":0,"content_filter_offsets":{"check_offset":95,"start_offset":95,"end_offset":129},"content_filter_results":{"hate":{"filtered":false,"severity":"safe"},"self_harm":{"filtered":false,"severity":"safe"},"sexual":{"filtered":false,"severity":"safe"},"violence":{"filtered":false,"severity":"safe"}},"delta":{"content":"?","role":null}}],"created":1706846880,"id":"chatcmpl-8nf8KfSBXm4NzcgQOvtjoSyPZOSun"}

data: {"choices":[{"finish_reason":"stop","index":0,"content_filter_offsets":{"check_offset":95,"start_offset":95,"end_offset":129},"content_filter_results":{"hate":{"filtered":false,"severity":"safe"},"self_harm":{"filtered":false,"severity":"safe"},"sexual":{"filtered":false,"severity":"safe"},"violence":{"filtered":false,"severity":"safe"}},"delta":{"content":null,"role":null}}],"created":1706846880,"id":"chatcmpl-8nf8KfSBXm4NzcgQOvtjoSyPZOSun"}

data: [DONE]`),
			},
			wantW: "test",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := &bytes.Buffer{}
			scanStream(tt.args.r, w)
			fmt.Println(w.String())
			// if gotW := w.String(); gotW != tt.wantW {
			// 	t.Errorf("scanStream() = %v, want %v", gotW, tt.wantW)
			// }
		})
	}
}
