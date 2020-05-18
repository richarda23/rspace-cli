## rspace-client eln addFolder

Creates a new Folder

### Synopsis


Create a new Folder, with an optional name and parent folder ID
	

```
rspace-client eln addFolder [flags]
```

### Examples

```

// make a new folder in folder with id FL1234
rspace eln addFolder --name MyFolder --folder FL1234
	
```

### Options

```
  -p, --folder string   An id for the folder that will contain the new folder
  -h, --help            help for addFolder
  -n, --name string     A name for the folder
```

### Options inherited from parent commands

```
      --config string         config file (default is $HOME/.rspace)
  -o, --outFile string        Output file for program output
  -f, --outputFormat string   Output format: one of 'json','table', 'csv' or 'quiet'  (default "table")
```

### SEE ALSO

* [rspace-client eln](rspace-client_eln.md)	 - Top-level command to work with RSpace ELN

###### Auto generated by spf13/cobra on 14-May-2020