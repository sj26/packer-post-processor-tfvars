package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/mitchellh/packer/common"
	"github.com/mitchellh/packer/helper/config"
	"github.com/mitchellh/packer/packer"
	"github.com/mitchellh/packer/template/interpolate"
)

type Config struct {
	common.PackerConfig `mapstructure:",squash"`

	OutputPath string `mapstructure:"output"`

	ctx interpolate.Context
}

type PostProcessor struct {
	config Config
}

func (p *PostProcessor) Configure(raws ...interface{}) error {
	err := config.Decode(&p.config, &config.DecodeOpts{
		Interpolate:        true,
		InterpolateContext: &p.config.ctx,
		InterpolateFilter: &interpolate.RenderFilter{
			Exclude: []string{},
		},
	}, raws...)
	if err != nil {
		return err
	}

	if p.config.OutputPath == "" {
		p.config.OutputPath = "packer.tfvars"
	}

	return nil
}

func (p *PostProcessor) PostProcess(ui packer.Ui, artifact packer.Artifact) (packer.Artifact, bool, error) {
	switch artifact.BuilderId() {
	case "mitchellh.amazonebs", "mitchellh.amazon.chroot", "mitchellh.amazon.instance":
		break
	default:
		err := fmt.Errorf("Unknown artifact type: %s\nCan only export from Amazon builders", artifact.BuilderId())
		return artifact, true, err
	}

	name := p.config.PackerBuildName
	amis := parseAmis(artifact.Id())
	vars := make(map[string]string)
	for region, id := range amis {
		vars["packer-"+region+"-"+name] = id
	}

	output, err := os.Create(p.config.OutputPath)
	if err != nil {
		return artifact, true, err
	}
	defer output.Close()

	enc := json.NewEncoder(output)
	if err := enc.Encode(&vars); err != nil {
		return artifact, true, err
	}

	return artifact, true, nil
}

func parseAmis(artifactId string) map[string]string {
	amis := make(map[string]string)

	for _, ami := range strings.Split(artifactId, ",") {
		pair := strings.SplitN(ami, ":", 2)
		region, id := pair[0], pair[1]

		amis[region] = id
	}

	return amis
}
