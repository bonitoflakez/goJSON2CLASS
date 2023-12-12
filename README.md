# goJSON2CLASS

`goJSON2CLASS` is a utility that converts JSON schema to classes

## Usage

```txt
Usage: goJSON2CLASS -l <target-lang> -s <schema.json> -o <output.ext>

	-l >> choose a language.
		Example: `-l rust` (default: nil)

	-s >> path to file containing JSON schema. (default: schema.json)
		Example: `-s schema.json`

	-o >> path to output file with extension. (default: output.txt)
		Example: `-o output.rs`

	-p >> define public if supported by language (default: false)
		Example: `-p`
```

## Supported Inputs

goJSON2CLASS only supports JSON Schema

## Supported Languages

C, Go, C++, Java, Rust, TypeScript

_If your favorite language is missing- please generate an issue or implement it by yourself._

---

### [Installation](./docs/INSTALLATION.md)

---

### Usage Example

#### Sample Schema

```json
{
  "title": "Root",
  "properties": {
    "property1": {
      "type": "string"
    },
    "property2": {
      "type": "integer"
    },
    "property3": {
      "type": "object",
      "title": "Property3",
      "properties": {
        "nestedProperty1": {
          "type": "boolean"
        },
        "nestedProperty2": {
          "type": "array",
          "items": {
            "type": "string"
          }
        },
        "nestedProperty3": {
          "type": "string"
        }
      }
    }
  }
}
```

#### Generating C Code

```txt
>> goJSON2CLASS -l c -s schema.json -o output.c
Done!
```

Output

```c
#include <stdio.h>
#include <stdlib.h>
#include <stdbool.h>

#define PROPERTY3_NESTEDPROPERTY2_SIZE 50

typedef struct Root Root;
typedef struct Property3 Property3;

struct Property3 {
    bool nestedProperty1;
    char* nestedProperty2[PROPERTY3_NESTEDPROPERTY2_SIZE];
    char* nestedProperty3;
};
struct Root {
    char* property1;
    int property2;
    Property3 property3;
};
```

---

### [Examples](./docs/Example.md)

---

### LICENSE

[GPL](./LICENSE)
