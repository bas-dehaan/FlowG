# FlowG: Seamless Lab Data Integration

FlowG is a Go package that facilitates the connection between data files from lab equipment and GLIMS, a Laboratory Information System developed by CliniSys. It streamlines common FlowG-related tasks by automatically monitoring lab equipment upload folders and converting data into the FlowG file format.

For detailed documentation, visit [pkg.go.dev](https://pkg.go.dev/github.com/bas-dehaan/FlowG).

## Key Features

- **Folder Monitoring**: Automatically monitors lab equipment upload directories for new data files.
- **FlowG File Generation**: Converts lab data into FlowG-compatible files for GLIMS integration.
- **Error Handling and Logging**: Provides configurable directories for processed data, error handling, and log management.

## Installation

To use FlowG in your Go project, import the package in your source code:

```go
import "github.com/bas-dehaan/FlowG"
```

## Configuration

FlowG requires some initial configuration to set up folder paths and logging preferences. Use the `SetConfig` function to define these parameters.

### Configuration Parameters

- **glimsDir**: The directory where lab equipment uploads new data files.
- **processedDir**: The directory to store successfully processed data files for archival purposes.
- **errorDir**: The directory to move files that fail to process correctly.
- **logDir**: The directory for storing log files.
- **logLvl**: The log level to control the verbosity of log messages. Options include `DEBUG`, `INFO`, `WARNING`, `ERROR`, and `CRITICAL`.

### Watching for New Files

Once configuration is complete, use the `FileWatch()` function to monitor the `glimsDir` for new files. When a new file is detected, FlowG will call your custom processing function.

### Processing Function

Your processing function should load all data into a slice of `SampleStruct`. Then, you can call `GlimsOutput()` to generate the FlowG file

## Example implementation

```go
package main

import (
    "github.com/bas-dehaan/FlowG"
    "path/filepath"
)

func main() {
    err := FlowG.SetConfig("glimsDir", "path/to/upload/folder")
    err = FlowG.SetConfig("processedDir", "processed/data/archival")
    err = FlowG.SetConfig("errorDir", "data/failed/to/process")
    err = FlowG.SetConfig("logDir", "log/storage/folder")
    err = FlowG.SetConfig("logLvl", FlowG.WARNING)
    if err != nil {
        panic(err)
    }
    
    FlowG.FileWatch(yourProcessingFunction)
}

func yourProcessingFunction(filePath string) bool {
    var samples []FlowG.SampleStruct
    filename := filepath.Base(filePath)
    
    // Your file processing, specific to your input filetype...
    
    ok := FlowG.GlimsOutput(filename, samples)
    return ok
}
```

## Documentation

For more detailed usage instructions, please visit the [pkg.go.dev documentation](https://pkg.go.dev/github.com/bas-dehaan/FlowG).
