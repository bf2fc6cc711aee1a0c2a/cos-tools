# .adoc generator

Generates `.adoc` files from the `.json` definitions of connectors

## Dependecies

- Python 3

## Usage

```shell
python adoc-generator.py
```

Or (run `chmod +x adoc-generator.py` first)

```shell
./adoc-generator.py
```

A Linux x86_64 binary is provided in the `/bin` folder and can be used standalone. (run `chmod +x adoc-generator` first)

```shell
./adoc-generator
```

## Options

- **`-h`**, **`--help`**: Show help message and exit
  ```shell
  ./adoc-generator.py  -h
  ```
- **`-f`**, **`--jsonfiles`**: JSON files to convert (default: None). Either this or `-s` must be provided.
  ```shell
  ./adoc-generator.py -f file1.json file2.json /path-to-directory/file3.json
  ```
- **`-s`**, **`--source`**: Directory with JSON files to convert (default: None). Either this or `-f` must be provided.
  ```shell
  ./adoc-generator.py -s example-inpput-folder
  ```
- **`-d`**, **`--destination`**: Directory where the .adoc files will be created (default: current working directory)
  ```shell
  ./adoc-generator.py -s example-inpput-folder -d example-output-folder
  ```
- **`-i`** , **`--ignore-properties`**: Configurarion properties to be ignored (default: ['error_handler', 'processors'])
  ```shell
  ./adoc-generator.py -s example-inpput-folder -i error_handler data_shape another_property
  ```
- **`-r`** , **`--recursive`**: Used with -s scan provided directory recursively. (default: False)
  ```shell
  ./adoc-generator.py -s example-inpput-folder -r
  ```

## Generating standalone binary (Optional)

To generate the binary for the current OS architecture you will need pyinstaller and then run it pointing to the script. (Needs Python 3 and Pip)

```shell
pip install pyinstaller

pyinstaller adoc_generator.py -F
```
