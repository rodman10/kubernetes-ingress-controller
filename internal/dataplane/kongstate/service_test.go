package kongstate

import (
	"io"
	"reflect"
	"testing"

	"github.com/kong/go-kong/kong"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"

	kongv1 "github.com/kong/kubernetes-ingress-controller/v2/pkg/apis/configuration/v1"
)

func TestOverrideService(t *testing.T) {
	assert := assert.New(t)

	testTable := []struct {
		inService      Service
		inKongIngresss kongv1.KongIngress
		outService     Service
		inAnnotation   map[string]string
	}{
		{
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("http"),
					Path:     kong.String("/"),
				},
			},
			kongv1.KongIngress{
				Proxy: &kongv1.KongIngressService{},
			},
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("http"),
					Path:     kong.String("/"),
				},
			},
			map[string]string{},
		},
		{
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("http"),
					Path:     kong.String("/"),
				},
			},
			kongv1.KongIngress{
				Proxy: &kongv1.KongIngressService{
					Protocol: kong.String("https"),
				},
			},
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("https"),
					Path:     kong.String("/"),
				},
			},
			map[string]string{},
		},
		{
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("http"),
					Path:     kong.String("/"),
				},
			},
			kongv1.KongIngress{
				Proxy: &kongv1.KongIngressService{
					Retries: kong.Int(0),
				},
			},
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("http"),
					Path:     kong.String("/"),
					Retries:  kong.Int(0),
				},
			},
			map[string]string{},
		},
		{
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("http"),
					Path:     kong.String("/"),
				},
			},
			kongv1.KongIngress{
				Proxy: &kongv1.KongIngressService{
					Path: kong.String("/new-path"),
				},
			},
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("http"),
					Path:     kong.String("/new-path"),
				},
			},
			map[string]string{},
		},
		{
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("http"),
					Path:     kong.String("/"),
				},
			},
			kongv1.KongIngress{
				Proxy: &kongv1.KongIngressService{
					Retries: kong.Int(1),
				},
			},
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("http"),
					Path:     kong.String("/"),
					Retries:  kong.Int(1),
				},
			},
			map[string]string{},
		},
		{
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("http"),
					Path:     kong.String("/"),
				},
			},
			kongv1.KongIngress{
				Proxy: &kongv1.KongIngressService{
					ConnectTimeout: kong.Int(100),
					ReadTimeout:    kong.Int(100),
					WriteTimeout:   kong.Int(100),
				},
			},
			Service{
				Service: kong.Service{
					Host:           kong.String("foo.com"),
					Port:           kong.Int(80),
					Name:           kong.String("foo"),
					Protocol:       kong.String("http"),
					Path:           kong.String("/"),
					ConnectTimeout: kong.Int(100),
					ReadTimeout:    kong.Int(100),
					WriteTimeout:   kong.Int(100),
				},
			},
			map[string]string{},
		},
		{
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("grpc"),
					Path:     nil,
				},
			},
			kongv1.KongIngress{
				Proxy: &kongv1.KongIngressService{
					Protocol: kong.String("grpc"),
				},
			},
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("grpc"),
					Path:     nil,
				},
			},
			map[string]string{},
		},
		{
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("https"),
					Path:     nil,
				},
			},
			kongv1.KongIngress{
				Proxy: &kongv1.KongIngressService{
					Protocol: kong.String("grpcs"),
				},
			},
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("grpcs"),
					Path:     nil,
				},
			},
			map[string]string{},
		},
		{
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("https"),
					Path:     kong.String("/"),
				},
			},
			kongv1.KongIngress{
				Proxy: &kongv1.KongIngressService{
					Protocol: kong.String("grpcs"),
				},
			},
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("grpcs"),
					Path:     nil,
				},
			},
			map[string]string{"konghq.com/protocol": "grpcs"},
		},
		{
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("https"),
					Path:     kong.String("/"),
				},
			},
			kongv1.KongIngress{
				Proxy: &kongv1.KongIngressService{
					Protocol: kong.String("grpcs"),
				},
			},
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("grpc"),
					Path:     nil,
				},
			},
			map[string]string{"konghq.com/protocol": "grpc"},
		},
		{
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("https"),
					Path:     kong.String("/"),
				},
			},
			kongv1.KongIngress{
				Proxy: &kongv1.KongIngressService{},
			},
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("grpcs"),
					Path:     nil,
				},
			},
			map[string]string{"konghq.com/protocol": "grpcs"},
		},
		{
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("https"),
					Path:     kong.String("/"),
				},
			},
			kongv1.KongIngress{
				Proxy: &kongv1.KongIngressService{
					Protocol: kong.String("grpcs"),
				},
			},
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("https"),
					Path:     kong.String("/"),
				},
			},
			map[string]string{"konghq.com/protocol": "https"},
		},
		{
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("https"),
					Path:     kong.String("/"),
				},
			},
			kongv1.KongIngress{
				Proxy: &kongv1.KongIngressService{},
			},
			Service{
				Service: kong.Service{
					Host:     kong.String("foo.com"),
					Port:     kong.Int(80),
					Name:     kong.String("foo"),
					Protocol: kong.String("https"),
					Path:     kong.String("/"),
				},
			},
			map[string]string{"konghq.com/protocol": "https"},
		},
	}

	for _, testcase := range testTable {
		testcase := testcase
		log := logrus.New()
		log.SetOutput(io.Discard)

		k8sServices := testcase.inService.K8sServices
		for _, svc := range k8sServices {
			testcase.inService.override(log, &testcase.inKongIngresss, svc)
			assert.Equal(testcase.inService, testcase.outService)
		}
	}

	assert.NotPanics(func() {
		log := logrus.New()
		log.SetOutput(io.Discard)

		var nilService *Service
		nilService.override(log, nil, nil)
	})
}

func TestOverrideServicePath(t *testing.T) {
	type args struct {
		service Service
		anns    map[string]string
	}
	tests := []struct {
		name string
		args args
		want Service
	}{
		{},
		{name: "basic empty service"},
		{
			name: "set to valid value",
			args: args{
				anns: map[string]string{
					"konghq.com/path": "/foo",
				},
			},
			want: Service{
				Service: kong.Service{
					Path: kong.String("/foo"),
				},
			},
		},
		{
			name: "does not set path if doesn't start with /",
			args: args{
				anns: map[string]string{
					"konghq.com/path": "foo",
				},
			},
			want: Service{},
		},
		{
			name: "overrides any other value",
			args: args{
				service: Service{
					Service: kong.Service{
						Path: kong.String("/foo"),
					},
				},
				anns: map[string]string{
					"konghq.com/path": "/bar",
				},
			},
			want: Service{
				Service: kong.Service{
					Path: kong.String("/bar"),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.service.overridePath(tt.args.anns)
			if !reflect.DeepEqual(tt.args.service, tt.want) {
				t.Errorf("overrideServicePath() got = %v, want %v", tt.args.service, tt.want)
			}
		})
	}
}

func TestOverrideConnectTimeout(t *testing.T) {
	type args struct {
		service Service
		anns    map[string]string
	}
	tests := []struct {
		name string
		args args
		want Service
	}{
		{
			name: "set to valid value",
			args: args{
				anns: map[string]string{
					"konghq.com/connect-timeout": "3000",
				},
			},
			want: Service{
				Service: kong.Service{
					ConnectTimeout: kong.Int(3000),
				},
			},
		},
		{
			name: "value cannot parse to int",
			args: args{
				anns: map[string]string{
					"konghq.com/connect-timeout": "burranyi yedigei",
				},
			},
			want: Service{},
		},
		{
			name: "overrides any other value",
			args: args{
				service: Service{
					Service: kong.Service{
						ConnectTimeout: kong.Int(2000),
					},
				},
				anns: map[string]string{
					"konghq.com/connect-timeout": "3000",
				},
			},
			want: Service{
				Service: kong.Service{
					ConnectTimeout: kong.Int(3000),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.service.overrideConnectTimeout(tt.args.anns)
			if !reflect.DeepEqual(tt.args.service, tt.want) {
				t.Errorf("overrideConnectTimeout() got = %v, want %v", tt.args.service, tt.want)
			}
		})
	}
}

func TestOverrideWriteTimeout(t *testing.T) {
	type args struct {
		service Service
		anns    map[string]string
	}
	tests := []struct {
		name string
		args args
		want Service
	}{
		{
			name: "set to valid value",
			args: args{
				anns: map[string]string{
					"konghq.com/write-timeout": "3000",
				},
			},
			want: Service{
				Service: kong.Service{
					WriteTimeout: kong.Int(3000),
				},
			},
		},
		{
			name: "value cannot parse to int",
			args: args{
				anns: map[string]string{
					"konghq.com/write-timeout": "burranyi yedigei",
				},
			},
			want: Service{},
		},
		{
			name: "overrides any other value",
			args: args{
				service: Service{
					Service: kong.Service{
						WriteTimeout: kong.Int(2000),
					},
				},
				anns: map[string]string{
					"konghq.com/write-timeout": "3000",
				},
			},
			want: Service{
				Service: kong.Service{
					WriteTimeout: kong.Int(3000),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.service.overrideWriteTimeout(tt.args.anns)
			if !reflect.DeepEqual(tt.args.service, tt.want) {
				t.Errorf("overrideWriteTimeout() got = %v, want %v", tt.args.service, tt.want)
			}
		})
	}
}

func TestOverrideReadTimeout(t *testing.T) {
	type args struct {
		service Service
		anns    map[string]string
	}
	tests := []struct {
		name string
		args args
		want Service
	}{
		{
			name: "set to valid value",
			args: args{
				anns: map[string]string{
					"konghq.com/read-timeout": "3000",
				},
			},
			want: Service{
				Service: kong.Service{
					ReadTimeout: kong.Int(3000),
				},
			},
		},
		{
			name: "value cannot parse to int",
			args: args{
				anns: map[string]string{
					"konghq.com/read-timeout": "burranyi yedigei",
				},
			},
			want: Service{},
		},
		{
			name: "overrides any other value",
			args: args{
				service: Service{
					Service: kong.Service{
						ReadTimeout: kong.Int(2000),
					},
				},
				anns: map[string]string{
					"konghq.com/read-timeout": "3000",
				},
			},
			want: Service{
				Service: kong.Service{
					ReadTimeout: kong.Int(3000),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.service.overrideReadTimeout(tt.args.anns)
			if !reflect.DeepEqual(tt.args.service, tt.want) {
				t.Errorf("overrideReadTimeout() got = %v, want %v", tt.args.service, tt.want)
			}
		})
	}
}

func TestOverrideRetries(t *testing.T) {
	type args struct {
		service Service
		anns    map[string]string
	}
	tests := []struct {
		name string
		args args
		want Service
	}{
		{
			name: "set to valid value",
			args: args{
				anns: map[string]string{
					"konghq.com/retries": "3000",
				},
			},
			want: Service{
				Service: kong.Service{
					Retries: kong.Int(3000),
				},
			},
		},
		{
			name: "value cannot parse to int",
			args: args{
				anns: map[string]string{
					"konghq.com/retries": "burranyi yedigei",
				},
			},
			want: Service{},
		},
		{
			name: "overrides any other value",
			args: args{
				service: Service{
					Service: kong.Service{
						Retries: kong.Int(2000),
					},
				},
				anns: map[string]string{
					"konghq.com/retries": "3000",
				},
			},
			want: Service{
				Service: kong.Service{
					Retries: kong.Int(3000),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.args.service.overrideRetries(tt.args.anns)
			if !reflect.DeepEqual(tt.args.service, tt.want) {
				t.Errorf("overrideRetries() got = %v, want %v", tt.args.service, tt.want)
			}
		})
	}
}
