Here are some ideas on how to use the CLI for various use-cases.

1. Uploading a folder of files to the Gallery

You have many files to upload and doing it through the web interface is tedious.

    rspace eln upload myFolder/ --recursive --add-summary

This will upload all files in directory (and its subdirectories) and upload them, generating a summary document with links to all the uploaded files.

If you want more control over what you upload you can combine with standard Unix tools, e.g.

    find . -name "*.jpg" | xargs -I file rspace eln upload file --add-summary

All the files uploaded will be placed in the 'ApiInbox' folder of their respective Galleries.

2. Downloading a folder of files  from the Gallery

You have a folder of files in RSpace Gallery and you want to download them all. Using the user interface to download each file one at a time is tiresome.

Use the `listTree` command to navigate through   file tree. Or, you can get the folder ID from the 'info' popup in the web application.

    rspace eln listTree 

Once you have a folder ID (in this example, 9), get the folder contents in 'quiet' mode. This just lists the IDs, which can then be piped into a second command to download the files

    rspace eln listTree --folder 9 -f quiet | xargs ./rs eln download

3. Sharing many items at once

You have a folder full of RSpace documents. You'd like to share them all with your group, but don't want them organised into a notebook, and sharing them one a time is tedious

To share 1 or more items we can use the `share` command. To share a single item:

    rspace eln share SD2345 --groups 123 --folder 567

will share document to the group with id '123' to folder with id '567'.


To share many items, we start with using `listTree` to identify the folder with the documents we want to share.

Now, we call `groups` to find the id of the group we want to share with, and the id of the group shared folder

    rspace eln groups

```
    Id      	Name                     	Type      	SharedFolderId  
    7       	g1                       	LAB_GROUP 	326
```

By listing the items to share in quiet mode, we can pipe the ids to the share command:

    rspace eln listTree --filter document,notebook --folder 2345 -f quiet | xargs -I item  share item --groups 7 --folder 326 --permission edit

Note that we include the `--filter` option - this is so that we don't attempt to share any subfolders, which is not possible in RSpace. 

If we wanted to share with 'read' permission only we can omit the `--permission` argument, or explicitly state `--permission read`.

4. Importing many Word / OpenOffice Documents

You have a many Word documents and you'd like to import them into RSpace as native RSpace documents. 

This works in a similar way to `upload`. You can supply one or more folders and files, and 
the command will filter for Word documents based on extension (.odt, .doc, .docx), then import them into RSpace.

    rspace eln importWord myfolder

If you're not quite sure what you'll end up importing, you can use the `--dry-run` flag which will show you what would be uploaded, but doesn't change anything on RSpace.

    rspace eln importWord myfolder --dry-run
    