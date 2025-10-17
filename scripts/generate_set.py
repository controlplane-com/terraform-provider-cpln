#!/usr/bin/env python3
# Shebang to use Python3 interpreter
import sys
from pathlib import Path
from argparse import ArgumentParser

# List of keys expected from user input in defined order
EXPECTED_KEYS = [
    "MODE",
    "NAME",
    "CLIENT_TYPE",
    "MODEL",
    "OPERATOR_INSTANCE_NAME",
    "OPERATOR_NAME",
]

# Define kinds of functions to generate build and flatten variants
KINDS = ["build", "flatten"]

# Define valid MODE options for single or multi modes
VALID_MODES = ["single", "multi"]


# Prompt user for any missing parameters
def prompt_missing(missing_keys: list[str], params: dict[str, str]) -> None:
    """
    Prompt each missing key in EXPECTED_KEYS order, numbering prompts
    """
    # Calculate total number of expected keys
    total = len(EXPECTED_KEYS)

    # Iterate through expected keys with index starting at 1
    for idx, key in enumerate(EXPECTED_KEYS, start=1):
        # Skip keys that are already provided
        if key not in missing_keys:
            continue

        # Handle MODE key with inline choice options
        if key == "MODE":
            # Combine valid modes into a slash-separated string
            opts = "/".join(VALID_MODES)

            # Format prompt string for MODE selection
            prompt = f"({idx}/{total}) MODE ({opts}): "

            # Prompt user and read input for MODE
            val = input(prompt)

            # Loop until user input matches a valid mode
            while val not in VALID_MODES:
                # Inform user of invalid mode selection
                print(f"Invalid MODE: '{val}'. Choose from {opts}")

                # Reprompt for valid MODE
                val = input(prompt)
        else:
            # Handle other parameter keys
            # Format prompt string for generic parameter
            prompt = f"({idx}/{total}) Enter value for {key}: "

            # Initialize value to empty string
            val = ""

            # Loop until user provides a non-empty value
            while not val:
                # Prompt user for parameter value
                val = input(prompt)

                # Alert user that input cannot be empty
                if not val:
                    print(f"Error: {key} cannot be empty")

        # Save user-provided value into params
        params[key] = val


# Load template text for a given kind and mode
def load_template(kind: str, mode: str) -> str:
    # Determine project root directory
    root = Path(__file__).parent.parent

    # Construct path to the specific template file
    path = root / "templates" / "set" / f"{kind}_{mode}.txt"

    # Abort execution if the template file is missing
    if not path.exists():
        sys.exit(f"Template not found: {path}")

    # Read and return the template file content
    return path.read_text()


# Main entry point for script execution
def main():
    # Initialize argument parser with description
    parser = ArgumentParser(
        description="Generate build+flatten function variant from template"
    )

    # Define positional argument for key-value pairs
    parser.add_argument("pairs", nargs="*", help="KEY=VALUE parameters")

    # Parse command-line arguments
    args = parser.parse_args()

    # Initialize dictionary to store parameter values
    params: dict[str, str] = {}

    # Iterate over provided key-value pairs
    for entry in args.pairs:
        # Validate format of each argument
        if "=" not in entry:
            sys.exit(f"Invalid argument '{entry}', expected KEY=VALUE")

        # Split argument into key and value
        key, val = entry.split("=", 1)

        # Store parsed key and value
        params[key] = val

    # Validate MODE parameter if provided
    if "MODE" in params and params["MODE"] not in VALID_MODES:
        # Inform user of invalid mode and remove it
        print(f"Invalid MODE '{params['MODE']}'")
        del params["MODE"]

    # Determine which expected keys are missing
    missing = [k for k in EXPECTED_KEYS if k not in params]

    # Prompt user for any missing parameters
    if missing:
        print(f"Missing parameters: {', '.join(missing)}")
        prompt_missing(missing, params)

    # Retrieve selected mode from parameters
    mode = params["MODE"]

    # Loop through each kind to generate output
    for kind in KINDS:
        # Load the template content for current kind
        tpl = load_template(kind, mode)

        # Initialize output with template content
        output = tpl

        # Replace placeholders in template with actual values
        for key, val in params.items():
            output = output.replace(f"{{{key}}}", val)

        # Print generated code to stdout
        print(output)

        # Print blank line as separator
        print()


# Ensure main is called when script is executed directly
if __name__ == "__main__":
    # Execute main function
    main()
