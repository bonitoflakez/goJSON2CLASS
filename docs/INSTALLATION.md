# Installation

```sh
# Clone this repo
>> git clone https://github.com/bonitoflakez/goJSON2CLASS.git
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
