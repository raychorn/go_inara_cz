package enums

import (
	"bytes"
	"fmt"
)

// Tritium=10269
// AgronomicTreatment=10268

// Alias hide the real type of the enum
// and users can use it to define the var for accepting enum
type Alias = map[string]int

type list struct { 
    Tritium Alias
    AgronomicTreatment Alias
}

// Enum for public use
var Enum = &list{ 
    Tritium: map[string]int{"Tritium":10269},
    AgronomicTreatment: map[string]int{"AgronomicTreatment": 10268},
}

// String ...
func String(m map[string]int) string {
    b := new(bytes.Buffer)
    for key, value := range m {
        fmt.Fprintf(b, "{\"%s\":%d}", key, value)
    }
    return b.String()
}
