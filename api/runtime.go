package api

import (
	"fmt"
	"github.com/tlarsendataguy/gobufrkit/bufr"
	"github.com/tlarsendataguy/gobufrkit/deserialize"
	"github.com/tlarsendataguy/gobufrkit/tdcfio"
)

type Config struct {
	DefinitionsPath string
	TablesPath      string

	// Only binary stream provides compressed data
	// in the format described by the BUFR Spec.
	InputType  tdcfio.InputType
	Compatible bool
	Verbose    bool
}

func (c *Config) toDeserializeConfig() *deserialize.Config {
	return &deserialize.Config{
		TablesPath: c.TablesPath,
		InputType:  c.InputType,
		Compatible: c.Compatible,
		Verbose:    c.Verbose,
	}
}

type Runtime struct {
	config   *Config
	scriptRt *ScriptRt
}

func NewRuntime(config *Config, pr tdcfio.PeekableReader) (*Runtime, error) {

	factory := deserialize.NewDefaultFactory(config.toDeserializeConfig(), pr)

	scriptRt := NewScriptRt(config.DefinitionsPath, factory)
	if err := scriptRt.Initialize(); err != nil {
		return nil, fmt.Errorf("cannot initialise script runtime: %w", err)
	}

	return &Runtime{
		config:   config,
		scriptRt: scriptRt,
	}, nil
}

func (rt *Runtime) Run() (*bufr.Message, error) {
	return rt.scriptRt.RunDeserializer()
}
