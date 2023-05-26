# goJSON2CLASS

`goJSON2CLASS` is a utility that converts JSON schema to classes

## Supported Inputs

- [x] JSON Schema
- [ ] JSON
- [ ] JSON API URLs

## Target Languages

- [x] Rust
- [ ] JavaScript
- [ ] TypeScript
- [ ] Java
- [ ] Go
- [ ] C
- [ ] C++

_If your favorite language is missing- please generate an issue or implement it by yourself._

## Installation

```sh
# Clone this repo
>> git clone https://github.com/salientarc/goJSON2CLASS.git
```

```sh
# go inside the directory and run build command
>> cd goJSON2CLASS && go build .
```

```txt
>>  .\goJSON2CLASS -h
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

## Example

Sample JSON Schema

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

Command (without `-p` flag)

```sh
>> .\goJSON2CLASS.exe -l rust -s schema.json -o output.rs
Done!
```

Generated Rust code

```rs
use serde::{Serialize, Deserialize};

#[derive(Debug, Serialize, Deserialize)]
struct Root {
        #[serde(rename = "property1")]
        property1: String,
        #[serde(rename = "property2")]
        property2: i64,
        #[serde(rename = "property3")]
        property3: Property3,
}

#[derive(Debug, Serialize, Deserialize)]
struct Property3 {
        #[serde(rename = "nestedProperty1")]
        nestedProperty1: bool,
        #[serde(rename = "nestedProperty2")]
        nestedProperty2: Vec<String>,
        #[serde(rename = "nestedProperty3")]
        nestedProperty3: String,
}
```

Command (with `-p` flag)

```sh
>> .\goJSON2CLASS.exe -l rust -s .\schema.json -o output.rs -p
Done!
```

Generated Rust code

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

## LICENSE

[GPL](./LICENSE)
