package compilers

import (
	"errors"
	"fmt"

	"github.com/steve-care-software/products/pangolin/domain/interpreters"
	"github.com/steve-care-software/products/pangolin/domain/linkers"
	"github.com/steve-care-software/products/pangolin/domain/middle"
	"github.com/steve-care-software/products/pangolin/domain/middle/variables/variable/value/computable"
	"github.com/steve-care-software/products/pangolin/domain/parsers"
)

type application struct {
	computableBuilder  computable.Builder
	interpreterBuilder interpreters.Builder
	parser             parsers.Parser
	middleAdapter      middle.Adapter
	programBuilder     linkers.ProgramBuilder
}

func createApplication(
	computableBuilder computable.Builder,
	interpreterBuilder interpreters.Builder,
	middleAdapter middle.Adapter,
	programBuilder linkers.ProgramBuilder,
) Application {
	out := application{
		computableBuilder:  computableBuilder,
		interpreterBuilder: interpreterBuilder,
		middleAdapter:      middleAdapter,
		programBuilder:     programBuilder,
	}

	return &out
}

// Execute executes the compiler application
func (app *application) Execute(script linkers.Script) (middle.Program, error) {
	langRef := script.Language()
	lang := langRef.Language()
	langApp := lang.Application()
	inVariable := langRef.Input()
	outVariable := langRef.Output()
	code := script.Code()

	codeValue, err := app.computableBuilder.Create().WithString(code).Now()
	if err != nil {
		return nil, err
	}

	input := map[string]computable.Value{
		inVariable: codeValue,
	}

	linkedProgram, err := app.programBuilder.Create().WithApplication(langApp).Now()
	if err != nil {
		return nil, err
	}

	interpreter, err := app.interpreterBuilder.Create().WithProgram(linkedProgram).Now()
	if err != nil {
		return nil, err
	}

	if !interpreter.IsApplication() {
		return nil, errors.New("the interpreter was expected to be an application interpreter")
	}

	stackFrame, err := interpreter.Application().Execute(input)
	if err != nil {
		return nil, err
	}

	computedCodeValue, err := stackFrame.Current().Fetch(outVariable)
	if err != nil {
		return nil, err
	}

	if !computedCodeValue.IsString() {
		str := fmt.Sprintf("the output variable (%s) was expected to contain code and therefore be a string", outVariable)
		return nil, errors.New(str)
	}

	pOutputCode := computedCodeValue.String()
	programIns, err := app.parser.ExecuteScript(*pOutputCode)
	if err != nil {
		return nil, err
	}

	if parsedProgram, ok := programIns.(parsers.Program); ok {
		return app.middleAdapter.ToProgram(parsedProgram)
	}

	return nil, errors.New("the parsed compiled output was expected to contain a parsers.Program instance")
}
