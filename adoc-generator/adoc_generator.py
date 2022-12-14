#!/usr/bin/env python3
"""Module providing CLI to convert json files describing connectors to adoc files."""

import argparse
import json
from pprint import pprint
import os
import sys
import textwrap

"""Auxiliary function to output errors to stderr."""


def eprint(*args, **kwargs):
    print(*args, file=sys.stderr, **kwargs)


def print_connector_header(connector_id, connector_name, connector_description):
    """Creates the main description of a connector."""
    connector_description = "{Empty}" if connector_description is None else connector_description
    return textwrap.dedent(
        f"""
        [id='{connector_id}']
        = {connector_name}
        ifdef::context[:parent-context: {{context}}]
        :context: connectors-reference
        :imagesdir: _images

        {connector_description}

        == Configuration properties

        The following table describes the configuration properties for the {connector_name} Connector.

        NOTE: Fields marked with an asterisk (*) are mandatory.

        """
    )


def print_property_table_header(*headers):
    """Creates a table header with one column for each value in headers"""
    table_header = """[width="100%",cols="2,^2,3,^2,^2,^3",options="header"]\n|===\n"""
    for header in headers:
        table_header += f"|{header}"
    return f"{table_header}\n"


def print_property_table_row(*rows):
    """Creates a table row with one column for each value in rows"""
    table_row = ""
    for row in rows:
        table_row += f"|{row}"
    return f"{table_row}\n"


def print_property_table_footer():
    """Creates a table footer"""
    return "|===\n"


def validate_args(
    file_path_args: list,
    source_dir_path_arg: str,
    destination_dir_path_arg: str,
    parser: argparse.ArgumentParser,
):
    """Validates the arguments passed on to the cli

    Args:
        file_path_args (): path to individual files to be processed
        source_dir_path_arg (str): path to a directory with files to be processed
        destination_dir_path_arg (str): path to a directory to output the generated adoc files
        parser (argparse.ArgumentParser): the argument parser instance
    """
    if file_path_args is None and source_dir_path_arg is None:
        eprint("error: -f, -s (or both) options must be provided!\n")
        parser.print_help()
        sys.exit(1)
    if file_path_args is not None:
        for file_path in file_path_args:
            if file_path is None or not os.path.exists(file_path):
                eprint(f'error: "{file_path}" does not exist or is not a file')
                sys.exit(2)
    if source_dir_path_arg is not None:
        if not os.path.exists(source_dir_path_arg) or not os.path.isdir(source_dir_path_arg):
            eprint(f'error: "{source_dir_path_arg}" does not exist or is not a directory')
            sys.exit(3)
    if destination_dir_path_arg is not None:
        if not os.path.exists(destination_dir_path_arg) or not os.path.isdir(destination_dir_path_arg):
            eprint(f'error: "{destination_dir_path_arg}" does not exist or is not a directory')
            sys.exit(4)


def get_composite_property(connector_schema, attribute_name: str, property_key: str, data_shape_key_or_value: str) -> str:
    """Retrieves the contents of a property which can be inside multiple ``oneOf`` objects
    There is a special case to get the description of the data_shape property. It will be inside the $defs object
    The format for camel and debezium connectors differs a little so there is also a special case for that.

    Args:
        connector_schema (): json object which represents the connector_json_definition["connector_type"]["schema"]
        attribute_name (str): attribute to get the value of
        property_key (str): which property to get the attribute from
    """
    property_content = connector_schema["properties"][property_key]
    if property_key == "data_shape" and attribute_name == "description":
        camel_data_shape_description = ""
        if property_key in connector_schema["$defs"]:  # camel connector
            if "consumes" in connector_schema["$defs"][property_key]:
                camel_data_shape_description = connector_schema["$defs"][
                    property_key]["consumes"]["properties"]["format"]["description"]
            if "produces" in connector_schema["$defs"][property_key]:
                camel_data_shape_description = connector_schema["$defs"][
                    property_key]["produces"]["properties"]["format"]["description"]
            return camel_data_shape_description
        # debezium connector
        elif property_key in connector_schema["properties"]:
            debezium_data_shape_description = f" {connector_schema['properties'][property_key]['properties'][data_shape_key_or_value]['description']}"
            return debezium_data_shape_description

    if attribute_name in property_content:
        return property_content[attribute_name]
    if "oneOf" in property_content:
        multiple_values = ""
        for each in property_content["oneOf"]:
            if attribute_name in each:
                multiple_values += each[attribute_name] if multiple_values == "" else f" / {each[attribute_name]}"
        return multiple_values
    return ""


def get_normal_property(property_content, attribute_name: str) -> str:
    """Retrieves the contents of a property"""
    return property_content[attribute_name] if attribute_name in property_content else ""


def get_required_property(connector_schema, attribute_name: str, property_key: str, data_shape_key_or_value: str, required_properties: list) -> str:
    """
    Retrieves the contents of a property which can be required/mandatory. If it is a ``*`` is appended to the returned value.
    There is a special case the data_shape property which normally does not have a "title" attribute. the data_shape property also
    differs a little between camel and debezium connectors
    """
    property_content = connector_schema["properties"][property_key]

    if property_key == "data_shape" and attribute_name == "title":
        if property_key in connector_schema['$defs']: # camel connector
            return_value = "Data Shape"
        elif property_key in connector_schema['properties']: # debezium connector
            return_value = connector_schema['properties'][property_key]['properties'][data_shape_key_or_value]['title']
    else:
        return_value = property_content[attribute_name] if attribute_name in property_content else ""

    return_value += "*" if property_key in required_properties else ""
    return return_value


def convert_to_adoc_from_json_file(json_file_path: str, destination_dir_path: str, ignored_properties: list):
    """Convert specified json file that describes a connector to an adoc file with information about the connector

    Args:
        json_file_path (str): path to json file to convert
        destination_dir_path (str): path to destination directory were the adoc files will be created
    """
    print(f"\nProcessing file: {json_file_path} ...")
    with open(json_file_path, encoding="utf-8") as read_file:
        output_adoc_file_path = os.path.join(
            destination_dir_path,
            os.path.splitext(os.path.basename(json_file_path))[0] + ".adoc",
        )
        json_data = json.load(read_file)
        if "connector_type" not in json_data:
            eprint(f'warn: connector_type property not found maybe "{json_file_path}" is not a valid connector description file')
            return

        properties_keys = json_data["connector_type"]["schema"]["properties"].keys()
        required_properties = json_data["connector_type"]["schema"]["required"]
        with open(output_adoc_file_path, "w", encoding="utf-8") as out_file:
            print(f"Creating file:   {output_adoc_file_path} ...")

            out_file.write(
                print_connector_header(
                    connector_id=json_data["connector_type"]["id"],
                    connector_name=json_data["connector_type"]["name"],
                    connector_description=json_data["connector_type"]["description"],
                )
            )
            out_file.write(print_property_table_header(
                "Name", "Property", "Description", "Type", "Default", "Example"))

            for property_key in properties_keys:
                if property_key not in ignored_properties:
                    property_content = json_data["connector_type"]["schema"]["properties"][property_key]
                    connector_schema = json_data["connector_type"]["schema"]
                    if property_key == "data_shape" and "key" in property_content["properties"]:
                        table_row = print_property_table_row(
                            "*{Empty}" + get_required_property(connector_schema, "title", property_key, "key", required_properties) + "*", "`" + property_key + "`",
                            get_composite_property(connector_schema, "description", property_key, "key"),
                            get_composite_property(connector_schema, "type", property_key, "key"),
                            get_normal_property(property_content, "default"),
                            get_normal_property(property_content, "example"),
                        )
                        table_row += print_property_table_row(
                            "*{Empty}" + get_required_property(connector_schema, "title", property_key, "value", required_properties) + "*", "`" + property_key + "`",
                            get_composite_property(connector_schema, "description", property_key, "value"),
                            get_composite_property(connector_schema, "type", property_key, "value"),
                            get_normal_property(property_content, "default"),
                            get_normal_property(property_content, "example"),
                        )

                    else:
                        table_row = print_property_table_row(
                            "*{Empty}" + get_required_property(
                                connector_schema, "title", property_key, "", required_properties) + "*", "`" + property_key + "`",
                            get_composite_property(connector_schema, "description", property_key, ""),
                            get_composite_property(connector_schema, "type", property_key, ""),
                            get_normal_property(property_content, "default"),
                            get_normal_property(property_content, "example"),
                        )

                    out_file.write(table_row)

            out_file.write(print_property_table_footer())

def main():
    """Main function"""
    parser = argparse.ArgumentParser(
        description="Parse JSON files that represents a connector and convert them to a .adoc files",
        formatter_class=argparse.ArgumentDefaultsHelpFormatter,
        add_help=False,
    )
    required_group = parser.add_argument_group("required arguments")
    required_group.add_argument(
        "-f",
        "--jsonfiles",
        help="List of JSON files to convert. \
        Either this or -s must be provided",
        nargs="*",
    )
    required_group.add_argument(
        "-s",
        "--source",
        help="Directory with JSON files to convert. \
        Either this or -f must be provided",
    )

    optional_group = parser.add_argument_group("optional arguments")
    optional_group.add_argument(
        "-h", "--help",
        action="help",
        help="Show this help message and exit")
    optional_group.add_argument(
        "-r",
        "--recursive",
        action="store_true",
        help="Used with -s, \
        --source will scan for JSON files recursively starting with the directory informed.",
    )
    optional_group.add_argument(
        "-d",
        "--destination",
        help="Directory where the .adoc files will be created.",
        default=".",
    )
    optional_group.add_argument(
        "-i",
        "--ignore-properties",
        help="Configurarion properties to be ignored.",
        nargs="*",
        default=["error_handler", "processors"],
    )

    args = parser.parse_args()

    file_path_args = vars(args)["jsonfiles"]
    source_dir_path_arg = vars(args)["source"]
    destination_dir_path_arg = vars(args)["destination"]
    ignored_properties = vars(args)["ignore_properties"]
    recursive_source_dir_scan = vars(args)["recursive"]

    validate_args(file_path_args, source_dir_path_arg, destination_dir_path_arg, parser)

    destination_dir_path_arg = os.path.join(destination_dir_path_arg, "")

    if file_path_args is not None:
        for file_path in file_path_args:
            convert_to_adoc_from_json_file(file_path, destination_dir_path_arg, ignored_properties)

    if source_dir_path_arg is not None:
        source_dir_path_arg = os.path.join(source_dir_path_arg, "")
        if recursive_source_dir_scan:
            for directory in [tuple[0] for tuple in os.walk(source_dir_path_arg)]:
                print(f"\nProcessing directory: {directory} ...")
                for file_name in os.listdir(directory):
                    file_path = os.path.join(directory, file_name)
                    if os.path.isfile(file_path) and file_name.lower().endswith((".json")):
                        convert_to_adoc_from_json_file(file_path, destination_dir_path_arg, ignored_properties)
        else:
            print(f"\nProcessing directory: {source_dir_path_arg} ...")
            for file_name in os.listdir(source_dir_path_arg):
                file_path = os.path.join(source_dir_path_arg, file_name)
                if os.path.isfile(file_path) and file_name.lower().endswith((".json")):
                    convert_to_adoc_from_json_file(file_path, destination_dir_path_arg, ignored_properties)

    print("\nDone!")


if __name__ == "__main__":
    main()
