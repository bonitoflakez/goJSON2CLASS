# Example

## Sample JSON Schema

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

## Generating C Code

```txt
>> ./goJSON2CLASS -l c -s schema.json -o output.rs
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

If we use `-p` flag here it shows "Public is not supported for `<target-lang>`"

```txt
>> ./goJSON2CLASS -l c -s schema.json -o output.c -p
Public is not supported for c
Choosing default settings
Done!
```

## Generating Rust Code

```sh
>> ./goJSON2CLASS -l rust -s schema.json -o output.c
```

Output

```rs
use serde::{Serialize, Deserialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct Root {
        #[serde(rename = "property1")]
        property1: String,
        #[serde(rename = "property2")]
        property2: i64,
        #[serde(rename = "property3")]
        property3: Property3,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct Property3 {
        #[serde(rename = "nestedProperty1")]
        nestedProperty1: bool,
        #[serde(rename = "nestedProperty2")]
        nestedProperty2: Vec<String>,
        #[serde(rename = "nestedProperty3")]
        nestedProperty3: String,
}
```

Generating Rust code with `-p` (Public flag)

```sh
>> ./goJSON2CLASS -l rust -s schema.json -o output.rs -p
Done!
```

Output

```rs
use serde::{Serialize, Deserialize};

#[derive(Debug, Serialize, Deserialize)]
pub struct Root {
        #[serde(rename = "property1")]
        pub property1: String,
        #[serde(rename = "property2")]
        pub property2: i64,
        #[serde(rename = "property3")]
        pub property3: Property3,
}

#[derive(Debug, Serialize, Deserialize)]
pub struct Property3 {
        #[serde(rename = "nestedProperty1")]
        pub nestedProperty1: bool,
        #[serde(rename = "nestedProperty2")]
        pub nestedProperty2: Vec<String>,
        #[serde(rename = "nestedProperty3")]
        pub nestedProperty3: String,
}
```
