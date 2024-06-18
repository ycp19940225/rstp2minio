package minio

import (
	"context"
	"testing"
)

// endpoint: '120.27.225.90:9000'
// bucket: 'admin'
// key: 'admin'
// secret: 'zenwell123456'
var (
	m, _ = New(
		WithEndpoint("120.27.225.90:9000"),
		WithKey("admin"),
		WithSecret("zenwell123456"),
		WithBucket("test"),
	)
)

func TestMinio_Token(t *testing.T) {
	type fields struct {
		m *Minio
	}
	type args struct {
		ctx    context.Context
		object string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "case1",
			fields: fields{
				m: m,
			},
			args: args{
				ctx:    context.Background(),
				object: "1.jpg",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotToken, err := tt.fields.m.Token(tt.args.ctx, tt.args.object)
			if (err != nil) != tt.wantErr {
				t.Errorf("Token() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("Token() gotToken = %v", gotToken)
		})
	}
}

func TestMinio_Get(t *testing.T) {
	type fields struct {
		m *Minio
	}
	type args struct {
		ctx    context.Context
		object string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "case1",
			fields: fields{
				m: m,
			},
			args: args{
				ctx:    context.Background(),
				object: "camera/upload/1.xlsx",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotReply, err := m.Get(tt.args.ctx, tt.args.object)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("Get() gotReply = %v", gotReply)
		})
	}
}

func TestMinio_Preview(t *testing.T) {
	type fields struct {
		m *Minio
	}
	type args struct {
		ctx    context.Context
		object string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "case1",
			fields: fields{
				m: m,
			},
			args: args{
				ctx:    context.Background(),
				object: "1.png",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotReply, err := m.Preview(tt.args.ctx, tt.args.object)
			if (err != nil) != tt.wantErr {
				t.Errorf("Preview() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("Preview() gotReply = %v", gotReply)
		})
	}
}
