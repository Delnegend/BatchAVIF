package libs

import (
	"fmt"
	"strings"
)

func ValidateConfig(mode string, ext []string, enc []string, fallback []string, repack []string) (string, error) {

	in_out_count := 0
	ext_has_scale := false
	enc_has_dim := false
	fallback_has_dim := false

	for _, p := range ext {
		if strings.Contains(p, "{{ input }}") {
			in_out_count++
		}
		if strings.Contains(p, "{{ output }}") {
			in_out_count++
		}
		if strings.Contains(p, "scale") {
			ext_has_scale = true
		}
	}
	for _, p := range enc {
		if strings.Contains(p, "{{ input }}") {
			in_out_count++
		}
		if strings.Contains(p, "{{ output }}") {
			in_out_count++
		}
		if strings.Contains(p, "{{ width }}") || strings.Contains(p, "{{ height }}") {
			enc_has_dim = true
		}
	}
	for _, p := range repack {
		if strings.Contains(p, "{{ input }}") {
			in_out_count++
		}
		if strings.Contains(p, "{{ output }}") {
			in_out_count++
		}
	}
	if in_out_count != 6 {
		return "", fmt.Errorf("invalid config: not enough {{ input }}/{{ output }}")
	}

	if len(fallback) > 0 {
		for _, p := range fallback {
			if strings.Contains(p, "{{ input }}") {
				in_out_count++
			}
			if strings.Contains(p, "{{ output }}") {
				in_out_count++
			}
			if strings.Contains(p, "{{ width }}") || strings.Contains(p, "{{ height }}") {
				fallback_has_dim = true
			}
		}
		if in_out_count != 8 {
			return "", fmt.Errorf("invalid config: not enough {{ input }}/{{ output }}")
		}
	}
	// validate mode
	if (mode != "file") && (mode != "pipe") {
		return "", fmt.Errorf("invalid config: mode must be file or pipe")
	}
	// no aomenc + pipe + scale
	if mode == "pipe" && ext_has_scale && (enc_has_dim || fallback_has_dim) {
		return "", fmt.Errorf("cannot scale the source with {{ width }} and/or {{ height }} in pipe mode")
	}
	return mode, nil
}
