
# Helm Values Manager

Helm Values Manager is a Go-based tool for merging multiple Helm values files, tracking overrides, and displaying the final merged values. It helps you identify how configuration values have been overridden between different environments and keeps track of base, override, and final values for easy debugging.

## Features

- **Track Base and Override Values**: Clearly display the base values from the initial file and how they are overridden in subsequent files.
- **Final Merged Values**: After processing all override files, the tool displays the final merged values.
- **List and Map Support**: Properly handles YAML lists and maps, ensuring that nested structures are merged and reported correctly.
- **Base Value Tracking**: Displays base values when present and logs them as `null` when they don’t exist in the base file.
- **Clear Logging Format**: Displays each key, along with the base value, all override values, and the final value.

## Installation

1. Clone this repository:

   \`\`\`bash
   git clone https://github.com/mwojciga/helm-values-manager.git
   cd helm-values-manager
   \`\`\`

2. Install dependencies and compile the Go program:

   \`\`\`bash
   go mod tidy
   go build
   \`\`\`

3. (Optional) Run the software without building:

   \`\`\`bash
   go run main.go <base-values-file> <override-values-file> <additional-override-files>
   \`\`\`

## Usage

To run the program, provide at least one base file and one or more override files. The base file is the first argument, and the subsequent files are the override files. The tool will process these files in order, log how values are overridden, and print the final merged values.

### Example

Assume you have the following Helm values files:

1. **Base file (`values.yaml`)**:

   \`\`\`yaml
   thanos:
     store:
       resources:
         requests:
           cpu: 500m
   \`\`\`

2. **Override file 1 (`values2.yaml`)**:

   \`\`\`yaml
   thanos:
     store:
       resources:
         requests:
           cpu: 1
   \`\`\`

3. **Override file 2 (`values3.yaml`)**:

   \`\`\`yaml
   thanos:
     store:
       resources:
         requests:
           memory: 6Gi
   \`\`\`

Run the program with:

\`\`\`bash
go run main.go values.yaml values2.yaml values3.yaml
\`\`\`

### Output

The tool will output something like:

\`\`\`
--- Override and Final Values Log ---
Key .thanos.store.resources.requests.cpu:
  Base value: 500m
  Override value: 1 (values2.yaml)
  Final value: 1
---
Key .thanos.store.resources.requests.memory:
  Base value: null
  Override value: 6Gi (values3.yaml)
  Final value: 6Gi
---
\`\`\`

### Options and Commands

- **Base Value Tracking**: The first argument is considered the base values file, and all subsequent files are treated as overrides.
- **Override Reporting**: The program logs the base value if it exists, followed by each override value, and finally the resulting merged value.

### Supported YAML Structures

- **Scalars**: Simple key-value pairs like `cpu: 500m`.
- **Lists**: Supports YAML lists like:
  \`\`\`yaml
  items:
    - item1
    - item2
  \`\`\`
- **Maps**: Nested YAML structures like:
  \`\`\`yaml
  resources:
    requests:
      cpu: 500m
  \`\`\`

## Error Handling

The tool will report an error if:
- A file cannot be read or parsed.
- The required arguments (base and at least one override file) are not provided.

## Development

To make changes, clone the repository and modify the Go code. Use \`go run main.go\` to test the code. Contributions are welcome via pull requests.

