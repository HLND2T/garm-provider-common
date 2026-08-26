// Copyright 2026 Cloudbase Solutions SRL
//
// Licensed under the Apache License, Version 2.0 (the "License");

package execution

import (
	"context"
	"testing"

	common "github.com/cloudbase/garm-provider-common/execution/common"
	executionv010 "github.com/cloudbase/garm-provider-common/execution/v0.1.0"
	"github.com/cloudbase/garm-provider-common/params"
)

type topLevelNamedProvider struct {
	legacyDeleteCalled bool
	instanceID         string
	instanceName       string
}

func (*topLevelNamedProvider) CreateInstance(context.Context, params.BootstrapInstance) (params.ProviderInstance, error) {
	return params.ProviderInstance{}, nil
}

func (p *topLevelNamedProvider) DeleteInstance(context.Context, string) error {
	p.legacyDeleteCalled = true
	return nil
}

func (p *topLevelNamedProvider) DeleteInstanceWithName(_ context.Context, instanceID, instanceName string) error {
	p.instanceID = instanceID
	p.instanceName = instanceName
	return nil
}

func (*topLevelNamedProvider) GetInstance(context.Context, string) (params.ProviderInstance, error) {
	return params.ProviderInstance{}, nil
}

func (*topLevelNamedProvider) ListInstances(context.Context, string) ([]params.ProviderInstance, error) {
	return nil, nil
}

func (*topLevelNamedProvider) RemoveAllInstances(context.Context) error { return nil }
func (*topLevelNamedProvider) Stop(context.Context, string, bool) error { return nil }
func (*topLevelNamedProvider) Start(context.Context, string) error      { return nil }
func (*topLevelNamedProvider) GetVersion(context.Context) string        { return "test" }

func TestRunPreservesNamedDeleteCapabilityThroughVersionInterface(t *testing.T) {
	concreteProvider := &topLevelNamedProvider{}
	var provider executionv010.ExternalProvider = concreteProvider
	environment := Environment{
		EnvironmentV010: executionv010.EnvironmentV010{
			Command:      common.DeleteInstanceCommand,
			InstanceID:   "vm-01",
			InstanceName: "runner-a",
		},
	}

	_, err := environment.Run(context.Background(), provider)
	if err != nil {
		t.Fatal(err)
	}
	if concreteProvider.legacyDeleteCalled {
		t.Fatal("legacy DeleteInstance was called")
	}
	if concreteProvider.instanceID != "vm-01" || concreteProvider.instanceName != "runner-a" {
		t.Fatalf("delete identity = %q/%q, want vm-01/runner-a", concreteProvider.instanceID, concreteProvider.instanceName)
	}
}
