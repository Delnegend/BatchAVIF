package libs

import (
	"fmt"
	"strings"
)

func ValidateConfig(ext []string, enc []string, fallback []string, repack []string) (string, error) {
	ext_in := false
	ext_out := false
	ext_has_scale := false
	enc_in := false
	enc_out := false
	enc_has_dim := false
	fallback_has_dim := false
	repack_in := false
	repack_out := false

	for _, p := range ext {
		if strings.Contains(p, "{{ input }}") {
			ext_in = true
		}
		if strings.Contains(p, "{{ output }}") {
			ext_out = true
		}
		if strings.Contains(p, "scale") {
			ext_has_scale = true
		}
	}
	for _, p := range enc {
		if strings.Contains(p, "{{ input }}") {
			enc_in = true
		}
		if strings.Contains(p, "{{ output }}") {
			enc_out = true
		}
		if strings.Contains(p, "{{ width }}") || strings.Contains(p, "{{ height }}") {
			enc_has_dim = true
		}
	}
	for _, p := range repack {
		if strings.Contains(p, "{{ input }}") {
			repack_in = true
		}
		if strings.Contains(p, "{{ output }}") {
			repack_out = true
		}
	}

	// must have {{ input }} in extractor and {{ output }} in repackager
	if !ext_in {
		return "", fmt.Errorf("extractor: {{ input }} not specified")
	}
	if !repack_out {
		return "", fmt.Errorf("repackager: {{ output }} not specified")
	}

	// region: fallback
	if len(fallback) > 0 {
		fallback_in := false
		fallback_out := false
		for _, p := range fallback {
			if strings.Contains(p, "{{ input }}") {
				fallback_in = true
			}
			if strings.Contains(p, "{{ output }}") {
				fallback_out = true
			}
			if strings.Contains(p, "{{ width }}") || strings.Contains(p, "{{ height }}") {
				fallback_has_dim = true
			}
		}
		// fallback encoder must also use the same mode with encoder
		if (enc_in != fallback_in) || (enc_out != fallback_out) {
			return "", fmt.Errorf("fallback must be specified in the same mode as encoder")
		}
	}
	// endregion

	// determine the mode
	var mode string
	if ext_out && enc_in && enc_out && repack_in {
		mode = "file"
	}
	if !ext_out && !enc_in && !enc_out && !repack_in {
		mode = "pipe"
	}

	// no aomenc + pipe + scale
	if mode == "pipe" && ext_has_scale && (enc_has_dim || fallback_has_dim) {
		return "", fmt.Errorf("cannot scale the source with {{ width }} and/or {{ height }} in pipe mode")
	}
	if mode == "" {
		return "", fmt.Errorf("cannot determine mode")
	}
	return mode, nil
}
