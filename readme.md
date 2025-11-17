# runner

A simple [Coveyor CI](https://conveyor.open.ug/) [driver](https://conveyor.open.ug/docs/concepts/drivers) that runs commands agains a codebase.

## Usage

### Resource Example

```json
{
  "name": "ubuntu-pipeline-6",
  "resource": "pipeline",
  "spec": {
    "image": "jimjuniorb/hello-node:1.2.3",
    "steps": [
      {
        "name": "print-working-directory",
        "command": "pwd"
      },
      {
        "name": "list-files",
        "command": "ls -l"
      },
      {
        "name": "show-os-info",
        "command": "cat /etc/os-release"
      }
    ]
  }
}
```