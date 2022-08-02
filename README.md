# docshtest
A small utility for testing shell sessions from Markdown documents. Like Python
[doctests](https://docs.python.org/3/library/doctest.html) expect for the Shell.

If you have the following Markdown document in a file called `document.md`

````markdown
Example document.

```console
$ mkdir test && cd test
$ echo "this block will be run" > myfile
$ cat myfile
this block will be run
```

```python
print("code fences with other tags are ignored")
```
````

Running `docshtest --run-highlighted-code-fences console` will extract all the
code fences with "console" as the highlighting language and execute them. Any
output lines will be compared against the actual output from the running the
commands and differences will be reported.

The expected output is matched against the output exactly, but after the
expected is exhausted, any extra output from the command will be ignored. The
output comes from the combined stdout and stderr of the command.

## CLI
By default, no blocks are run as doc tests. You have to specify which blocks to
run with flags. Any combination of flags can be provided. If the interactive
session is not parsed correctly the program will exit with an error.

`--run-highlighted-code-fences <language>`: Run code fences with the selected
languages. You can select more than one by providing this flag multiple times.

`--run-code-fences`: Run code fences with the no selected language (just
backticks).

`--run-code-blocks`: Run code blocks. Code blocks are any bit of code indented
with 4 spaces or a tab.
