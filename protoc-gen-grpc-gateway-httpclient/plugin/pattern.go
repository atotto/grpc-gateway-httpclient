package plugin

import (
	"fmt"
	"strings"

	"github.com/gogo/protobuf/protoc-gen-gogo/generator"
	"github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway/httprule"
	"github.com/grpc-ecosystem/grpc-gateway/utilities"
)

type op struct {
	code    utilities.OpCode
	operand int
}

func MustParsePattern(pattern string) string {
	p, _ := ParsePattern(pattern)
	return p
}

func ParsePattern(pattern string) (string, error) {
	parsed, err := httprule.Parse(pattern)
	if err != nil {
		return "", err
	}
	tmpl := parsed.Compile()

	if tmpl.Version != 1 {
		return "", fmt.Errorf("unsupported version: %d", tmpl.Version)
	}

	ops := tmpl.OpCodes
	pool := tmpl.Pool

	l := len(ops)
	if l%2 != 0 {
		return "", fmt.Errorf("odd number of ops codes: %d", l)
	}

	var request strings.Builder
	request.WriteString(`"`)

	var (
		typedOps        []op
		stack, maxstack int
		tailLen         int
		pushMSeen       bool
		vars            []string
	)
	for i := 0; i < l; i += 2 {
		op := op{code: utilities.OpCode(ops[i]), operand: ops[i+1]}
		switch op.code {
		case utilities.OpNop:
			continue
		case utilities.OpPush:
			if pushMSeen {
				tailLen++
			}
			stack++
		case utilities.OpPushM:
			if pushMSeen {
				return "", fmt.Errorf("pushM appears twice")
			}
			pushMSeen = true
			stack++
		case utilities.OpLitPush:
			if op.operand < 0 || len(pool) <= op.operand {
				return "", fmt.Errorf("negative literal index: %d", op.operand)
			}

			request.WriteString("/")
			request.WriteString(pool[op.operand])

			if pushMSeen {
				tailLen++
			}
			stack++
		case utilities.OpConcatN:
			if op.operand <= 0 {
				return "", fmt.Errorf("negative concat size: %d", op.operand)
			}
			stack -= op.operand
			if stack < 0 {
				return "", fmt.Errorf("stack underflow")
			}
			stack++
		case utilities.OpCapture:
			if op.operand < 0 || len(pool) <= op.operand {
				return "", fmt.Errorf("variable name index out of bound: %d", op.operand)
			}
			v := pool[op.operand]
			op.operand = len(vars)
			vars = append(vars, v)
			stack--
			if stack < 0 {
				return "", fmt.Errorf("stack underflow")
			}

			request.WriteString(`/"+ req.`)
			request.WriteString(generator.CamelCase(v))
			request.WriteString(`+ "`)

		default:
			return "", fmt.Errorf("invalid opcode: %d", op.code)
		}

		if maxstack < stack {
			maxstack = stack
		}
		typedOps = append(typedOps, op)
	}

	request.WriteString(`"`)
	return request.String(), nil
}
