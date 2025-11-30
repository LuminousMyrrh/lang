package eval

import "lang/internal/env"

func unwrapBuiltinValue(v any) any {
    if instEnv, ok := v.(*env.Env); ok {
        if instEnv.Parent != nil {
            pName := instEnv.Parent.Type
            if pName == "string" || pName == "int" || pName == "float" {
                if valueSym, ok := instEnv.Symbols["value"]; ok {
                    return valueSym.Value()
                }
            }
        }
    }
    return v
}
