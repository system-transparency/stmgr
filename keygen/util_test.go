package keygen

import (
	"encoding/pem"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadPEM(t *testing.T) {
	for _, table := range []struct {
		name    string
		data    []byte
		wantErr error
	}{
		{
			name: "success",
			data: []byte(
				"-----BEGIN PRIVATE KEY-----\n" +
					"MC4CAQAwBQYDK2VwBCIEIKrt4aIojIpXfdmw2aVWHNCNGZDvDSL+t1CI6STXjh7F\n" +
					"-----END PRIVATE KEY-----\n",
			),
			wantErr: nil,
		},
		{
			name:    "no pem",
			data:    []byte("not pem data"),
			wantErr: ErrNoPEMBlock,
		},
		{
			name: "trailing",
			data: []byte(
				"-----BEGIN PRIVATE KEY-----\n" +
					"MC4CAQAwBQYDK2VwBCIEIKrt4aIojIpXfdmw2aVWHNCNGZDvDSL+t1CI6STXjh7F\n" +
					"-----END PRIVATE KEY-----\n" +
					"trailing data",
			),
			wantErr: ErrTrailing,
		},
	} {
		t.Run(table.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "test.pem")

			want, rest := pem.Decode(table.data)
			if len(rest) != 0 {
				want = nil
			}

			if err := os.WriteFile(path, table.data, os.ModePerm); err != nil {
				t.Errorf("Failed writing file: %v", err)
			}

			if got, err := LoadPEM(path); !errors.Is(err, table.wantErr) {
				t.Errorf("LoadPEM() err = %q, want %q", err, table.wantErr)
			} else if !reflect.DeepEqual(got, want) {
				t.Errorf("LoadPEM() = %q, want %q", got, want)
			}
		})
	}
}

func TestWritePEM(t *testing.T) {
	for _, table := range []struct {
		name    string
		data    []byte
		wantErr error
	}{
		{
			name: "success",
			data: []byte(
				"-----BEGIN PRIVATE KEY-----\n" +
					"MC4CAQAwBQYDK2VwBCIEIKrt4aIojIpXfdmw2aVWHNCNGZDvDSL+t1CI6STXjh7F\n" +
					"-----END PRIVATE KEY-----\n",
			),
			wantErr: nil,
		},
		{
			name:    "no pem",
			data:    []byte("not pem data"),
			wantErr: ErrNoPEMBlock,
		},
	} {
		t.Run(table.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "test.pem")

			want, rest := pem.Decode(table.data)
			if len(rest) != 0 {
				want = nil
			}

			if err := WritePEM(want, path); !errors.Is(err, table.wantErr) {
				t.Errorf("WritePEM() err = %q, want %q", err, table.wantErr)
			}

			if table.wantErr == nil {
				bytes, err := os.ReadFile(path)
				if err != nil {
					t.Errorf("Failed to read file: %v", err)
				}

				if got, rest := pem.Decode(bytes); len(rest) != 0 {
					t.Errorf("WritePEM() = %q + %q, want %q", got, rest, want)
				} else if !reflect.DeepEqual(got, want) {
					t.Errorf("WritePEM() = %q, want %q", got, want)
				}
			}
		})
	}
}
