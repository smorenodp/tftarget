# Tftarget
## Description
Tftarget is a project that focus in allowing people to launch certain terraform commands with targets easier and faster that normal. Right now when using terraform you have to specify each and every one of the resources you want to target manually, which is fine until you have to detail more than ten targets or even every resource except one.

With Tftarget you can select the targets from a list dynamically and launch all the commands you need without having to copy paste the names.

## Installation
You only have to compile the go binary with go build.

## Usage
To launch tftarget you only need to invoke the binary in a terraform project directory or, using the dir flag to pass the directory.
```bash
tftarget
tftarget -dir <folder>
```
The tftarget automatically looks for the variable file in a vars/ folder found in said directory with the actual workspace loaded.

Once you launch the binary, the screen will show a spinner screen while it's loading the resources where there are changes

![Spinner Screen](https://github.com/smorenodp/tftarget/blob/main/images/spinner_screen.png)

Once it loads every resource the screen changes into a list containing the name of the resource and the type of change which can be update, replace, delete and create each with a different color.

![List Screen](https://github.com/smorenodp/tftarget/blob/main/images/list_screen.png)

In this screen you can use different keys to accomplish certain tasks:

* Space bar - To select / deselect the resources to target
* Forwardslash (/) - To filter through the resources
* Lowercase p - To launch a plan with the selected resources
* Lowercase a - To launch an apply with the selected resources
* Lowercase q - Quit the tftarget process
* Escape - Exit the filter
* Asterisk - Select all the resources

With the arrow keys you can navigate through the resources and change the page if there is more than one.

Once you have selected the resources and launch either a plan or apply, the screen changes to show the output of the command and ask input in case the command need certain parameters. Once the command finish you can return to the list pressing Enter or quit the process pressing q.

