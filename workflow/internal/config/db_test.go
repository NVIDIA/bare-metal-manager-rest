package config

import (
	"reflect"
	"testing"
)

func TestNewDBConfig(t *testing.T) {
	type args struct {
		host     string
		port     int
		name     string
		user     string
		password string
	}

	dbcfg := DBConfig{
		Host:     "localhost",
		Port:     5432,
		Name:     "forge",
		User:     "forge",
		Password: "test123",
	}

	tests := []struct {
		name string
		args args
		want *DBConfig
	}{
		{
			name: "initialize database config",
			args: args{
				host:     dbcfg.Host,
				port:     dbcfg.Port,
				name:     dbcfg.Name,
				user:     dbcfg.User,
				password: dbcfg.Password,
			},
			want: &dbcfg,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewDBConfig(tt.args.host, tt.args.port, tt.args.name, tt.args.user, tt.args.password)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewDBConfig() = %v, want %v", got, tt.want)
			}

			if got := got.GetHostPort(); got != tt.want.GetHostPort() {
				t.Errorf("GetHostPort() = %v, want %v", got, tt.want.GetHostPort())
			}
		})
	}
}
