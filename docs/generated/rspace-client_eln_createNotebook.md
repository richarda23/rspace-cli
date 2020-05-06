## rspace-client eln createNotebook

Creates a new notebook

### Synopsis

Create a new notebook, with an optional name and parent folder
	

```
rspace-client eln createNotebook [flags]
```

### Examples

```

		// create a new notebook 'MyNotebook' in folder FL1234
		rspace eln createNotebook --name MyNotebook --folder FL1234

		//create an unnamed notebook in home folder
		rspace eln createNotebook
	
```

### Options

```
  -p, --folder string   An id for the folder that will contain the new notebook
  -h, --help            help for createNotebook
  -n, --name string     A name for the notebook
```

### Options inherited from parent commands

```
      --config string         config file (default is $HOME/.rspace)
  -o, --outFile string        Output file for program output
  -f, --outputFormat string   Output format: one of 'json','table', 'csv' or 'quiet'  (default "table")
```

### SEE ALSO

* [rspace-client eln](rspace-client_eln.md)	 - Top-level command to work with RSpace ELN

###### Auto generated by spf13/cobra on 6-May-2020