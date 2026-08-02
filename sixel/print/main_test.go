package print

import "testing"

func TestResolveSubcommand(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		args       []string
		wantErr    string
		wantRemain []string
	}{
		{
			name:    "missing subcommand",
			args:    nil,
			wantErr: "no subcommand provided; available subcommands: image, table",
		},
		{
			name:    "unknown subcommand",
			args:    []string{"other"},
			wantErr: `unknown subcommand "other"; available subcommands: image, table`,
		},
		{
			name:       "image subcommand",
			args:       []string{"image", "-i", "png"},
			wantRemain: []string{"-i", "png"},
		},
		{
			name:       "table subcommand",
			args:       []string{"table", "-i", "csv"},
			wantRemain: []string{"-i", "csv"},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			fn, remaining, err := resolveSubcommand(tt.args)
			if tt.wantErr != "" {
				if err == nil || err.Error() != tt.wantErr {
					t.Fatalf("resolveSubcommand() error = %v, want %q", err, tt.wantErr)
				}
				if fn != nil {
					t.Fatal("resolveSubcommand() returned function for error case")
				}
				return
			}

			if err != nil {
				t.Fatalf("resolveSubcommand() error = %v", err)
			}
			if fn == nil {
				t.Fatal("resolveSubcommand() returned nil function")
			}
			if len(remaining) != len(tt.wantRemain) {
				t.Fatalf("resolveSubcommand() remaining len = %d, want %d", len(remaining), len(tt.wantRemain))
			}
			for i := range remaining {
				if remaining[i] != tt.wantRemain[i] {
					t.Fatalf("resolveSubcommand() remaining[%d] = %q, want %q", i, remaining[i], tt.wantRemain[i])
				}
			}
		})
	}
}
