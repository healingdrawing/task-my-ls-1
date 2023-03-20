<h3 align="center">my-ls-1</h3>



<p align="center"> My-ls is a project, consisting on creating your own ls command.
    <br> 
</p>

## 📝 Table of Contents

- [About](#about)
- [Usage](#usage)
- [Deployment](#deployment)
- [Built Using](#built_using)
- [Authors](#authors)

## 🧐 About <a name = "about"></a>

my-ls-1 is a command-line tool written in Go that allows users to list the files and folders in a directory. The tool is designed to behave similarly to the original ls command with a few variations, including incorporating the following flags:

-l : Display the files and directories in a long format, including file permissions, ownership, size, and modification date.
-R : Recursively list subdirectories.
-a : Include hidden files and directories in the list.
-r : Reverse the order of the list.
-t : Sort the list by modification time.



## 🎈 Usage <a name="usage"></a>

To install my-ls-1, follow these steps:

Clone the repository: git clone https://github.com/<Mijan>/my-ls-1.git
Change to the directory where you cloned the repository: cd my-ls
Build the tool: go build

### Installing

To install my-ls-1, follow these steps:

Clone the repository: git clone https://github.com/<Mijan>/my-ls-1.git
Change to the directory where you cloned the repository: cd my-ls
Build the tool: go build

## 🔧 Running the tests <a name = "tests"></a>

Once you have built the tool, you can use it by running the following command:

./my-ls-1 [flags] [path]

The flags are optional and can be any combination of -l, -R, -a, -r, and -t. The path is also optional and can be any valid directory path. If no path is specified, the tool will default to the current directory.


### And coding style tests

Here are some examples of how to use my-ls-1:

List the files and directories in the current directory:
./my-ls-1

List the files and directories in a specific directory:
./my-ls-1 /path/to/directory

List the files and directories in a specific directory in long format:
./my-ls-1 -l /path/to/directory

List the files and directories in a specific directory recursively:
./my-ls-1 -R /path/to/directory

List all files and directories in a specific directory, including hidden files and directories:
./my-ls-1 -a /path/to/directory

Sort the files and directories in a specific directory by modification time:
./my-ls-1 -t /path/to/directory

Sort the files and directories in a specific directory in reverse order:
./my-ls-1 -r /path/to/directory



## ✍️ Authors <a name = "authors"></a>

-[@lenivaya10003]https://01.gritlab.ax/git/lenivaya10003
-[@MariaSagulin]https://01.gritlab.ax/git/MariaSagulin
-[@Mijan]https://01.gritlab.ax/git/Mijan