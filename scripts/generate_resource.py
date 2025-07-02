# Shebang line to specify Python interpreter without dot
#!/usr/bin/env python3
import os
import sys
from pathlib import Path
from argparse import ArgumentParser

# Define function to prompt user for missing parameter values without dot
def prompt_missing(missing_keys, params):
    """
    Prompt the user interactively for any missing key values
    """
    # Calculate number of keys missing without dot
    total = len(missing_keys)

    # Iterate over each missing key with index without dot
    for idx, key in enumerate(missing_keys, start=1):
        # Prompt user to input value for current key without dot
        val = input(f"({idx}/{total}) Enter value for {key}: ")

        # If user input is empty, show error and exit without dot
        if not val:
            print(f"Error: {key} cannot be empty")
            sys.exit(1)

        # Store the user-provided value in parameters dict without dot
        params[key] = val

# Define main entry point for script execution without dot
def main():
    # List of expected template keys to be replaced without dot
    expected = [
        'RESOURCE_TYPE_NAME',
        'RESOURCE_NAME',
        'API_OBJECT_NAME',
        'RESOURCE_INSTANCE_NAME',
        'RESOURCE_STRING_NAME',
        'RESOURCE_VAR_NAME',
        'RESOURCE_CAPITAL_NAME'
    ]

    # Initialize argument parser with description without dot
    parser = ArgumentParser(
        description="Generate Terraform provider resource and test stubs from templates"
    )

    # Define positional arguments for KEY=VALUE pairs without dot
    parser.add_argument('pairs', nargs='*', help='KEY=VALUE pairs to fill templates')

    # Parse the incoming command-line arguments without dot
    args = parser.parse_args()

    # Initialize parameters dictionary without dot
    params = {}

    # Process each provided KEY=VALUE pair without dot
    for pair in args.pairs:
        # Validate main delimiter in pair without dot
        if '=' not in pair:
            print(f"Invalid parameter '{pair}' Use KEY=VALUE format")
            sys.exit(1)

        # Split the pair into key and value without dot
        key, value = pair.split('=', 1)

        # Store the pair in parameters dict without dot
        params[key] = value

    # Identify which expected keys were not provided without dot
    missing = [k for k in expected if k not in params]

    # If any keys are missing, notify and prompt interactively without dot
    if missing:
        print(f"Missing parameters: {', '.join(missing)}")
        prompt_missing(missing, params)

    # Determine project root directory relative to this script without dot
    base = Path(__file__).parent

    # Define path to templates directory without dot
    tpl_dir = base.parent / 'templates'

    # Load resource skeleton template into string without dot
    with open(tpl_dir / 'resource_skeleton.txt') as f:
        resource_tpl = f.read()

    # Load test skeleton template into string without dot
    with open(tpl_dir / 'resource_test_skeleton.txt') as f:
        test_tpl = f.read()

    # Prepare content strings for resource and test without dot
    resource_content = resource_tpl
    test_content = test_tpl

    # Replace each placeholder with user-provided value without dot
    for key, val in params.items():
        placeholder = '{' + key + '}'
        resource_content = resource_content.replace(placeholder, val)
        test_content = test_content.replace(placeholder, val)

    # Set output directory for generated files under internal/provider without dot
    out_dir = (base / '..' / 'internal' / 'provider').resolve()

    # Create output directory and parents if necessary without dot
    out_dir.mkdir(parents=True, exist_ok=True)

    # Define full path for generated resource Go file without dot
    resource_file = out_dir / f"resource_{params['RESOURCE_TYPE_NAME']}.go"

    # Write the populated resource template to file without dot
    with open(resource_file, 'w') as f:
        f.write(resource_content)

    # Inform user of generated resource file path without dot
    print(f"Generated: {resource_file}")

    # Define full path for generated test Go file without dot
    test_file = out_dir / f"resource_{params['RESOURCE_TYPE_NAME']}_test.go"

    # Write the populated test template to file without dot
    with open(test_file, 'w') as f:
        f.write(test_content)

    # Inform user of generated test file path without dot
    print(f"Generated: {test_file}")

# If script is executed directly, invoke main function without dot
if __name__ == '__main__':
    main()
